package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/VladimirDronik/touchon-server/helpers"
	"github.com/fasthttp/router"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpadaptor"
)

func New(name string, cfg map[string]string, ringBuffer fmt.Stringer, logger *logrus.Logger) (*Server, error) {
	if logger == nil {
		return nil, errors.Wrap(errors.New("logger is nil"), "http.New")
	}

	ctx, cancel := context.WithCancel(context.Background())

	o := &Server{
		name: name,
		httpServer: &fasthttp.Server{
			ReadTimeout:          5 * time.Second,
			WriteTimeout:         5 * time.Second,
			NoDefaultContentType: true,
		},
		cfg:        cfg,
		router:     router.New(),
		ringBuffer: ringBuffer,
		logger:     logger,
		ctx:        ctx,
		cancel:     cancel,
	}

	// Обработчик для Swagger'а
	o.router.GET("/swagger/{filepath:*}", fasthttpadaptor.NewFastHTTPHandler(
		httpSwagger.Handler(
			httpSwagger.DeepLinking(true),
			httpSwagger.DocExpansion("none"),
			httpSwagger.DomID("swagger-ui"),
		),
	))

	// Служебные эндпоинты
	o.router.GET("/_/info", JsonHandlerWrapper(o.handleGetInfo))
	o.router.GET("/_/log", o.handleGetLog)

	o.httpServer.Handler = o.RequestWrapper(o.router.Handler)

	return o, nil
}

type Server struct {
	name       string
	httpServer *fasthttp.Server
	router     *router.Router
	ringBuffer fmt.Stringer
	logger     *logrus.Logger
	cfg        map[string]string
	ctx        context.Context
	cancel     context.CancelFunc
}

func (o *Server) AddHandler(method, path string, handler RequestHandler) {
	o.router.Handle(method, path, JsonHandlerWrapper(handler))
}

func (o *Server) GetContext() context.Context {
	return o.ctx
}

func (o *Server) GetServer() *fasthttp.Server {
	return o.httpServer
}

func (o *Server) GetRouter() *router.Router {
	return o.router
}

func (o *Server) GetLogger() *logrus.Logger {
	return o.logger
}

func (o *Server) GetConfig() map[string]string {
	return o.cfg
}

func (o *Server) Start(bindAddr string) error {
	go func() {
		o.logger.Infof("HTTP (%s): %s-сервер запущен", bindAddr, o.name)

		if err := o.httpServer.ListenAndServe(bindAddr); err != nil {
			o.logger.Fatalf("HTTP (%s): %v", o.name, err)
		}

		o.logger.Infof("HTTP (%s): %s-сервер остановлен", bindAddr, o.name)
	}()

	return nil
}

func (o *Server) Shutdown() error {
	o.cancel()

	if err := o.httpServer.Shutdown(); err != nil {
		return errors.Wrap(err, "Shutdown")
	}

	return nil
}

// RequestWrapper добавляет CORS заголовки и content type
func (o *Server) RequestWrapper(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		ref := string(ctx.Request.Header.Peek("Origin"))
		if ref == "" {
			ref = "*"
		}

		ctx.Response.Header.Set("Access-Control-Allow-Origin", ref)
		ctx.Response.Header.Set("Access-Control-Allow-Headers", "api-key,token,content-type")
		ctx.Response.Header.Set("Access-Control-Allow-Methods", "*")
		ctx.Response.Header.Set("Access-Control-Allow-Credentials", "true")

		// Хрому нужен OK на запрос OPTIONS
		if string(ctx.Method()) == http.MethodOptions {
			ctx.SetStatusCode(http.StatusOK)
			return
		}

		if len(ctx.Response.Header.ContentType()) == 0 {
			ctx.Response.Header.SetContentType("application/json; charset=UTF-8")
		}

		if next != nil {
			next(ctx)
		}

		helpers.DumpRequestCtx(o.logger, ctx)
	}
}
