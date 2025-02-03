package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/VladimirDronik/touchon-server/helpers"
	"github.com/valyala/fasthttp"
)

// Возвращает объект шторы
// @Security TokenAuth
// @Summary Возвращает объект шторы
// @Tags Curtains
// @Description Возвращает объект шторы
// @ID GetCurtain
// @Produce json
// @Param itemId query int true "ID" Format(int) default(10)
// @Success      200 {object} Response[model.ViewItem]
// @Failure      400 {object} Response[any]
// @Failure      500 {object} Response[any]
// @Router /private/curtain [get]
func (o *Server) getCurtain(ctx *fasthttp.RequestCtx) (interface{}, int, error) {
	id, err := helpers.GetUintParam(ctx, "itemId")
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	curtain, err := o.store.Curtains().GetCurtain(id)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	return curtain, http.StatusOK, nil
}

type setCurtainOpenPercentRequest struct {
	Value float32 `json:"value" default:"50"`
}

// Установка процента открытия шторы
// @Security TokenAuth
// @Summary Установка процента открытия шторы
// @Tags Curtains
// @Description Установка процента открытия шторы
// @ID SetCurtainOpenPercent
// @Produce json
// @Param itemId query int true "ID" Format(int) default(10)
// @Param body body setCurtainOpenPercentRequest true "Body"
// @Success      200 {object} Response[any]
// @Failure      400 {object} Response[any]
// @Failure      500 {object} Response[any]
// @Router /private/curtain/open-percent [patch]
// setCurtainOpenPercentHandler обработчик для установки
func (o *Server) setCurtainOpenPercent(ctx *fasthttp.RequestCtx) (interface{}, int, error) {
	id, err := helpers.GetUintParam(ctx, "itemId")
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	var requestData setCurtainOpenPercentRequest

	if err := json.Unmarshal(ctx.Request.Body(), &requestData); err != nil {
		return nil, http.StatusBadRequest, err
	}

	if requestData.Value < 0 || requestData.Value > 100 {
		return nil, http.StatusBadRequest, errors.New("value is invalid")
	}

	if err = o.store.Curtains().SetCurtainOpenPercent(id, requestData.Value); err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return nil, http.StatusOK, nil
}
