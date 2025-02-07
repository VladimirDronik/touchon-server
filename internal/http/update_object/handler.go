package update_object

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
	"touchon-server/internal/context"
	"touchon-server/internal/model"
	"touchon-server/internal/objects"
	memStore "touchon-server/internal/store/memstore"
	_ "touchon-server/lib/http/server"
)

// Обновление объекта
// @Summary Обновление объекта
// @Tags Objects
// @Description Обновление объекта
// @ID UpdateObject
// @Accept json
// @Produce json
// @Param object body Request true "Объект"
// @Success      200 {object} server.Response[Response]
// @Failure      400 {object} server.Response[any]
// @Failure      500 {object} server.Response[any]
// @Router /objects [put]
func Handler(ctx *fasthttp.RequestCtx) (_ interface{}, _ int, e error) {
	accessLevel, err := context.GetAccessLevel(ctx)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	req := &Request{}
	if err := json.Unmarshal(ctx.Request.Body(), req); err != nil {
		return nil, http.StatusBadRequest, err
	}

	objModel, err := objects.LoadObject(req.ID, "", "", model.ChildTypeAll)
	if err != nil {
		return nil, 0, err
	}

	if req.ParentID > 0 {
		objModel.SetParentID(&req.ParentID)
	} else {
		objModel.SetStatus(model.StatusDisabled)
	}

	objModel.SetZoneID(req.ZoneID)
	objModel.SetName(req.Name)

	for k, v := range req.Props {
		dstProp, err := objModel.GetProps().Get(k)
		if err != nil {
			return nil, http.StatusBadRequest, err
		}

		if !dstProp.Editable.Check(accessLevel, objModel.GetProps()) {
			continue
		}

		//Если меняется период опроса у датчика
		if dstProp.Code == "update_interval" {
			if _, err := updateSensorCronTask(req); err != nil {
				return nil, http.StatusBadRequest, err
			}
		}

		// TODO: убрал NPE, нужно провести рефакторинг данного блока
		//Если меняем адрес размещения устройства, то проверяем возможность поменять настройки порта на контроллере
		if dstProp.Code == "address" && objModel.GetCategory() != "controller" {
			err := func() error {
				interfaceConnection, err := objModel.GetProps().Get("interface") //req.Props["interface"].(string)
				if err != nil {
					return nil
				}

				interfaceConnectionString, err := interfaceConnection.GetStringValue()
				if err != nil {
					return nil
				}

				newAddress, _ := req.Props["address"].(string)
				title := "[" + strconv.Itoa(req.ID) + "] " + req.Name

				//Получаем тип объекта
				objectType := objModel.GetType()
				objectID := objModel.GetID()

				//Переводим старый порт в состояние дефолта (NC)
				oldAddress, _ := dstProp.GetIntValue()
				oldAddressString := strconv.Itoa(oldAddress)
				if oldAddress == 0 {
					oldAddressString, err = dstProp.GetStringValue()
					if err != nil {
						return nil
					}
				}

				if newAddress != oldAddressString && newAddress != "" {
					//Ищем все устройства, которые висят на данном порту
					objectsToReset, relatedObjects, err := objects.FindRelatedObjects(newAddress, interfaceConnectionString, objectID, objectType)
					if err != nil {
						return err
					}

					objects.ResetParentAndAddress(objectsToReset)
					objects.ResetPortToDefault(objectsToReset, relatedObjects)
				}

				//Настраиваем новый порт
				options, err := objects.FillOptions(objectType, req.Props)
				if err != nil {
					return err
				}

				err = objects.ConfigureDevice(interfaceConnectionString, newAddress, options, title)
				if err != nil {
					return err
				}

				return nil
			}()

			if err != nil {
				return nil, http.StatusBadRequest, err
			}
		}

		if err := dstProp.SetValue(v); err != nil {
			return nil, http.StatusBadRequest, err
		}
	}

	if err := objModel.GetProps().Check(accessLevel); err != nil {
		return nil, http.StatusBadRequest, err
	}

	if len(req.Children) > 0 {
		if err := setChildrenProps(objModel.GetChildren(), req.Children, accessLevel); err != nil {
			return nil, http.StatusBadRequest, err
		}
	}

	if err := objModel.Save(); err != nil {
		return nil, http.StatusBadRequest, err
	}

	if err := memStore.I.SaveObject(objModel); err != nil {
		return nil, http.StatusBadRequest, err
	}

	return nil, http.StatusOK, nil
}

func setChildrenProps(objModelChildren *objects.Children, children []Child, accessLevel model.AccessLevel) error {
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
