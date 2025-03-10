package http

import (
	"net/http"
	"strings"
	"sync"

	"github.com/valyala/fasthttp"
	"touchon-server/internal/g"
	"touchon-server/internal/helpers"
	"touchon-server/internal/model"
	"touchon-server/internal/store"
	"touchon-server/internal/store/memstore"
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
	ID     int                    `json:"id"`
	Type   string                 `json:"type"`
	Values map[string]interface{} `json:"values,omitempty"`
	Error  string                 `json:"error,omitempty"`
}

// Получение значений датчиков
// @Summary Получение значений датчиков
// @Tags Service
// @Description Получение значений датчиков
// @ID ServiceSensors
// @Produce json
// @Param tags query string false "Тэги"
// @Success      200 {object} http.Response[[]SensorValues]
// @Failure      400 {object} http.Response[any]
// @Failure      500 {object} http.Response[any]
// @Router /_/sensors [get]
func (o *Server) handleGetSensors(ctx *fasthttp.RequestCtx) (interface{}, int, error) {
	filters := map[string]interface{}{"category": model.CategorySensor}
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

	rows, err := store.I.ObjectRepository().GetObjects(filters, tags, 0, 1000)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	r := make([]*SensorValues, 0, len(rows))

	for _, row := range rows {
		objModel, err := memstore.I.GetObject(row.ID)
		if err != nil {
			return nil, http.StatusInternalServerError, err
		}

		rItem := &SensorValues{
			ID:   objModel.GetID(),
			Type: objModel.GetType(),
		}
		r = append(r, rItem)

		m, err := objModel.GetMethods().Get("check")
		if err != nil {
			rItem.Error = err.Error()
			continue
		}

		if _, err = m.Func(nil); err != nil {
			rItem.Error = err.Error()
			continue
		}

		state, err := objModel.GetState()
		if err != nil {
			rItem.Error = err.Error()
			continue
		}

		rItem.Values = state.GetPayload()
	}

	return r, http.StatusOK, nil
}

var (
	authOldTokenMU sync.Mutex
	authOldToken   = "disable_auth"
)

// Включение/отключение авторизации (только в debug mode)
// @Summary Включение/отключение авторизации (только в debug mode)
// @Tags Service
// @Description Включение/отключение авторизации (только в debug mode)
// @ID ServiceAuth
// @Produce json
// @Success      200 {object} http.Response[string]
// @Failure      400 {object} http.Response[any]
// @Failure      500 {object} http.Response[any]
// @Router /_/switch_auth [post]
func (o *Server) handleSwitchAuth(ctx *fasthttp.RequestCtx) (interface{}, int, error) {
	authOldTokenMU.Lock()
	defer authOldTokenMU.Unlock()

	authOldToken, g.Config["token_secret"] = g.Config["token_secret"], authOldToken

	if g.Config["token_secret"] == "disable_auth" {
		return "Авторизация отключена", http.StatusOK, nil
	} else {
		return "Авторизация включена", http.StatusOK, nil
	}
}
