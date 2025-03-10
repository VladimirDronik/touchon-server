package http

import (
	"encoding/json"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
	"touchon-server/internal/helpers"
	"touchon-server/internal/model"
	"touchon-server/internal/objects"
	"touchon-server/internal/store"
	memStore "touchon-server/internal/store/memstore"
	libHelpers "touchon-server/lib/helpers"
	"touchon-server/lib/interfaces"
)

type getObjectsTypesResponseItem struct {
	Category string `json:"category"`
	Type     string `json:"type"`
	Name     string `json:"name"`
}

// Получение категорий типов объектов
// @Summary Получение категорий типов объектов
// @Tags Objects
// @Description Получение категорий типов объектов
// @ID GetObjectsTypes
// @Produce json
// @Param tags query string false "Тэги, через точку с запятой, либо пустое значение" example(sensor; htu21d; temperature)
// @Success      200 {object} http.Response[[]getObjectsTypesResponseItem]
// @Failure      400 {object} http.Response[any]
// @Failure      500 {object} http.Response[any]
// @Router /objects/types [get]
func (o *Server) handleGetObjectsTypes(ctx *fasthttp.RequestCtx) (interface{}, int, error) {
	tags := libHelpers.PrepareTags(helpers.GetParam(ctx, "tags"), ";")
	m := objects.GetCategoriesAndTypes()

	r := make([]*getObjectsTypesResponseItem, 0, len(m))
	for objCat, cat := range m {
		// Эти объекты отдельно через API не создаются
		if objCat == model.CategoryPort || objCat == model.CategorySensorValue || objCat == model.CategoryServer {
			continue
		}

		for objType, objectAttr := range cat {
			if len(tags) == 0 || compareTags(tags, objectAttr.Tags) {
				r = append(r, &getObjectsTypesResponseItem{
					Category: objCat,
					Type:     objType,
					Name:     objectAttr.Name,
				})
			}
		}
	}

	sort.Slice(r, func(i, j int) bool {
		switch {
		case r[i].Category != r[j].Category:
			return r[i].Category < r[j].Category
		default:
			return r[i].Type < r[j].Type
		}
	})

	return r, http.StatusOK, nil
}

// compareTags функция находит вхождение введенных тэгов и тэгов, которые есть у объекта
func compareTags(inputsTags []string, objectTags []string) bool {
	trueCnt := 0

	for _, inputTag := range inputsTags {
		inputTag = strings.Trim(inputTag, " ")

		for _, objectTag := range objectTags {
			if inputTag == objectTag {
				trueCnt++
				break
			}
		}
	}

	return trueCnt == len(inputsTags)
}

// Получение модели объекта
// @Summary Получение модели объекта
// @Tags Objects
// @Description Получение модели объекта
// @ID GetObjectModel
// @Produce json
// @Param category query string true "Категория объекта" Enums(controller,sensor,regulator,generic_input,relay,modbus,conditioner)
// @Param type query string true "Тип объекта" example(mega_d,htu21d,regulator,generic_input,relay,wb_mrm2_mini)
// @Success      200 {object} http.Response[objects.ObjectModel]
// @Failure      400 {object} http.Response[any]
// @Failure      500 {object} http.Response[any]
// @Router /objects/model [get]
func (o *Server) handleGetObjectModel(ctx *fasthttp.RequestCtx) (interface{}, int, error) {
	objCat := helpers.GetParam(ctx, "category")
	if objCat == "" {
		return nil, http.StatusBadRequest, errors.New("category is empty")
	}

	objType := helpers.GetParam(ctx, "type")
	if objType == "" {
		return nil, http.StatusBadRequest, errors.New("type is empty")
	}

	obj, err := objects.GetObjectModel(objCat, objType, true)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	// Внутренние объекты не создаем через API
	if obj.GetFlags().Has(objects.CreationForbidden) {
		return nil, http.StatusBadRequest, objects.CreationForbidden.Err()
	}

	return obj, http.StatusOK, nil
}

// Получение объекта
// @Summary Получение объекта
// @Tags Objects
// @Description Получение объекта
// @ID GetObject
// @Produce json
// @Param id path int true "ID объекта"
// @Param without_children query bool true "Без дочерних объектов" Enums(true,false)
// @Success      200 {object} http.Response[objects.ObjectModel]
// @Failure      400 {object} http.Response[any]
// @Failure      500 {object} http.Response[any]
// @Router /objects/{id} [get]
func (o *Server) handleGetObject(ctx *fasthttp.RequestCtx) (interface{}, int, error) {
	objectID, err := helpers.GetUintPathParam(ctx, "id")
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	withoutChildren, err := helpers.GetBoolParam(ctx, "without_children", false)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	objModel, err := objects.LoadObject(objectID, "", "", !withoutChildren)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	return objModel, http.StatusOK, nil
}

// Получение объекта по его свойствам
// @Summary Получение объекта по его свойствам
// @Tags Objects
// @Description Получение объекта по его свойствам
// @ID GetObjectByProps
// @Produce json
// @Param props query string true "Массив свойств объекта" default(type=i2c,number=0)
// @Param parent_id query string false "ID родительского объекта" default(2)
// @Success      200 {object} http.Response[int]
// @Failure      400 {object} http.Response[any]
// @Failure      500 {object} http.Response[any]
// @Router /objects/by_props [get]
func (o *Server) handleGetObjectByProps(ctx *fasthttp.RequestCtx) (interface{}, int, error) {
	props, err := helpers.GetMapParam(ctx, "props")
	if err != nil {
		return nil, http.StatusBadRequest, errors.Wrap(err, "param object_id is not valid")
	}

	if len(props) == 0 {
		return nil, http.StatusBadRequest, errors.New("param props is empty")
	}

	parentID, err := helpers.GetIntParam(ctx, "parent_id")
	if err != nil {
		return nil, http.StatusBadRequest, errors.Wrap(err, "param object_id is not valid")
	}

	objectID, err := store.I.ObjectRepository().GetObjectIDByProps(props, parentID)
	if err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(err, "GetObjectIDByProps")
	}

	return objectID, http.StatusOK, nil
}

// Удаление объекта
// @Summary Удаление объекта
// @Tags Objects
// @Description Удаление объекта
// @ID DeleteObject
// @Produce json
// @Param id path int true "ID объекта"
// @Success      200 {object} http.Response[any]
// @Failure      400 {object} http.Response[any]
// @Failure      500 {object} http.Response[any]
// @Router /objects/{id} [delete]
func (o *Server) handleDeleteObject(ctx *fasthttp.RequestCtx) (interface{}, int, error) {
	objectID, err := helpers.GetUintPathParam(ctx, "id")
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	objModel, err := objects.LoadObject(objectID, "", "", false)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	if objModel.GetFlags().Has(objects.DeletionForbidden) {
		return nil, http.StatusBadRequest, objects.DeletionForbidden.Err()
	}

	return o.DeleteObject(objectID)
}

func (o *Server) DeleteObject(objectID int) (interface{}, int, error) {
	//У порта обнуляем тип и режим
	objModel, err := objects.LoadObject(objectID, "", "", true)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	addressObject, _ := objModel.GetProps().GetIntValue("address")
	addressObjectString := strconv.Itoa(addressObject)
	if addressObject == 0 {
		addressObjectString, _ = objModel.GetProps().GetStringValue("address")
	}

	interfaceObject, _ := objModel.GetProps().GetStringValue("interface")

	//ищем другие объекты с таким же адресом
	_, relatedObjects, err := objects.FindRelatedObjects(addressObjectString, interfaceObject, objectID, objModel.GetType())
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	var objectToReset = make(map[int]string)
	objectToReset[objectID] = addressObjectString

	objects.ResetPortToDefault(objectToReset, relatedObjects)

	if err != nil {
		//TODO: тут сформировать запись в лог, что не могли изменить состояние порта на дефолтное и убрать вывод ошибки, чтобы неуспешность действия не было критичным
		return nil, http.StatusBadRequest, err
	}

	if err := memStore.I.DeleteObject(objectID); err != nil {
		return nil, http.StatusBadRequest, err
	}

	if err := store.I.ObjectRepository().DelObject(objectID); err != nil {
		return nil, http.StatusBadRequest, err
	}

	if err := objModel.DeleteChildren(); err != nil {
		return nil, http.StatusBadRequest, err
	}

	// Удаление действий в action-router

	// Удаляем все возможные события объекта
	if err := store.I.EventsRepo().DeleteEvent(interfaces.TargetTypeObject, objectID, "all"); err != nil {
		return nil, http.StatusInternalServerError, err
	}

	// Удаляем все действия для сторонних событий, где может фигурировать объект
	if err := store.I.EventActionsRepo().DeleteActionByObject(interfaces.TargetTypeObject, objectID); err != nil {
		return nil, http.StatusInternalServerError, err
	}

	// Удаляем все действия крона для объекта
	if err := store.I.CronRepo().DeleteTask(objectID, interfaces.TargetTypeObject); err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return nil, http.StatusOK, nil
}

// Получение всех тегов
// @Summary Получение всех тегов
// @Tags Objects
// @Description Получение всех тегов
// @ID GetAllObjectsTags
// @Produce json
// @Param related query bool true "Только те тэги, которые привязаны к созданным объектам" enums(true, false)
// @Success      200 {object} http.Response[[]string]
// @Failure      400 {object} http.Response[any]
// @Failure      500 {object} http.Response[any]
// @Router /objects/tags [get]
func (o *Server) handleGetAllObjectsTags(ctx *fasthttp.RequestCtx) (interface{}, int, error) {
	related, err := helpers.GetBoolParam(ctx, "related", false)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	if !related {
		var tags = make(map[string]int)
		catAndTypes := objects.GetCategoriesAndTypes()
		for objCat, _ := range catAndTypes {
			for _, obj := range catAndTypes[objCat] {
				for _, tag := range obj.Tags {
					tags[tag] = 0
				}
			}
		}
		return tags, http.StatusOK, nil
	}

	tags, err := store.I.ObjectRepository().GetAllTags()
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return tags, http.StatusOK, nil
}

// Получение состояния объекта
// @Summary Получение состояния объекта
// @Tags Objects
// @Description Получение состояния объекта
// @ID GetState
// @Produce json
// @Param id path int true "ID объекта" default(391)
// @Success      200 {object} http.Response[any]
// @Failure      400 {object} http.Response[any]
// @Failure      500 {object} http.Response[any]
// @Router /objects/{id}/state [get]
func (o *Server) handleGetObjectState(ctx *fasthttp.RequestCtx) (interface{}, int, error) {
	objectID, err := helpers.GetUintPathParam(ctx, "id")
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	obj, err := memStore.I.GetObject(objectID)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	state, err := obj.GetState()
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	return state, http.StatusOK, nil
}

// Запуск метода объекта
// @Summary Запуск метода объекта
// @Tags Objects
// @Description Запуск метода объекта
// @ID ExecMethod
// @Produce json
// @Param id path int true "ID объекта" default(391)
// @Param method path string true "Название метода" default(check)
// @Param args body object false "Параметры запуска метода" default({})
// @Success      200 {object} http.Response[any]
// @Failure      400 {object} http.Response[any]
// @Failure      500 {object} http.Response[any]
// @Router /objects/{id}/exec/{method} [post]
func (o *Server) handleExecMethod(ctx *fasthttp.RequestCtx) (interface{}, int, error) {
	objectID, err := helpers.GetUintPathParam(ctx, "id")
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	methodName := helpers.GetPathParam(ctx, "method")

	obj, err := memStore.I.GetObject(objectID)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	method, err := obj.GetMethods().Get(methodName)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	args := make(map[string]interface{}, 10)
	body := ctx.Request.Body()
	if len(body) >= 2 {
		if err := json.Unmarshal(ctx.Request.Body(), &args); err != nil {
			return nil, http.StatusBadRequest, err
		}
	}

	r, err := method.Func(args)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return r, http.StatusOK, nil
}
