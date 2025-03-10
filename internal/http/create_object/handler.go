package create_object

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
	"touchon-server/internal/g"
	"touchon-server/internal/helpers"
	"touchon-server/internal/model"
	"touchon-server/internal/objects"
	"touchon-server/internal/store"
	memStore "touchon-server/internal/store/memstore"
	_ "touchon-server/lib/http/server"
	"touchon-server/lib/interfaces"
)

// Создание объекта (с действиями)
// @Summary Создание объекта (с действиями)
// @Tags Objects
// @Description Создание объекта (с действиями)
// @ID CreateObject
// @Accept json
// @Produce json
// @Param object body Request true "Объект"
// @Success      200 {object} server.Response[Response]
// @Failure      400 {object} server.Response[any]
// @Failure      500 {object} server.Response[any]
// @Router /objects [post]
func Handler(ctx *fasthttp.RequestCtx) (_ interface{}, _ int, e error) {
	req := &Request{}
	if err := json.Unmarshal(ctx.Request.Body(), req); err != nil {
		return nil, http.StatusBadRequest, err
	}

	req.Object.Enabled = true

	// Проверяем названия событий на валидность
	for _, item := range req.Events {
		// TODO
		//if _, err := event.GetMaker(item.Name); err != nil {
		//	return nil, http.StatusBadRequest, err
		//}

		if len(item.Actions) == 0 {
			return nil, http.StatusBadRequest, errors.Errorf("event %q: actions list is empty", item.Name)
		}
	}

	// Сохраняем объект в базу и memstore
	objectID, err := createObject(req)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	resp := &Response{ID: objectID}

	// Если не закончили транзакцию, удаляем созданный объект
	defer func() {
		if e != nil {
			if err := deleteObject(objectID); err != nil {
				e = err
				g.Logger.Error(err)
			}
		}
	}()

	if req.Object.Category != model.CategoryController && req.Object.Category != model.CategoryRS485 {
		status, err := deviceConfiguration(*req, objectID)
		if err != nil {
			return nil, status, err
		}
	}

	// Если событий нет, то уходим
	if len(req.Events) == 0 {
		return resp, http.StatusOK, nil
	}

	// Если транзакцию не закончили, удаляем события со всеми действиями
	defer func() {
		if e != nil {
			for _, ev := range req.Events {
				if err := store.I.EventsRepo().DeleteEvent(interfaces.TargetTypeObject, objectID, ev.Name); err != nil {
					e = err
					g.Logger.Error(err)
				}
			}
		}
	}()

	for _, ev := range req.Events {
		for _, act := range ev.Actions {
			if err := g.HttpServer.CreateEventAction(interfaces.TargetTypeObject, objectID, ev.Name, act); err != nil {
				return nil, http.StatusInternalServerError, err
			}
		}
	}

	return resp, http.StatusOK, nil
}

func createObject(req *Request) (int, error) {
	objModel, err := objects.LoadObject(0, req.Object.Category, req.Object.Type, true)
	if err != nil {
		return 0, errors.Wrap(err, "createObject")
	}

	if objModel.GetFlags().Has(objects.CreationForbidden) {
		return 0, errors.Wrap(objects.CreationForbidden.Err(), "createObject")
	}

	objModel.SetParentID(req.Object.ParentID)
	objModel.SetZoneID(req.Object.ZoneID)
	objModel.SetName(req.Object.Name)
	objModel.SetEnabled(req.Object.Enabled)

	// Выставляем сначала значения по умолчанию
	for _, p := range objModel.GetProps().GetAll().GetValueList() {
		if p.DefaultValue != nil {
			if err := p.SetValue(p.DefaultValue); err != nil {
				return 0, errors.Wrap(err, "createObject")
			}
		}
	}

	for k, v := range req.Object.Props {
		dstProp, err := objModel.GetProps().Get(k)
		if err != nil {
			return 0, errors.Wrap(err, "createObject")
		}

		if !dstProp.Editable.Check(objModel.GetProps()) {
			continue
		}

		if err := dstProp.SetValue(v); err != nil {
			return 0, errors.Wrap(err, "createObject")
		}
	}

	if err := objModel.GetProps().Check(); err != nil {
		return 0, errors.Wrap(err, "createObject")
	}

	if err := setChildrenDefaultPropValues(objModel.GetChildren()); err != nil {
		return 0, errors.Wrap(err, "createObject")
	}

	if len(req.Object.Children) > 0 {
		if err := setChildrenProps(objModel.GetChildren(), req.Object.Children); err != nil {
			return 0, errors.Wrap(err, "createObject")
		}
	}

	if err := objModel.Save(); err != nil {
		return 0, errors.Wrap(err, "createObject")
	}

	if err := memStore.I.SaveObject(objModel); err != nil {
		return 0, errors.Wrap(err, "createObject")
	}

	return objModel.GetID(), nil
}

func setChildrenProps(objModelChildren *objects.Children, children []*Child) error {
	if objModelChildren.Len() != len(children) {
		return errors.Wrap(errors.Errorf("objModelChildren.Len() != len(children), %d != %d", objModelChildren.Len(), len(children)), "setChildrenProps")
	}

	modelChildren := objModelChildren.GetAll()

	for i, child := range children {
		objModel := modelChildren[i]

		for k, v := range child.Props {
			dstProp, err := objModel.GetProps().Get(k)
			if err != nil {
				return errors.Wrap(err, "setChildrenProps")
			}

			if !dstProp.Editable.Check(objModel.GetProps()) {
				continue
			}

			if err := dstProp.SetValue(v); err != nil {
				return errors.Wrap(err, "setChildrenProps")
			}
		}

		if err := objModel.GetProps().Check(); err != nil {
			return errors.Wrap(err, "setChildrenProps")
		}

		if len(child.Children) > 0 {
			if err := setChildrenProps(objModel.GetChildren(), child.Children); err != nil {
				return err
			}
		}
	}

	return nil
}

// setChildrenDefaultPropValues выставляет значения св-в по умолчанию
func setChildrenDefaultPropValues(objModelChildren *objects.Children) error {
	for _, objModel := range objModelChildren.GetAll() {
		for _, p := range objModel.GetProps().GetAll().GetValueList() {
			if p.GetValue() == nil && p.DefaultValue != nil {
				if err := p.SetValue(p.DefaultValue); err != nil {
					return errors.Wrap(err, "setChildrenDefaultPropValues")
				}
			}
		}

		if err := objModel.GetProps().Check(); err != nil {
			return errors.Wrap(err, "setChildrenDefaultPropValues")
		}

		if objModel.GetChildren().Len() > 0 {
			if err := setChildrenDefaultPropValues(objModel.GetChildren()); err != nil {
				return err
			}
		}
	}

	return nil
}

func deleteObject(objectID int) error {
	if err := memStore.I.DeleteObject(objectID); err != nil {
		return errors.Wrap(err, "deleteObject")
	}

	if err := store.I.ObjectRepository().DelObject(objectID); err != nil {
		return errors.Wrap(err, "deleteObject")
	}

	return nil
}

func deviceConfiguration(req Request, objectID int) (_ int, e error) {
	fastConfig := false
	//Настройка портов контроллера, либо конфигурирование другого устройства, на котором располагается объект
	interfaceConnection, _ := req.Object.Props["interface"].(string)
	addressObject, _ := req.Object.Props["address"].(string)
	typeObject := req.Object.Type

	//ищем контроллер, если опция быстрого конфига выключена, то не конфигурим порты на лету
	objContr, err := objects.LoadObject(*req.Object.ParentID, model.CategoryController, "", false)
	if err != nil {
		g.Logger.Error(err)
	}
	if objContr != nil {
		fastConfig, err = objContr.GetProps().GetBoolValue("fast_config")
		if err != nil {
			g.Logger.Error(err)
		}
	}

	if fastConfig == false {
		return 0, nil
	}

	//Проверяем назначен ли адрес на какой-либо другой объект
	objectsToReset, _, err := objects.FindRelatedObjects(addressObject, interfaceConnection, objectID, typeObject)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	//для найденных объектов на контроллере подчищаем порты
	for _, addressReset := range objectsToReset {
		resetPorts := strings.Split(addressReset, ";")
		for _, resetPort := range resetPorts {
			if resetPort == addressObject {
				e = objects.ConfigureDevice("NC", resetPort, nil, "")
			}
		}
	}

	title := "[" + strconv.Itoa(objectID) + "]" + req.Object.Name
	options, err := objects.FillOptions(typeObject, req.Object.Props)
	if err != nil {
		return http.StatusBadRequest, err
	}
	e = objects.ConfigureDevice(interfaceConnection, addressObject, options, title)

	if err := helpers.ResetParentAndAddress(objectsToReset); err != nil {
		return http.StatusBadRequest, err
	}

	return http.StatusOK, nil
}
