package interfaces

import (
	"context"

	"github.com/fasthttp/router"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

type RequestHandler func(ctx *fasthttp.RequestCtx) (body interface{}, status int, err error)

type HttpServer interface {
	GetGzipResponse() bool
	SetGzipResponseIfPossible(v bool)
	AddHandler(method, path string, handler RequestHandler)
	GetContext() context.Context
	GetServer() *fasthttp.Server
	GetRouter() *router.Router
	GetLogger() *logrus.Logger
	GetConfig() map[string]string
	Start(bindAddr string) error
	Shutdown() error
	SetRequestID(ctx *fasthttp.RequestCtx)
	RequestWrapper(next fasthttp.RequestHandler) fasthttp.RequestHandler

	CreateCronTask(task *CronTask) error
	CreateEventAction(targetType string, targetID int, eventName string, act *EventAction) error
	DeleteObject(objectID int) (interface{}, int, error)
}

type WSServer interface {
	Send(event string, data interface{})
	Start(bindAddr string) error
	Shutdown() error
}
