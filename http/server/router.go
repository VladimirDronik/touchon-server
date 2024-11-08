package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/VladimirDronik/touchon-server/info"
	"github.com/valyala/fasthttp"
)

type RequestHandler func(ctx *fasthttp.RequestCtx) (body interface{}, status int, err error)

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

// Meta Метаинформация о запросе/ответе
type Meta struct {
	Duration      float64 `json:"duration"`       // Длительность запроса
	ContentLength string  `json:"content_length"` // Размер данных ответа
}

// Response Ответ сервиса
type Response[T any] struct {
	Meta  Meta   `json:"meta"`            // Метаинформация о запросе/ответе
	Data  T      `json:"data,omitempty"`  // Полезная нагрузка, зависит от запроса
	Error string `json:"error,omitempty"` // Описание возвращенной ошибки
}

// JsonHandlerWrapper ответ в формате JSON оборачивает в единый формат и добавляет метаданные.
func JsonHandlerWrapper(f RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		var r Response[any]
		const magic = "CoNtEnTLeNgTh"

		start := time.Now()
		data, status, err := f(ctx)
		r.Meta.Duration = float64(int(time.Since(start).Seconds()*1000)) / 1000
		r.Meta.ContentLength = magic
		ctx.Response.SetStatusCode(status)

		switch {
		case err != nil:
			r.Error = err.Error()
		case data != nil:
			r.Data = data
		}

		var buf bytes.Buffer

		enc := json.NewEncoder(&buf)
		enc.SetIndent("", "  ")
		if err := enc.Encode(r); err != nil {
			buf.Reset()
			r.Data = nil
			r.Error = err.Error()
			_ = enc.Encode(r)
		}

		// Выставляем размер ответа
		contLength := buf.Len() - len(magic)
		body := strings.Replace(buf.String(), magic, strconv.Itoa(contLength/1024)+"K", 1)
		_, _ = ctx.WriteString(body)
	}
}
