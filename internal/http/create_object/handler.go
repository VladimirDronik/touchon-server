package create_object

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
	"touchon-server/internal/context"
	"touchon-server/internal/helpers"
	"touchon-server/internal/model"
	"touchon-server/internal/objects"
	"touchon-server/internal/store"
	memStore "touchon-server/internal/store/memstore"
	httpClient "touchon-server/lib/http/client"
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
	accessLevel, err := context.GetAccessLevel(ctx)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	req := &Request{}
	if err := json.Unmarshal(ctx.Request.Body(), req); err != nil {
		return nil, http.StatusBadRequest, err
	}

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
	objectID, err := createObject(req, accessLevel)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	resp := &Response{ID: objectID}

	// Если не закончили транзакцию, удаляем созданный объект
	defer func() {
		if e != nil {
			if err := deleteObject(objectID); err != nil {
				e = err
				context.Logger.Error(err)
			}
		}
	}()

	//Настройка портов контроллера, либо конфигурирование другого устройства, на котором располагается объект
	interfaceConnection, _ := req.Object.Props["interface"].(string)
	addressObject, _ := req.Object.Props["address"].(string)
	typeObject := req.Object.Type

	//Проверяем назначен ли адрес на какой-либо другой объект
	objectsToReset, _, err := objects.FindRelatedObjects(addressObject, interfaceConnection, objectID, typeObject)
	if err != nil {
		return nil, http.StatusInternalServerError, err
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
		return nil, http.StatusBadRequest, err
	}
	e = objects.ConfigureDevice(interfaceConnection, addressObject, options, title)

	if err := helpers.ResetParentAndAddress(objectsToReset); err != nil {
		return nil, http.StatusBadRequest, err
	}

	//Если объект является сенсором, то создаем в экшен-роутере действия для его проверки
	if req.Object.Category == "sensor" {
		_, e = createSensorCronTask(objectID, req)
	}

	// Если событий нет, то уходим
	if len(req.Events) == 0 {
		return resp, http.StatusOK, nil
	}

	arBaseUrl := "http://" + context.Config["action_router_addr"]

	// Если транзакцию не закончили, удаляем события со всеми действиями
	defer func() {
		if e != nil {
			for _, ev := range req.Events {
				params := map[string]string{
					"target_type": string(interfaces.TargetTypeObject),
					"target_id":   strconv.Itoa(objectID),
					"event_name":  ev.Name,
				}

				if _, err := httpClient.I.DoRequest("DELETE", arBaseUrl+"/events", params, nil, nil); err != nil {
					e = err
					context.Logger.Error(err)
				}
			}
		}
	}()

	for _, ev := range req.Events {
		params := map[string]string{
			"target_type": string(interfaces.TargetTypeObject),
			"target_id":   strconv.Itoa(objectID),
			"event_name":  ev.Name,
		}

		for _, act := range ev.Actions {
			if _, err := httpClient.I.DoRequest("POST", arBaseUrl+"/events/actions", params, nil, act); err != nil {
				return nil, http.StatusInternalServerError, err
			}
		}
	}

	return resp, http.StatusOK, nil
}

func createObject(req *Request, accessLevel model.AccessLevel) (int, error) {
	objModel, err := objects.LoadObject(0, req.Object.Category, req.Object.Type, model.ChildTypeAll)
	if err != nil {
		return 0, errors.Wrap(err, "createObject")
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

		if !dstProp.Editable.Check(accessLevel, objModel.GetProps()) {
			continue
		}

		if err := dstProp.SetValue(v); err != nil {
			return 0, errors.Wrap(err, "createObject")
		}
	}

	if err := objModel.GetProps().Check(accessLevel); err != nil {
		return 0, errors.Wrap(err, "createObject")
	}

	if err := setChildrenDefaultPropValues(objModel.GetChildren(), accessLevel); err != nil {
		return 0, errors.Wrap(err, "createObject")
	}

	if len(req.Object.Children) > 0 {
		if err := setChildrenProps(objModel.GetChildren(), req.Object.Children, accessLevel); err != nil {
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

func setChildrenProps(objModelChildren *objects.Children, children []*Child, accessLevel model.AccessLevel) error {
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

			if !dstProp.Editable.Check(accessLevel, objModel.GetProps()) {
				continue
			}

			if err := dstProp.SetValue(v); err != nil {
				return errors.Wrap(err, "setChildrenProps")
			}
		}

		if err := objModel.GetProps().Check(accessLevel); err != nil {
			return errors.Wrap(err, "setChildrenProps")
		}

		if len(child.Children) > 0 {
			if err := setChildrenProps(objModel.GetChildren(), child.Children, accessLevel); err != nil {
				return err
			}
		}
	}

	return nil
}

// setChildrenDefaultPropValues выставляет значения св-в по умолчанию
func setChildrenDefaultPropValues(objModelChildren *objects.Children, accessLevel model.AccessLevel) error {
	for _, objModel := range objModelChildren.GetAll() {
		for _, p := range objModel.GetProps().GetAll().GetValueList() {
			if p.GetValue() == nil && p.DefaultValue != nil {
				if err := p.SetValue(p.DefaultValue); err != nil {
					return errors.Wrap(err, "setChildrenDefaultPropValues")
				}
			}
		}

		if err := objModel.GetProps().Check(accessLevel); err != nil {
			return errors.Wrap(err, "setChildrenDefaultPropValues")
		}

		if objModel.GetChildren().Len() > 0 {
			if err := setChildrenDefaultPropValues(objModel.GetChildren(), accessLevel); err != nil {
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
