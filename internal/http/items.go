package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
	"touchon-server/internal/g"
	"touchon-server/internal/model"
	"touchon-server/internal/objects"
	"touchon-server/internal/store"
	memStore "touchon-server/internal/store/memstore"
	"touchon-server/lib/events/item"
	"touchon-server/lib/helpers"
	"touchon-server/lib/interfaces"
	"touchon-server/lib/messages"
)

// Создание элемента
// @Security TokenAuth
// @Summary Создание элемента
// @Tags Items
// @Description Создание элемента
// @ID CreateItem
// @Accept json
// @Produce json
// @Param item body model.ViewItem true "Элемент"
// @Success      200 {object} Response[model.ViewItem]
// @Failure      400 {object} Response[any]
// @Failure      500 {object} Response[any]
// @Router /private/items [post]
func (o *Server) handleCreateItem(ctx *fasthttp.RequestCtx) (interface{}, int, error) {
	return o.saveItem(ctx.Request.Body())
}

// Создание датчика
// @Security TokenAuth
// @Summary Создание датчика
// @Tags Items
// @Description Создание датчика
// @ID CreateSensor
// @Accept json
// @Produce json
// @Param item body model.Sensor true "Датчик"
// @Success      200 {object} Response[model.Sensor]
// @Failure      400 {object} Response[any]
// @Failure      500 {object} Response[any]
// @Router /private/item/sensor [post]
func (o *Server) handleCreateSensor(ctx *fasthttp.RequestCtx) (interface{}, int, error) {
	sensor := &model.Sensor{}
	if err := json.Unmarshal(ctx.Request.Body(), sensor); err != nil {
		return nil, http.StatusInternalServerError, err
	}

	sensor.MinThreshold = nil
	sensor.MaxThreshold = nil

	sensorObj, err := memStore.I.GetObject(sensor.ObjectID)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	//sensorValue, err := store.I.ObjectRepository().GetObjectByParent(sensor.ObjectID, sensor.Type)
	//if err != nil {
	//	return nil, http.StatusInternalServerError, errors.Wrap(err, "Get valueSensor by sensorID")
	//}

	//Если включена регулировка датчика, то выставляем пороговые значения у параметра сенсора
	if sensor.Adjustment == true {
		//objModel, err := objects.LoadObject(sensorValue.ID, "", "", model.ChildTypeNobody)
		//if err != nil {
		//	return nil, http.StatusInternalServerError, errors.Wrap(err, "LoadObject By ID")
		//}
		//
		//minThreshold, err := objModel.GetProps().Get("min_threshold")
		//if err != nil {
		//	return nil, http.StatusInternalServerError, errors.Wrap(err, "GetProps For Object")
		//}
		//minThreshold.SetValue(sensor.MinThreshold)
		//
		//maxThreshold, err := objModel.GetProps().Get("max_threshold")
		//if err != nil {
		//	return nil, http.StatusInternalServerError, errors.Wrap(err, "GetProps For Object")
		//}
		//maxThreshold.SetValue(sensor.MaxThreshold)
		//
		//if err := objModel.Save(); err != nil {
		//	return nil, http.StatusInternalServerError, err
		//}
		//if err := memStore.I.SaveObject(objModel); err != nil {
		//	return nil, http.StatusInternalServerError, errors.Wrap(err, "updateSensorValue")
		//}
	}

	//Ищем регулятор для сенсора, если находим, то включаем, если не находим, то создаем
	//reg, err := store.I.ObjectRepository().GetObjectByParent(sensorValue.ID, "regulator")
	//if err != nil {
	//	return nil, http.StatusInternalServerError, errors.Wrap(err, "Get regulator by sensorID")
	//}
	//
	//var objRegulator objects.Object
	//regID := 0
	//
	//if &reg != nil {
	//	regID = reg.ID
	//}
	//
	//objRegulator, err = objects.LoadObject(regID, "regulator", "regulator", model.ChildTypeNobody)
	//if err != nil {
	//	return nil, http.StatusInternalServerError, errors.Wrap(err, "LoadObject By ID")
	//}
	//
	//objRegulator.SetEnabled(sensor.Adjustment)
	//
	//if err := objRegulator.Save(); err != nil {
	//	return nil, http.StatusInternalServerError, errors.Wrap(err, "Save regulator in storage")
	//}
	//if err := memStore.I.SaveObject(objRegulator); err != nil {
	//	return nil, http.StatusInternalServerError, errors.Wrap(err, "Save regulator in memory")
	//}

	item := &model.ViewItem{
		Type:    "sensor",
		Enabled: true,
		ZoneID:  &sensor.ZoneID,
		Title:   sensor.Title,
		Icon:    sensor.Icon,
		Params:  "{}",
		Auth:    "",
		Sort:    0,
	}

	if err := store.I.Items().SaveItem(item); err != nil {
		return nil, http.StatusInternalServerError, err
	}
	sensor.ViewItemID = item.ID

	event := &model.TrEvent{
		EventName:  "object.sensor.on_check",
		TargetType: "object",
		TargetID:   *sensorObj.GetParentID(),
		Enabled:    1,
	}

	eventID, err := store.I.Events().AddEvent(event)
	if err != nil || eventID == 0 {
		return nil, http.StatusInternalServerError, err
	}

	eventAction := &model.EventActions{
		EventID:    eventID,
		TargetType: "item",
		TargetID:   item.ID,
		Type:       "method",
		Name:       "set_value",
		Args:       "{\"param\":\"" + sensor.Type + "\"}",
	}

	_, err = store.I.Events().AddEventAction(eventAction)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return sensor, http.StatusOK, store.I.Devices().SaveSensor(sensor)
}

// Удаление датчика
// @Security TokenAuth
// @Summary Удаление датчика
// @Tags Items
// @Description Удаление датчика
// @ID DeleteSensor
// @Accept json
// @Produce json
// @Param id query int true "ID итема датчика"
// @Success      200 {object} Response[any]
// @Failure      400 {object} Response[any]
// @Failure      500 {object} Response[any]
// @Router /private/item/sensor [delete]
func (o *Server) handleDeleteSensor(ctx *fasthttp.RequestCtx) (interface{}, int, error) {
	itemID, err := helpers.GetUintParam(ctx, "id")
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	store.I.Events().DeleteEvent(itemID)
	store.I.Devices().DeleteSensor(itemID)

	return o.deleteItem(itemID)
}

func (o *Server) handleUpdateSensor(ctx *fasthttp.RequestCtx) (interface{}, int, error) {
	sensor := &model.Sensor{}
	if err := json.Unmarshal(ctx.Request.Body(), sensor); err != nil {
		return nil, http.StatusInternalServerError, err
	}
	//TODO: Сделать редактирование итема датчика
	return nil, 0, nil
}

type targetRequest struct {
	ItemID      int     `json:"item_id"`
	TargetValue float32 `json:"target"`
}

// Обновление target у датчика
// @Security TokenAuth
// @Summary Обновление target у датчика
// @Tags Items
// @Description Обновление target у датчика
// @ID UpdateTarget
// @Accept json
// @Produce json
// @Param item body targetRequest true "Датчик"
// @Success      200 {object} Response[any]
// @Failure      400 {object} Response[any]
// @Failure      500 {object} Response[any]
// @Router /private/item/sensor/value [patch]
func (o *Server) handleSetTargetSensor(ctx *fasthttp.RequestCtx) (interface{}, int, error) {
	targetReq := &targetRequest{}
	if err := json.Unmarshal(ctx.Request.Body(), targetReq); err != nil {
		return nil, http.StatusBadRequest, err
	}

	sensorItem, err := store.I.Devices().GetSensor(targetReq.ItemID)
	if err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(err, "Sensor not found: "+strconv.Itoa(targetReq.ItemID))
	}

	sensorObj, err := memStore.I.GetObject(sensorItem.ObjectID) //objects.LoadObject(sensorItem.ObjectID, "", "", model.ChildTypeAll)
	if err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(err, "Object not loaded: "+strconv.Itoa(sensorItem.ObjectID))
	}

	childrens := sensorObj.GetChildren().GetAll()
	for _, child := range childrens {
		child.SetEnabled(true)
		target, err := child.GetProps().Get("target_sp")
		if err != nil {
			return nil, http.StatusInternalServerError, errors.Wrap(err, "Property 'target_sp' not found for object: "+strconv.Itoa(sensorItem.ObjectID))
		}
		target.SetValue(targetReq.TargetValue)
		if err := child.Save(); err != nil {
			return nil, http.StatusInternalServerError, errors.Wrap(err, "Unable to save object: "+strconv.Itoa(sensorItem.ObjectID))
		}

		if err := memStore.I.SaveObject(child); err != nil {
			return nil, http.StatusInternalServerError, errors.Wrap(err, "Unable to save to memory object: "+strconv.Itoa(sensorItem.ObjectID))
		}
		break
	}

	return nil, http.StatusOK, nil
}

// Обновление элемента
// @Security TokenAuth
// @Summary Обновление элемента
// @Tags Items
// @Description Обновление элемента
// @ID UpdateItem
// @Accept json
// @Produce json
// @Param item body model.ViewItem true "Элемент"
// @Success      200 {object} Response[model.ViewItem]
// @Failure      400 {object} Response[any]
// @Failure      500 {object} Response[any]
// @Router /private/items [put]
func (o *Server) handleUpdateItem(ctx *fasthttp.RequestCtx) (interface{}, int, error) {
	return o.saveItem(ctx.Request.Body())
}

func (o *Server) saveItem(requestBody []byte) (*model.ViewItem, int, error) {
	item := &model.ViewItem{}
	if err := json.Unmarshal(requestBody, item); err != nil {
		return nil, http.StatusInternalServerError, err
	}

	if item.Params == "" {
		item.Params = "{}"
	}

	if *item.ZoneID == 0 {
		item.ZoneID = nil
	}

	item.Enabled = true

	if err := store.I.Items().SaveItem(item); err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return item, http.StatusOK, nil
}

// Получение итема
// @Security TokenAuth
// @Summary Получение данных итема
// @Tags Items
// @Description Получение данных итема
// @ID GetItem
// @Produce json
// @Param id query int true "ID" default(323)
// @Success      200 {object} Response[model.ViewItem]
// @Failure      400 {object} Response[any]
// @Failure      500 {object} Response[any]
// @Router /private/item [get]
func (o *Server) getItem(ctx *fasthttp.RequestCtx) (interface{}, int, error) {
	id, err := helpers.GetUintParam(ctx, "id")
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	item, err := store.I.Items().GetItem(id)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	return item, http.StatusOK, nil
}

// Изменение итема
// @Security TokenAuth
// @Summary Изменение данных итема
// @Tags Items
// @Description Изменение данных итема
// @ID PatchItem
// @Produce json
// @Param body body model.ViewItem true "Body"
// @Success      200 {object} Response[any]
// @Failure      400 {object} Response[any]
// @Failure      500 {object} Response[any]
// @Router /private/item [patch]
func (o *Server) updateItem(ctx *fasthttp.RequestCtx) (interface{}, int, error) {
	req := &model.ViewItem{}
	if err := json.Unmarshal(ctx.Request.Body(), req); err != nil {
		return nil, http.StatusInternalServerError, err
	}

	if err := store.I.Items().UpdateItem(req); err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return req, http.StatusOK, nil
}

// Удаление итема
// @Security TokenAuth
// @Summary Удаление итема
// @Tags Items
// @Description Удаление итема
// @ID DeleteItem
// @Produce json
// @Param id query int true "ID"
// @Success      200 {object} Response[any]
// @Failure      400 {object} Response[any]
// @Failure      500 {object} Response[any]
// @Router /private/item [delete]
func (o *Server) handleDeleteItem(ctx *fasthttp.RequestCtx) (interface{}, int, error) {
	itemID, err := helpers.GetUintParam(ctx, "id")
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	return o.deleteItem(itemID)
}

func (o *Server) deleteItem(itemID int) (interface{}, int, error) {
	if err := store.I.Items().DeleteItem(itemID); err != nil {
		return nil, http.StatusBadRequest, err
	}

	return nil, http.StatusOK, nil
}

// Получение данных димера
// @Security TokenAuth
// @Summary Получение данных димера
// @Tags Items
// @Description Получение данных димера
// @ID GetDimmer
// @Produce json
// @Param id query int true "ID" default(323)
// @Success      200 {object} Response[model.Dimmer]
// @Failure      400 {object} Response[any]
// @Failure      500 {object} Response[any]
// @Router /private/item/dimer [get]
func (o *Server) getDimmer(ctx *fasthttp.RequestCtx) (interface{}, int, error) {
	id, err := helpers.GetUintParam(ctx, "id")
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	dimmer, err := store.I.Devices().GetDimmer(id)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	return dimmer, http.StatusOK, nil
}

// Получение данных термостата
// @Security TokenAuth
// @Summary Получение данных термостата
// @Tags Items
// @Description Получение данных термостата
// @ID GetThermostat
// @Produce json
// @Param id query int true "ID" default(307)
// @Success      200 {object} Response[model.Sensor]
// @Failure      400 {object} Response[any]
// @Failure      500 {object} Response[any]
// @Router /private/item/thermostat [get]
func (o *Server) getThermostat(ctx *fasthttp.RequestCtx) (interface{}, int, error) {
	id, err := helpers.GetUintParam(ctx, "id")
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	thermostat, err := store.I.Devices().GetSensor(id)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	// Backward compatibility
	if thermostat.Enabled {
		thermostat.Status = "on"
	} else {
		thermostat.Status = "off"
	}

	return thermostat, http.StatusOK, nil
}

// Получение данных датчика
// @Security TokenAuth
// @Summary Получение данных датчика
// @Tags Items
// @Description Получение данных датчика
// @ID GetSensor
// @Produce json
// @Param itemId query int true "ID" default(306)
// @Success      200 {object} Response[model.Sensor]
// @Failure      400 {object} Response[any]
// @Failure      500 {object} Response[any]
// @Router /private/item/sensor [get]
func (o *Server) getSensor(ctx *fasthttp.RequestCtx) (interface{}, int, error) {
	id, err := helpers.GetUintParam(ctx, "itemId")
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	sensor, err := store.I.Devices().GetSensor(id)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	item, err := store.I.Items().GetItem(sensor.ViewItemID)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	sensor.Title = item.Title

	//Берем у объекта текущее значение
	sensorVal, err := objects.LoadObject(sensor.ObjectID, "", "", true)
	if err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(err, "Error of LoadObject: "+strconv.Itoa(sensor.ObjectID))
	}
	sensor.Current, err = sensorVal.GetProps().GetFloatValue("value")
	if err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(err, "Error of get value for sensor: "+strconv.Itoa(sensor.ObjectID))
	}

	minThresholdOfStoredValue, err := sensorVal.GetProps().GetFloatValue("min_threshold")
	if err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(err, "Error of get value for min_threshold: "+strconv.Itoa(sensor.ObjectID))
	}
	if sensor.MinThreshold == nil || *sensor.MinThreshold < minThresholdOfStoredValue {
		sensor.MinThreshold = new(float32)
		*sensor.MinThreshold = minThresholdOfStoredValue
	}
	maxThresholdOfStoredValue, err := sensorVal.GetProps().GetFloatValue("max_threshold")
	if sensor.MaxThreshold == nil || *sensor.MaxThreshold > maxThresholdOfStoredValue {
		sensor.MaxThreshold = new(float32)
		*sensor.MaxThreshold = maxThresholdOfStoredValue
	}

	//У регулятора датчика получаем целевое показание
	sensorChildrens := sensorVal.GetChildren().GetAll()
	for _, children := range sensorChildrens {
		sensor.Target, err = children.GetProps().GetFloatValue("target_sp")
		if err != nil {
			return nil, http.StatusInternalServerError, errors.Wrap(err, "Error of get target for sensor: "+strconv.Itoa(sensor.ObjectID))
		}
	}

	sensor.History, err = store.I.History().GetHistory(id, model.HistoryItemTypeDeviceObject, "")
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	// Backward compatibility
	if sensor.Enabled {
		sensor.Status = "on"
	} else {
		sensor.Status = "off"
	}

	return sensor, http.StatusOK, nil
}

// Получение списка счетчиков
// @Security TokenAuth
// @Summary Получение списка счетчиков
// @Tags Counters
// @Description Получение списка счетчиков
// @ID GetCountersList
// @Produce json
// @Success      200 {object} Response[[]model.Counter]
// @Failure      400 {object} Response[any]
// @Failure      500 {object} Response[any]
// @Router /private/counters-list [get]
func (o *Server) getCountersList(ctx *fasthttp.RequestCtx) (interface{}, int, error) {
	counters, err := store.I.Items().GetCountersList()
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	return counters, http.StatusOK, nil
}

// Получение счетчика
// @Security TokenAuth
// @Summary Получение счетчика
// @Tags Counters
// @Description Получение счетчика
// @ID GetCounter
// @Produce json
// @Param counterId query int true "ID" default(3)
// @Success      200 {object} Response[model.Counter]
// @Failure      400 {object} Response[any]
// @Failure      500 {object} Response[any]
// @Router /private/counter [get]
func (o *Server) getCounter(ctx *fasthttp.RequestCtx) (interface{}, int, error) {
	id, err := helpers.GetUintParam(ctx, "counterId")
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	counter, err := store.I.Items().GetCounter(id)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	counter.History, err = store.I.History().GetHistory(id, model.HistoryItemTypeCounterObject, "")
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	return counter, http.StatusOK, nil
}

type itemChangeRequest struct {
	ItemID int                    `json:"item_id" default:"16"`
	Event  string                 `json:"event" default:"onChange"`
	State  string                 `json:"state" enums:"on,off"`
	Params map[string]interface{} `json:"-,omitempty"`

	// Backward compatibility
	ParamsString string `json:"params,omitempty" default:"{}"`
}

type itemSortRequest struct {
	ZoneId  int   `json:"zone_id"`
	ItemIDs []int `json:"item_ids"`
}

// Добавление события о нажатии/отпускании кнопки
// @Security TokenAuth
// @Summary Добавление события о нажатии/отпускании кнопки
// @Tags Items
// @Description Добавление события о нажатии/отпускании кнопки
// @ID ItemChange
// @Produce json
// @Param body body itemChangeRequest true "Body"
// @Success      200 {object} Response[any]
// @Failure      400 {object} Response[any]
// @Failure      500 {object} Response[any]
// @Router /private/item-change [post]
func (o *Server) itemChange(ctx *fasthttp.RequestCtx) (interface{}, int, error) {
	req := &itemChangeRequest{}
	if err := json.Unmarshal(ctx.Request.Body(), req); err != nil {
		return nil, http.StatusInternalServerError, err
	}

	// Backward compatibility
	params := []byte(req.ParamsString)
	if json.Valid(params) {
		if err := json.Unmarshal(params, &req.Params); err != nil {
			return nil, http.StatusInternalServerError, err
		}
	}

	if err := store.I.Items().ChangeItem(req.ItemID, req.State); err != nil {
		return nil, http.StatusInternalServerError, err
	}

	var msg interfaces.Message
	var err error

	switch req.State {
	case "on":
		msg, err = item.NewOnChangeStateOn(req.ItemID)
	case "off":
		msg, err = item.NewOnChangeStateOff(req.ItemID)
	default:
		msg, err = messages.NewEvent(req.Event, interfaces.TargetTypeItem, req.ItemID)
	}

	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	msg.SetPayload(req.Params)

	if err := g.Msgs.Send(msg); err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return nil, http.StatusOK, nil
}

type setItemsOrderRequest struct {
	ZoneID  int   `json:"zone_id"`
	ItemIDs []int `json:"item_ids"`
}

// Установка порядка отображения элементов
// @Security TokenAuth
// @Summary Установка порядка отображения элементов
// @Tags Items
// @Description Установка порядка отображения элементов
// @ID SetItemsOrder
// @Produce json
// @Param body body setItemsOrderRequest true "Body"
// @Success      200 {object} Response[any]
// @Failure      400 {object} Response[any]
// @Failure      500 {object} Response[any]
// @Router /private/items/order [patch]
func (o *Server) setItemsOrder(ctx *fasthttp.RequestCtx) (interface{}, int, error) {
	req := &setItemsOrderRequest{}

	if err := json.Unmarshal(ctx.Request.Body(), &req); err != nil {
		return nil, http.StatusInternalServerError, err
	}

	if err := store.I.Items().SetOrder(req.ItemIDs, req.ZoneID); err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return nil, http.StatusOK, nil
}
