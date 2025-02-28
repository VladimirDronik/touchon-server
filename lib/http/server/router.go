package server

import (
	"bytes"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/valyala/fasthttp"
	"touchon-server/lib/interfaces"
)

type TextPlain = string

// Meta Метаинформация о запросе/ответе
type Meta struct {
	Duration      float64 `json:"duration"`       // Длительность запроса
	ContentLength string  `json:"content_length"` // Размер данных ответа
}

// Response Ответ сервиса
type Response[T any] struct {
	Meta    Meta   `json:"meta"`               // Метаинформация о запросе/ответе
	Success bool   `json:"success"`            // Подтверждение отсутствия ошибок в выводе
	Data    T      `json:"response,omitempty"` // Полезная нагрузка, зависит от запроса
	Error   string `json:"error,omitempty"`    // Описание возвращенной ошибки
}

// JsonHandlerWrapper ответ в формате JSON оборачивает в единый формат и добавляет метаданные.
func JsonHandlerWrapper(f interfaces.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		var r Response[any]
		const magic = "CoNtEnTLeNgTh"

		start := time.Now()
		data, status, err := f(ctx)
		r.Meta.Duration = float64(int(time.Since(start).Seconds()*1000000)) / 1000000
		r.Meta.ContentLength = magic
		ctx.Response.SetStatusCode(status)

		switch {
		case err != nil:
			r.Error = err.Error()
		case data != nil:
			r.Data = data
		}

		r.Success = err == nil

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
