package http

import (
	"encoding/json"
	"net/http"

	"github.com/VladimirDronik/touchon-server/events/item"
	"github.com/VladimirDronik/touchon-server/helpers"
	"github.com/VladimirDronik/touchon-server/mqtt/messages"
	"github.com/valyala/fasthttp"
	"translator/internal/model"
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
	item := &model.ViewItem{}
	event := &model.Event{}

	if err := json.Unmarshal(ctx.Request.Body(), sensor); err != nil {
		return nil, http.StatusInternalServerError, err
	}

	item.Type = "sensor"
	item.Enabled = true
	item.ZoneID = sensor.ZoneID
	itemID, err := o.store.Items().SaveItem(item)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	event.EventName = "object.sensor.on_check"
	event.TargetType = "object"
	event.TargetID = sensor.ObjectID
	event.Value = sensor.Type
	event.ItemID = itemID
	sensor.ViewItemID = itemID
	sensor.Enabled = true

	o.store.Events().AddEvent(event)

	return sensor, http.StatusOK, o.store.Devices().SaveSensor(sensor)
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

	o.store.Events().DeleteEvent(itemID)
	o.store.Devices().DeleteSensor(itemID)

	return o.deleteItem(itemID)
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

	item.Enabled = true

	if _, err := o.store.Items().SaveItem(item); err != nil {
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

	item, err := o.store.Items().GetItem(id)
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

	if err := o.store.Items().UpdateItem(req); err != nil {
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
	if err := o.store.Items().DeleteItem(itemID); err != nil {
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

	dimmer, err := o.store.Devices().GetDimmer(id)
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

	thermostat, err := o.store.Devices().GetSensor(id)
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
// @Tags Sensors
// @Description Получение данных датчика
// @ID GetSensor
// @Produce json
// @Param itemId query int true "ID" default(306)
// @Success      200 {object} Response[model.Sensor]
// @Failure      400 {object} Response[any]
// @Failure      500 {object} Response[any]
// @Router /private/sensor [get]
func (o *Server) getSensor(ctx *fasthttp.RequestCtx) (interface{}, int, error) {
	id, err := helpers.GetUintParam(ctx, "itemId")
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	sensor, err := o.store.Devices().GetSensor(id)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	sensor.History, err = o.store.History().GetHistory(id, model.HistoryItemTypeDeviceObject, "")
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
	counters, err := o.store.Items().GetCountersList()
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

	counter, err := o.store.Items().GetCounter(id)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	counter.History, err = o.store.History().GetHistory(id, model.HistoryItemTypeCounterObject, "")
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

	if err := o.store.Items().ChangeItem(req.ItemID, req.State); err != nil {
		return nil, http.StatusInternalServerError, err
	}

	var msg messages.Message
	var err error

	switch req.State {
	case "on":
		msg, err = item.NewOnChangeStateOnMessage("translator/item/event", req.ItemID)
	case "off":
		msg, err = item.NewOnChangeStateOffMessage("translator/item/event", req.ItemID)
	default:
		msg, err = messages.NewEvent(req.Event, messages.TargetTypeItem, req.ItemID, nil)
	}

	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	msg.SetTopic("translator/item/event")
	msg.SetPayload(req.Params)

	if err := o.mqttClient.Send(msg); err != nil {
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

	if err := o.store.Items().SetOrder(req.ItemIDs, req.ZoneID); err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return nil, http.StatusOK, nil
}
