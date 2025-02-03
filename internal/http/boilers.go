package http

import (
	"encoding/json"
	"net/http"

	"github.com/valyala/fasthttp"
	"touchon-server/internal/model"
	"touchon-server/internal/store"
	"touchon-server/lib/helpers"
)

// Получить данные котла
// @Security TokenAuth
// @Summary Получить данные котла
// @Tags Boilers
// @Description Получить данные котла
// @ID GetBoiler
// @Produce json
// @Param id query int true "ID" Format(int) default(1)
// @Success      200 {object} Response[model.Boiler]
// @Failure      400 {object} Response[any]
// @Failure      500 {object} Response[any]
// @Router /private/boiler [get]
func (o *Server) getBoiler(ctx *fasthttp.RequestCtx) (interface{}, int, error) {
	id, err := helpers.GetUintParam(ctx, "id")
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	boiler, err := store.I.Boilers().GetBoiler(id)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	return boiler, http.StatusOK, nil
}

type setBoilerOutlineStatusRequest struct {
	Outline string `json:"outline" enums:"heating,water"`
	Status  string `json:"status" enums:"on,off"`
}

// Включение или выключения контура котла
// @Security TokenAuth
// @Summary Включение или выключения контура котла
// @Tags Boilers
// @Description Включение или выключения контура котла
// @ID SetBoilerOutlineStatus
// @Produce json
// @Param boilerId query int true "ID" Format(int) default(1)
// @Param body body setBoilerOutlineStatusRequest true "Body"
// @Success      200 {object} Response[any]
// @Failure      400 {object} Response[any]
// @Failure      500 {object} Response[any]
// @Router /private/boiler/outline-status [patch]
func (o *Server) setBoilerOutlineStatus(ctx *fasthttp.RequestCtx) (interface{}, int, error) {
	id, err := helpers.GetUintParam(ctx, "boilerId")
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	var data setBoilerOutlineStatusRequest

	if err := json.Unmarshal(ctx.Request.Body(), &data); err != nil {
		return nil, http.StatusBadRequest, err
	}

	if err := store.I.Boilers().SetOutlineStatus(id, data.Outline, data.Status); err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return nil, http.StatusOK, nil
}

type setBoilerHeatingModeRequest struct {
	Mode string `json:"mode" enums:"auto,manual"`
}

// Установка режима отопления котла
// @Security TokenAuth
// @Summary Установка режима отопления котла
// @Tags Boilers
// @Description Установка режима отопления котла
// @ID SetBoilerHeatingMode
// @Produce json
// @Param boilerId query int true "ID" Format(int) default(1)
// @Param body body setBoilerHeatingModeRequest true "Body"
// @Success      200 {object} Response[any]
// @Failure      400 {object} Response[any]
// @Failure      500 {object} Response[any]
// @Router /private/boiler/heating-mode [patch]
func (o *Server) setBoilerHeatingMode(ctx *fasthttp.RequestCtx) (interface{}, int, error) {
	id, err := helpers.GetUintParam(ctx, "boilerId")
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	var data setBoilerHeatingModeRequest

	if err := json.Unmarshal(ctx.Request.Body(), &data); err != nil {
		return nil, http.StatusBadRequest, err
	}

	if err := store.I.Boilers().SetHeatingMode(id, data.Mode); err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return nil, http.StatusOK, nil
}

type setBoilerHeatingTemperatureRequest struct {
	Temperature float32 `json:"Temperature" default:"70"`
}

// Установка температуры отопления для котла
// @Security TokenAuth
// @Summary Установка температуры отопления для котла
// @Tags Boilers
// @Description Установка температуры отопления для котла
// @ID SetBoilerHeatingTemperature
// @Produce json
// @Param boilerId query int true "ID" Format(int) default(1)
// @Param body body setBoilerHeatingTemperatureRequest true "Body"
// @Success      200 {object} Response[any]
// @Failure      400 {object} Response[any]
// @Failure      500 {object} Response[any]
// @Router /private/boiler/heating-temp [patch]
// setBoilerHeatingTemperature ручная установка температуры отопления для котла
func (o *Server) setBoilerHeatingTemperature(ctx *fasthttp.RequestCtx) (interface{}, int, error) {
	id, err := helpers.GetUintParam(ctx, "boilerId")
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	var data setBoilerHeatingTemperatureRequest

	if err := json.Unmarshal(ctx.Request.Body(), &data); err != nil {
		return nil, http.StatusBadRequest, err
	}

	if err := store.I.Boilers().SetHeatingTemperature(id, data.Temperature); err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return nil, http.StatusOK, nil
}

type boilerPresetsUpdateRequest struct {
	BoilerId int                   `json:"id_boiler" default:"1"`
	Presets  []*model.BoilerPreset `json:"presets"` // default:"[{\"temp_out\":0,\"temp_coolant\":45},{\"temp_out\":10.0,\"temp_coolant\":30.0},{\"temp_out\":-10,\"temp_coolant\":50.0},{\"temp_out\":-20,\"temp_coolant\":60.0}]"
}

// Обновить предустановки котла
// @Security TokenAuth
// @Summary Обновить предустановки котла
// @Tags Boilers
// @Description Обновить предустановки котла
// @ID UpdateBoilerPresets
// @Produce json
// @Param body body boilerPresetsUpdateRequest true "Body" default({"id_boiler":1,"presets":[{"temp_out":0,"temp_coolant":45},{"temp_out":10.0,"temp_coolant":30.0},{"temp_out":-10,"temp_coolant":50.0},{"temp_out":-20,"temp_coolant":60.0}]})
// @Success      200 {object} Response[any]
// @Failure      400 {object} Response[any]
// @Failure      500 {object} Response[any]
// @Router /private/boiler/presets [put]
func (o *Server) updateBoilerPresets(ctx *fasthttp.RequestCtx) (interface{}, int, error) {
	var request boilerPresetsUpdateRequest

	// Декодируем JSON из тела запроса
	if err := json.Unmarshal(ctx.Request.Body(), &request); err != nil {
		return nil, http.StatusBadRequest, err
	}

	if err := store.I.Boilers().UpdateBoilerPresets(request.BoilerId, request.Presets); err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return nil, http.StatusOK, nil
}
