package http

import (
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
	"touchon-server/internal/model"
	"touchon-server/internal/objects"
	"touchon-server/internal/store"
	memStore "touchon-server/internal/store/memstore"
	"touchon-server/lib/helpers"
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
	tags := helpers.GetParam(ctx, "tags")
	m := objects.GetCategoriesAndTypes()

	tagsSlice := strings.Split(tags, ";")

	// Эти объекты отдельно через API не создаются
	delete(m, model.CategoryPort)
	delete(m, model.CategorySensorValue)

	r := make([]*getObjectsTypesResponseItem, 0, len(m))
	for objCat, cat := range m {
		for objType, objectAttr := range cat {
			if compareTags(tagsSlice, objectAttr.Tags) || tags == "" {
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

	for _, objectTag := range objectTags {
		for _, inputTag := range inputsTags {
			if strings.Trim(inputTag, " ") == objectTag {
				trueCnt++
			}
		}
	}

	if trueCnt == len(inputsTags) {
		return true
	}
	return false
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

	obj, err := objects.GetObjectModel(model.Category(objCat), objType)
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

	children := model.ChildTypeAll
	if v := helpers.GetParam(ctx, "without_children"); v != "" {
		withoutChildren, err := strconv.ParseBool(v)
		if err != nil {
			return nil, http.StatusBadRequest, err
		}

		if withoutChildren {
			children = model.ChildTypeNobody
		}
	}

	objModel, err := objects.LoadObject(objectID, "", "", children)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	return objModel, http.StatusOK, nil
}

func parseGetObjectsParams(ctx *fasthttp.RequestCtx) (map[string]interface{}, error) {
	m := make(map[string]interface{}, 10)

	type P struct {
		QueryParamName string
		Type           string
		FieldName      string
	}

	params := []P{
		{"filter_by_id", "int", "id"},
		{"filter_by_parent_id", "int", "parent_id"},
		{"filter_by_zone_id", "int", "zone_id"},
		{"filter_by_category", "string", "category"},
		{"filter_by_type", "string", "type"},
		{"filter_by_name", "string", "name"},
		{"filter_by_status", "string", "status"},
		{"offset", "int", "offset"},
		{"limit", "int", "limit"},
		{"children", "int", "children"},
		{"type_struct", "string", "type_struct"},
		{"with_methods", "string", "with_methods"},
		{"type_children", "string", "type_children"},
		{"with_tags", "string", "with_tags"},
		{"simple_tree", "string", "simple_tree"},
	}

	for _, p := range params {
		switch p.Type {
		case "string":
			if v := helpers.GetParam(ctx, p.QueryParamName); v != "" {
				m[p.FieldName] = v
			}

		case "int":
			v, err := helpers.GetUintParam(ctx, p.QueryParamName)
			if err != nil {
				return nil, errors.Wrap(err, "parseGetObjectsParams")
			}
			if v > 0 || p.FieldName == "offset" || p.FieldName == "limit" {
				m[p.FieldName] = v
			}

		default:
			return nil, errors.Wrap(errors.Errorf("unexpected param type %q", p.Type), "parseGetObjectsParams")
		}
	}

	return m, nil
}

// Получение объекта по его свойствам
// @Summary Получение объекта по его свойствам
// @Tags Objects
// @Description Получение объекта по его свойствам
// @ID GetObjectByProps
// @Produce json
// @Param props query string false "Массив свойств объекта" example(address=1, interface=I2C)
// @Param parent_id query string false "ID родительского объекта"
// @Success      200 {object} http.Response[int]
// @Failure      400 {object} http.Response[any]
// @Failure      500 {object} http.Response[any]
// @Router /objects/by_props [get]
func (o *Server) handleGetObjectByProps(ctx *fasthttp.RequestCtx) (interface{}, int, error) {
	propsReq := helpers.GetParam(ctx, "props")
	parentID, err := helpers.GetIntParam(ctx, "parent_id")
	if err != nil {
		return nil, http.StatusBadRequest, errors.Wrap(err, "param object_id is not valid")
	}

	if propsReq == "" {
		return nil, http.StatusBadRequest, errors.New("props is empty")
	}

	propsSlice := strings.Split(propsReq, ",")
	if len(propsSlice) == 0 {
		return nil, http.StatusBadRequest, errors.New("props string is bad")
	}

	propsMap := make(map[string]string, 10)
	for _, p := range propsSlice {
		prop := strings.Split(p, "=")
		if len(prop) < 2 {
			return nil, http.StatusBadRequest, errors.New("props string is bad")
		}
		propsMap[prop[0]] = prop[1]
	}

	objectID, err := store.I.ObjectRepository().GetObjectIDByProps(propsMap, parentID)
	if err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(err, "GetObjectIDByProps")
	}

	return objectID, http.StatusOK, nil
}

type GetObjectsResponse struct {
	Total int                  `json:"total"`
	List  []*model.StoreObject `json:"list"`
}

// Получение списка объектов
// @Summary Получение объектов
// @Tags Objects
// @Description Получение объектов
// @ID GetObjects
// @Produce json
// @Param filter_by_id query string false "ID объекта"
// @Param filter_by_parent_id query string false "ID родительского объекта"
// @Param filter_by_zone_id query string false "ID зоны"
// @Param filter_by_category query string false "Категория" Enums(controller,sensor,regulator,generic_input, sensor_value)
// @Param filter_by_type query string false "Тип" example(mega_d, htu21d, regulator, generic_input)
// @Param filter_by_name query string false "Название"
// @Param filter_by_status query string false "Статус" Enums(ON,OFF,Enable,N/A)
// @Param tags query string false "Тэги"
// @Param offset query string false "Смещение" default(0)
// @Param limit query string false "Лимит" default(20)
// @Param children query string false "Возвращать дочерние объекты (0-без детей, 1-дети, 2-дети+внуки и т.д.)" default(1)
// @Param type_children query string false "Тип выводимых дочерних элементов (all - все, internal - только внутр., external - только внешние)" Enums(all, internal, external)
// @Param type_struct query string false "Тип структуры в ответе" Enums(easy, full)
// @Param with_methods query string false "Добавить методы в структуру" Enums(true, false)
// @Param with_tags query string false "Добавлять тэги в структуру" Enums(true, false)
// @Param simple_tree query string false "Дерево объектов будет строиться от отфильтрованного объекта вниз, но не вверх" Enums(true, false)
// @Success      200 {object} http.Response[GetObjectsResponse]
// @Failure      400 {object} http.Response[any]
// @Failure      500 {object} http.Response[any]
// @Router /objects [get]
func (o *Server) handleGetObjects(ctx *fasthttp.RequestCtx) (interface{}, int, error) {
	tags := helpers.PrepareTags(helpers.GetParam(ctx, "tags"))
	withTagsFlag := true

	params, err := parseGetObjectsParams(ctx)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	if v, ok := params["category"]; ok {
		objCat := v.(string)

		m := objects.GetCategoriesAndTypes()
		if _, ok := m[objCat]; !ok || objCat == string(model.CategoryPort) || objCat == string(model.CategorySensorValue) {
			return nil, http.StatusInternalServerError, errors.Errorf("bad category value")
		}
	}

	offset, _ := params["offset"].(int)
	limit, _ := params["limit"].(int)
	childrenAge, _ := params["children"].(int)
	typeStruct := params["type_struct"]
	withMethods := params["with_methods"]
	typeChildren := params["type_children"]
	withTags := params["with_tags"]
	simpleTree := params["simple_tree"]
	delete(params, "offset")
	delete(params, "limit")
	delete(params, "children")
	delete(params, "type_struct")
	delete(params, "with_methods")
	delete(params, "type_children")
	delete(params, "with_tags")
	delete(params, "simple_tree")

	if limit == 0 {
		limit = 20
	}

	var childType model.ChildType
	switch typeChildren {
	case "all":
		childType = model.ChildTypeAll
	case "external":
		childType = model.ChildTypeExternal
	case "internal":
		childType = model.ChildTypeInternal
	default:
		childType = model.ChildTypeExternal
	}

	rows, err := store.I.ObjectRepository().GetObjects(params, tags, offset, limit)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	// TODO фильтруем внутренние объекты

	// TODO нужно подсчитывать кол-во за вычетом внутренних объектов
	total, err := store.I.ObjectRepository().GetTotal(params, tags)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	// Создаем общий список всех эл-ов
	m := make(map[int]*model.StoreObject, len(rows))
	items := make([]*model.StoreObject, 0, len(m))
	for _, row := range rows {
		if withMethods == "true" { // если опция показа методов включена
			obj, err := objects.GetObjectModel(row.Category, row.Type)
			if err != nil {
				return nil, http.StatusBadRequest, err
			}

			methods := obj.GetMethods().GetAll()

			for _, method := range methods {
				row.Methods = append(row.Methods, model.Method{
					Name:        method.Name,
					Description: method.Description,
				})
			}
		}

		if withTags == "false" {
			row.Tags = nil
			withTagsFlag = false
		}

		if simpleTree == "true" {
			row.ParentID = nil
		}

		m[row.ID] = row

		if typeStruct == "easy" { //если упрощенный вывод включен, то убираем всё лишнее
			row.ParentID = nil
			row.Status = ""
			row.Category = ""
			row.Children = nil
			row.ZoneID = nil
			items = append(items, row)
		}
	}

	if typeStruct == "easy" { //если упрощенный вывод включен, то сразу выводим плоскую модель без вложенностей
		return GetObjectsResponse{
			Total: total,
			List:  items,
		}, http.StatusOK, nil
	}

	//Рекурсивно загружаем детей
	if err := loadChildren(m, rows, store.I.ObjectRepository(), childrenAge, childType, withTagsFlag); err != nil {
		return nil, http.StatusInternalServerError, err
	}

	// Рекурсивно загружаем родителей
	if err := loadParents(m, rows, store.I.ObjectRepository()); err != nil {
		return nil, http.StatusInternalServerError, err
	}

	// Детей добавляем к родителям
	for _, item := range m {
		if item.ParentID != nil {
			parent, ok := m[*item.ParentID]
			if !ok {
				return nil, http.StatusInternalServerError, errors.Errorf("Parent(ID=%d) not found for object(ID=%d)", *item.ParentID, item.ID)
			}

			parent.Children = append(parent.Children, item)
		}
	}

	//Выбираем эл-ты верхнего уровня
	items = make([]*model.StoreObject, 0, len(m))
	for _, item := range m {
		if item.ParentID == nil {
			items = append(items, item)
		}
	}

	sortObjectsTree(items)

	return GetObjectsResponse{
		Total: total,
		List:  items,
	}, http.StatusOK, nil
}

// Получение объектов по тегам
// @Summary Получение объектов по тегам
// @Tags Objects
// @Description Получение объектов  по тегам
// @ID GetObjectsByTags
// @Produce json
// @Param tags query string true "Тэги" default(controller,mega_d)
// @Param offset query string false "Смещение" default(0)
// @Param limit query string false "Лимит" default(20)
// @Param children query string false "Возвращать дочерние объекты (0-без детей, 1-дети, 2-дети+внуки и т.д.)" default(1)
// @Success      200 {object} http.Response[GetObjectsResponse]
// @Failure      400 {object} http.Response[any]
// @Failure      500 {object} http.Response[any]
// @Router /objects/by_tags [get]
func (o *Server) handleGetObjectsByTags(ctx *fasthttp.RequestCtx) (interface{}, int, error) {
	tags := helpers.PrepareTags(helpers.GetParam(ctx, "tags"))
	if len(tags) == 0 {
		return nil, http.StatusBadRequest, errors.New("tags is empty")
	}

	offset, err := helpers.GetUintParam(ctx, "offset")
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	limit, err := helpers.GetUintParam(ctx, "limit")
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	if limit < 1 {
		return nil, http.StatusBadRequest, errors.New("limit == 0")
	}

	childrenAge, err := helpers.GetUintParam(ctx, "children")
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	rows, err := store.I.ObjectRepository().GetObjectsByTags(tags, offset, limit, model.ChildTypeAll)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	total, err := store.I.ObjectRepository().GetTotalByTags(tags, model.ChildTypeAll)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	// Создаем общий список всех эл-ов
	m := make(map[int]*model.StoreObject, len(rows))
	for _, row := range rows {
		m[row.ID] = row
	}

	// Рекурсивно загружаем детей
	if err := loadChildren(m, rows, store.I.ObjectRepository(), childrenAge, model.ChildTypeAll, true); err != nil {
		return nil, http.StatusInternalServerError, err
	}

	// Рекурсивно загружаем родителей
	if err := loadParents(m, rows, store.I.ObjectRepository()); err != nil {
		return nil, http.StatusInternalServerError, err
	}

	// Детей добавляем к родителям
	for _, item := range m {
		if item.ParentID != nil {
			parent, ok := m[*item.ParentID]
			if !ok {
				return nil, http.StatusInternalServerError, errors.Errorf("Parent(ID=%d) not found for object(ID=%d)", *item.ParentID, item.ID)
			}

			parent.Children = append(parent.Children, item)
		}
	}

	// Выбираем эл-ты верхнего уровня
	items := make([]*model.StoreObject, 0, len(m))
	for _, item := range m {
		if item.ParentID == nil {
			items = append(items, item)
		}
	}

	sortObjectsTree(items)

	return GetObjectsResponse{
		Total: total,
		List:  items,
	}, http.StatusOK, nil
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

	objModel, err := objects.LoadObject(objectID, "", "", model.ChildTypeNobody)
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
	objModel, err := objects.LoadObject(objectID, "", "", model.ChildTypeAll)
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
// @Param related query bool true "Только те тэги, которые привязаны к созданным объектам" default(true)
// @Success      200 {object} http.Response[[]string]
// @Failure      400 {object} http.Response[any]
// @Failure      500 {object} http.Response[any]
// @Router /objects/tags [get]
func (o *Server) handleGetAllObjectsTags(ctx *fasthttp.RequestCtx) (interface{}, int, error) {
	related, err := helpers.GetBoolParam(ctx, "related")
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	if related == false {
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
