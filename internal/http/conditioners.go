package http

import (
	"encoding/json"
	"net/http"
	"touchon-server/internal/model"
	memStore "touchon-server/internal/store/memstore"

	"github.com/valyala/fasthttp"
	"touchon-server/internal/helpers"
	"touchon-server/internal/store"
)

// Получение данных кондиционера
// @Security TokenAuth
// @Summary Получение данных кондиционера
// @Tags Conditioners
// @Description Получение данных кондиционера
// @ID GetConditioner
// @Produce json
// @Param itemId query int true "ID" Format(int) default(16)
// @Success      200 {object} Response[model.CobditionerParams]
// @Failure      400 {object} Response[any]
// @Failure      500 {object} Response[any]
// @Router /private/item/conditioner [get]
// getConditioner
func (o *Server) getConditioner(ctx *fasthttp.RequestCtx) (interface{}, int, error) {
	itemID, err := helpers.GetUintParam(ctx, "itemId")
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	conditioner, err := store.I.Conditioners().GetConditioner(itemID)

	//conditioner, err := store.I.Conditioners().GetConditioner(itemId)
	//if err != nil {
	//	return nil, http.StatusBadRequest, err
	//}

	//conditioner.History, err = store.I.History().GetHistory(id, model.HistoryItemTypeDeviceObject, "")
	//if err != nil {
	//	return nil, http.StatusInternalServerError, err
	//}

	return conditioner, http.StatusOK, nil
}

type setConditionerTemperatureRequest struct {
	Value uint `json:"value" default:"27"`
}

// Установка температуры для кондиционера
// @Security TokenAuth
// @Summary Установка температуры для кондиционера
// @Tags Conditioners
// @Description Установка температуры для кондиционера
// @ID SetConditionerTemperature
// @Produce json
// @Param itemId query int true "ID" Format(int) default(16)
// @Param body body setConditionerTemperatureRequest true "Body"
// @Success      200 {object} Response[any]
// @Failure      400 {object} Response[any]
// @Failure      500 {object} Response[any]
// @Router /private/conditioner/temp [patch]
func (o *Server) setConditionerTemperature(ctx *fasthttp.RequestCtx) (interface{}, int, error) {
	id, err := helpers.GetUintParam(ctx, "itemId")
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	var data setConditionerTemperatureRequest

	if err := json.Unmarshal(ctx.Request.Body(), &data); err != nil {
		return nil, http.StatusBadRequest, err
	}

	if err := store.I.Conditioners().SetConditionerTemperature(id, float32(data.Value)); err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return nil, http.StatusOK, nil
}

type setConditionerModeRequest struct {
	Mode  string `json:"mode" enums:"eco_mode,silent_mode,turbo_mode,sleep_mode"`
	Value bool   `json:"value" default:"true"`
}

// Установка режима для кондиционера
// @Security TokenAuth
// @Summary Установка режима для кондиционера
// @Tags Conditioners
// @Description Установка режима для кондиционера
// @ID SetConditionerMode
// @Produce json
// @Param itemId query int true "ID" Format(int) default(16)
// @Param body body setConditionerModeRequest true "Body"
// @Success      200 {object} Response[any]
// @Failure      400 {object} Response[any]
// @Failure      500 {object} Response[any]
// @Router /private/conditioner/mode [patch]
func (o *Server) setConditionerMode(ctx *fasthttp.RequestCtx) (interface{}, int, error) {
	id, err := helpers.GetUintParam(ctx, "itemId")
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	var data setConditionerModeRequest

	if err := json.Unmarshal(ctx.Request.Body(), &data); err != nil {
		return nil, http.StatusBadRequest, err
	}

	if err := store.I.Conditioners().SetConditionerMode(id, data.Mode, data.Value); err != nil {
		return nil, http.StatusBadRequest, err
	}

	return nil, http.StatusOK, nil
}

type setConditionerOperatingModeRequest struct {
	Mode string `json:"mode" enums:"auto,cooling,heating,dehumidification,ventilation"`
}

// Установка режима работы кондиционера
// @Security TokenAuth
// @Summary Установка режима работы кондиционера
// @Tags Conditioners
// @Description Установка режима работы кондиционера
// @ID SetConditionerOperatingMode
// @Produce json
// @Param itemId query int true "ID" Format(int) default(16)
// @Param body body setConditionerOperatingModeRequest true "Body"
// @Success      200 {object} Response[any]
// @Failure      400 {object} Response[any]
// @Failure      500 {object} Response[any]
// @Router /private/conditioner/operating-mode [patch]
func (o *Server) setConditionerOperatingMode(ctx *fasthttp.RequestCtx) (interface{}, int, error) {
	id, err := helpers.GetUintParam(ctx, "itemId")
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	var data setConditionerOperatingModeRequest

	if err := json.Unmarshal(ctx.Request.Body(), &data); err != nil {
		return nil, http.StatusBadRequest, err
	}

	if err := store.I.Conditioners().SetConditionerOperatingMode(id, data.Mode); err != nil {
		return nil, http.StatusBadRequest, err
	}

	return nil, http.StatusOK, nil
}

type setConditionerDirectionRequest struct {
	Plane     string `json:"lamely" enums:"vertical,horizontal"`
	Direction string `json:"direction" enums:"auto,swing,first_position,second_position,third_position,fourth_position,fifth_position"`
}

// Установка направления ламелей
// @Security TokenAuth
// @Summary Установка направления ламелей
// @Tags Conditioners
// @Description Установка направления ламелей
// @ID SetConditionerDirection
// @Produce json
// @Param itemId query int true "ID" Format(int) default(16)
// @Param body body setConditionerDirectionRequest true "Body"
// @Success      200 {object} Response[any]
// @Failure      400 {object} Response[any]
// @Failure      500 {object} Response[any]
// @Router /private/conditioner/air-direction [patch]
func (o *Server) setConditionerDirection(ctx *fasthttp.RequestCtx) (interface{}, int, error) {
	id, err := helpers.GetUintParam(ctx, "itemId")
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	var data setConditionerDirectionRequest

	if err := json.Unmarshal(ctx.Request.Body(), &data); err != nil {
		return nil, http.StatusBadRequest, err
	}

	if err := store.I.Conditioners().SetConditionerDirection(id, data.Plane, data.Direction); err != nil {
		return nil, http.StatusBadRequest, err
	}

	return nil, http.StatusOK, nil
}

type setConditionerFanSpeedRequest struct {
	Speed string `json:"speed" enums:"auto,first,second,third,fourth,fifth"`
}

// Установка скорости вентилятора
// @Security TokenAuth
// @Summary Установка скорости вентилятора
// @Tags Conditioners
// @Description Установка скорости вентилятора
// @ID SetConditionerFanSpeed
// @Produce json
// @Param itemId query int true "ID" Format(int) default(16)
// @Param body body setConditionerFanSpeedRequest true "Body"
// @Success      200 {object} Response[any]
// @Failure      400 {object} Response[any]
// @Failure      500 {object} Response[any]
// @Router /private/conditioner/fan-speed [patch]
func (o *Server) setConditionerFanSpeed(ctx *fasthttp.RequestCtx) (interface{}, int, error) {
	id, err := helpers.GetUintParam(ctx, "itemId")
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	var data setConditionerFanSpeedRequest

	if err := json.Unmarshal(ctx.Request.Body(), &data); err != nil {
		return nil, http.StatusBadRequest, err
	}

	if err := store.I.Conditioners().SetConditionerFanSpeed(id, data.Speed); err != nil {
		return nil, http.StatusBadRequest, err
	}

	return nil, http.StatusOK, nil
}

type setConditionerExtraModeRequest struct {
	Mode  string `json:"mode" enums:"ionisation,self_cleaning,anti_mold,sound,on_duty_heating,soft_top"`
	Value bool   `json:"value" default:"true"`
}

// Установка дополнительного режима
// @Security TokenAuth
// @Summary Установка дополнительного режима
// @Tags Conditioners
// @Description Установка дополнительного режима
// @ID SetConditionerExtraMode
// @Produce json
// @Param itemId query int true "ID" Format(int) default(16)
// @Param body body setConditionerExtraModeRequest true "Body"
// @Success      200 {object} Response[any]
// @Failure      400 {object} Response[any]
// @Failure      500 {object} Response[any]
// @Router /private/conditioner/extra-mode [patch]
func (o *Server) setConditionerExtraMode(ctx *fasthttp.RequestCtx) (interface{}, int, error) {
	id, err := helpers.GetUintParam(ctx, "itemId")
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	var data setConditionerExtraModeRequest

	if err := json.Unmarshal(ctx.Request.Body(), &data); err != nil {
		return nil, http.StatusBadRequest, err
	}

	if err := store.I.Conditioners().SetConditionerExtraMode(id, data.Mode, data.Value); err != nil {
		return nil, http.StatusBadRequest, err
	}

	return nil, http.StatusOK, nil
}

// Создание итема кондиционера
// @Security TokenAuth
// @Summary Создание итема кондиционера
// @Tags Conditioners
// @Description Создание итема кондиционера
// @ID CreateConditioner
// @Accept json
// @Produce json
// @Param item body model.ConditionerItem true "Кондиционер"
// @Success      200 {object} Response[model.ConditionerItem]
// @Failure      400 {object} Response[any]
// @Failure      500 {object} Response[any]
// @Router /private/item/conditioner [post]
func (o *Server) handleCreateConditioner(ctx *fasthttp.RequestCtx) (interface{}, int, error) {
	conditioner := &model.ConditionerItem{}

	if err := json.Unmarshal(ctx.Request.Body(), conditioner); err != nil {
		return nil, http.StatusInternalServerError, err
	}

	condObj, err := memStore.I.GetObject(conditioner.ObjectID)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	item := &model.ViewItem{
		Type:    "sensor",
		Enabled: true,
		ZoneID:  condObj.GetZoneID(),
		Title:   conditioner.Title,
		Icon:    conditioner.Icon,
		Params:  "{}",
		Auth:    "",
		Sort:    0,
	}

	if err := store.I.Items().SaveItem(item); err != nil {
		return nil, http.StatusInternalServerError, err
	}
	conditioner.ViewItemID = item.ID

	event := &model.TrEvent{
		EventName:  "on_change_state",
		TargetType: "object",
		TargetID:   condObj.GetID(),
		Enabled:    true,
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
		Name:       "set_state",
		Args:       "{\"param\":\"state\"}",
		Enabled:    true,
	}

	_, err = store.I.Events().AddEventAction(eventAction)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return conditioner, http.StatusOK, store.I.Devices().SaveConditioner(conditioner)
}
