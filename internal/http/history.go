package http

import (
	"net/http"
	"strconv"

	"github.com/VladimirDronik/touchon-server/helpers"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
	"translator/internal/model"
)

// Возвращает историю изменений значений
// @Security TokenAuth
// @Summary Возвращает историю изменений значений
// @Tags History
// @Description Возвращает историю изменений значений
// @ID GetObjectHistory
// @Produce json
// @Param itemId   query int true  "ID"     default(10)
// @Param itemType query string true  "Type"   Enums(device,counter)
// @Param filter   query string false "Filter" Enums(day,week,month,year)
// @Success      200 {object} Response[model.HistoryPoints]
// @Failure      400 {object} Response[any]
// @Failure      500 {object} Response[any]
// @Router /private/history [get]
func (o *Server) getObjectHistory(ctx *fasthttp.RequestCtx) (interface{}, int, error) {
	id, err := helpers.GetUintParam(ctx, "itemId")
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	itemType := model.HistoryItemType(helpers.GetParam(ctx, "itemType"))
	filter := model.HistoryFilter(helpers.GetParam(ctx, "filter"))

	points, err := o.store.History().GetHistory(id, itemType, filter)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return points, http.StatusOK, nil
}

// Генерирует тестовые данные для истории (в продакшене - 404)
// @Security TokenAuth
// @Summary Генерирует тестовые данные для истории (в продакшене - 404)
// @Tags History
// @Description Генерирует тестовые данные для истории (в продакшене - 404)
// @ID GenerateHistory
// @Produce json
// @Param itemId   query int true  "ID"     default(10)
// @Param itemType query string true  "Type"   Enums(device,counter)
// @Param filter   query string true "Filter" Enums(day,week,month,year)
// @Param startDate query string true "Start date" Format(date-time) Default(2006-01-02 15:04)
// @Param endDate query string true "End date" Format(date-time) Default(2006-01-02 15:04)
// @Param min query number true "Min" Format(float) Default(0.1)
// @Param max query number true "Max" Format(float) Default(1.0)
// @Success      200 {object} Response[any]
// @Failure      400 {object} Response[any]
// @Failure      500 {object} Response[any]
// @Router /private/generate-history [get]
func (o *Server) generateHistory(ctx *fasthttp.RequestCtx) (interface{}, int, error) {
	if o.GetLogger().Level < logrus.DebugLevel {
		return nil, http.StatusNotFound, errors.New("not found")
	}

	id, err := helpers.GetUintParam(ctx, "itemId")
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	itemType := model.HistoryItemType(helpers.GetParam(ctx, "itemType"))
	startDate := helpers.GetParam(ctx, "startDate")
	endDate := helpers.GetParam(ctx, "endDate")

	minParam := helpers.GetParam(ctx, "min")
	maxParam := helpers.GetParam(ctx, "max")

	if minParam == "" || maxParam == "" {
		return nil, http.StatusBadRequest, errors.New("invalid parameter")
	}

	minValue, err := strconv.ParseFloat(minParam, 64)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	maxValue, err := strconv.ParseFloat(maxParam, 64)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	filter := model.HistoryFilter(helpers.GetParam(ctx, "filter"))
	if filter == "" {
		return nil, http.StatusBadRequest, errors.New("filter is empty")
	}

	err = o.store.History().GenerateHistory(id, itemType, filter, startDate, endDate, float32(minValue), float32(maxValue))
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	return nil, http.StatusOK, nil
}
