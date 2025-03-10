package http

import (
	"encoding/json"
	"net/http"

	"github.com/valyala/fasthttp"
	"touchon-server/internal/helpers"
	"touchon-server/internal/store"
)

// Получение источника света
// @Security TokenAuth
// @Summary Получение источника света
// @Tags Lights
// @Description Получение источника света
// @ID GetLight
// @Produce json
// @Param itemId query int true "ID" default(6)
// @Success      200 {object} Response[model.ViewItem]
// @Failure      400 {object} Response[any]
// @Failure      500 {object} Response[any]
// @Router /private/light [get]
func (o *Server) getLight(ctx *fasthttp.RequestCtx) (interface{}, int, error) {
	id, err := helpers.GetUintParam(ctx, "itemId")
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	light, err := store.I.Lights().GetLight(id)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	return light, http.StatusOK, nil
}

type setLightHSVColorRequest struct {
	Hue        int      `json:"hue" default:"361"`
	Saturation *float32 `json:"saturation" default:"0.3"`
	Brightness *float32 `json:"brightness" default:"0.6"`
}

// Установка цвета в формате HSV
// @Security TokenAuth
// @Summary Установка цвета в формате HSV
// @Tags Lights
// @Description Установка цвета в формате HSV
// @ID SetLightHSVColor
// @Produce json
// @Param itemId query int true "ID" default(6)
// @Param body body setLightHSVColorRequest true "Body"
// @Success      200 {object} Response[any]
// @Failure      400 {object} Response[any]
// @Failure      500 {object} Response[any]
// @Router /private/light/hsv [patch]
func (o *Server) setLightHSVColor(ctx *fasthttp.RequestCtx) (interface{}, int, error) {
	id, err := helpers.GetUintParam(ctx, "itemId")
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	var requestData setLightHSVColorRequest

	if err := json.Unmarshal(ctx.Request.Body(), &requestData); err != nil {
		return nil, http.StatusBadRequest, err
	}

	if err = store.I.Lights().SetLightHSVColor(id, requestData.Hue, *requestData.Saturation, *requestData.Brightness); err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return nil, http.StatusOK, nil
}

type setLightCCTColorRequest struct {
	Cct *int `json:"cct" default:"7001"`
}

// Установка цвета в формате CCT
// @Security TokenAuth
// @Summary Установка цвета в формате CCT
// @Tags Lights
// @Description Установка цвета в формате CCT
// @ID SetLightCCTColor
// @Produce json
// @Param itemId query int true "ID" default(6)
// @Param body body setLightCCTColorRequest true "Body"
// @Success      200 {object} Response[any]
// @Failure      400 {object} Response[any]
// @Failure      500 {object} Response[any]
// @Router /private/light/cct [patch]
func (o *Server) setLightCCTColor(ctx *fasthttp.RequestCtx) (interface{}, int, error) {
	id, err := helpers.GetUintParam(ctx, "itemId")
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	var requestData setLightCCTColorRequest

	if err := json.Unmarshal(ctx.Request.Body(), &requestData); err != nil {
		return nil, http.StatusBadRequest, err
	}

	if err = store.I.Lights().SetLightCCTColor(id, *requestData.Cct); err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return nil, http.StatusOK, nil
}

type setLightBrightnessRequest struct {
	Brightness *float32 `json:"brightness" default:"0.6"`
}

// Установка яркости
// @Security TokenAuth
// @Summary Установка яркости
// @Tags Lights
// @Description Установка яркости
// @ID SetLightBrightness
// @Produce json
// @Param itemId query int true "ID" default(6)
// @Param body body setLightBrightnessRequest true "Body"
// @Success      200 {object} Response[any]
// @Failure      400 {object} Response[any]
// @Failure      500 {object} Response[any]
// @Router /private/light/brightness [patch]
func (o *Server) setLightBrightness(ctx *fasthttp.RequestCtx) (interface{}, int, error) {
	id, err := helpers.GetUintParam(ctx, "itemId")
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	var requestData setLightBrightnessRequest

	if err := json.Unmarshal(ctx.Request.Body(), &requestData); err != nil {
		return nil, http.StatusBadRequest, err
	}

	if err = store.I.Lights().SetBrightness(id, *requestData.Brightness); err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return nil, http.StatusOK, nil
}
