package http

import (
	"net/http"
	"strings"

	"github.com/valyala/fasthttp"
	"touchon-server/internal/model"
	"touchon-server/internal/objects"
	"touchon-server/internal/store"
	"touchon-server/lib/helpers"
	"touchon-server/lib/info"
)

// Получить информацию о сервисе
// @Summary Получить информацию о сервисе
// @Tags Service
// @Description Получить информацию о сервисе
// @ID ServiceInfo
// @Produce json
// @Success      200 {object} http.Response[info.Info]
// @Failure      500 {object} http.Response[any]
// @Router /_/info [get]
func (o *Server) handleGetInfo(ctx *fasthttp.RequestCtx) (interface{}, int, error) {
	nfo, err := info.GetInfo()
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return nfo, http.StatusOK, nil
}

// Получить логи
// @Summary Получить логи
// @Tags Service
// @Description Получить логи
// @ID ServiceLog
// @Produce text/plain
// @Success      200
// @Router /_/log [get]
func (o *Server) handleGetLog(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.SetContentType("text/plain; charset=UTF-8")
	if o.ringBuffer != nil {
		_, _ = ctx.WriteString(o.ringBuffer.String())
	}
}

type SensorValues struct {
	ID     int                `json:"id"`
	Type   string             `json:"type"`
	Values map[string]float32 `json:"values,omitempty"`
	Error  string             `json:"error,omitempty"`
}

// Получение значений датчиков
// @Summary Получение значений датчиков
// @Tags Service
// @Description Получение значений датчиков
// @ID ServiceSensors
// @Produce json
// @Success      200 {object} http.Response[[]SensorValues]
// @Failure      400 {object} http.Response[any]
// @Failure      500 {object} http.Response[any]
// @Router /_/sensors [get]
func (o *Server) handleGetSensors(ctx *fasthttp.RequestCtx) (interface{}, int, error) {
	filters := map[string]interface{}{"category": string(model.CategorySensor)}
	tags := strings.Split(helpers.GetParam(ctx, "tags"), ",")

	tagsMap := make(map[string]bool, len(tags))
	for _, tag := range tags {
		tag = strings.TrimSpace(tag)
		if tag != "" {
			tagsMap[tag] = true
		}
	}

	tags = tags[:0]
	for tag := range tagsMap {
		tags = append(tags, tag)
	}

	rows, err := store.I.ObjectRepository().GetObjects(filters, tags, 0, 1000, model.ChildTypeNobody)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	r := make([]*SensorValues, 0, len(rows))

	for _, row := range rows {
		objModel, err := objects.LoadObject(row.ID, "", "", model.ChildTypeInternal)
		if err != nil {
			return nil, http.StatusInternalServerError, err
		}

		rItem := &SensorValues{
			ID:     objModel.GetID(),
			Type:   objModel.GetType(),
			Values: make(map[string]float32, 5),
		}

		m, err := objModel.GetMethods().Get("check")
		if err != nil {
			rItem.Error = err.Error()
			r = append(r, rItem)
			continue
		}

		if _, err = m.Func(nil); err != nil {
			rItem.Values = nil
			rItem.Error = err.Error()
			r = append(r, rItem)
			continue
		}

		for _, valueObj := range objModel.GetChildren().GetAll() {
			if valueObj.GetCategory() != model.CategorySensorValue {
				continue
			}

			v, err := valueObj.GetProps().GetFloatValue("value")
			if err != nil {
				rItem.Values = nil
				rItem.Error = err.Error()
				r = append(r, rItem)
				break
			}

			rItem.Values[valueObj.GetType()] = v
		}

		r = append(r, rItem)
	}

	return r, http.StatusOK, nil
}
