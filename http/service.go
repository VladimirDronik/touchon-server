package http

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/fasthttp/router"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpadaptor"
)

func New(ringBuffer fmt.Stringer, logger *logrus.Logger) (*Service, error) {
	if logger == nil {
		return nil, errors.Wrap(errors.New("logger is nil"), "http.New")
	}

	o := &Service{
		httpServer: &fasthttp.Server{
			ReadTimeout:          5 * time.Second,
			WriteTimeout:         5 * time.Second,
			NoDefaultContentType: true,
		},
		router:     router.New(),
		ringBuffer: ringBuffer,
		logger:     logger,
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
	o.router.GET("/_/info", handlerWrapper(o.handleGetInfo))
	o.router.GET("/_/log", handlerWrapper(o.handleGetLog))

	o.httpServer.Handler = o.requestWrapper

	return o, nil
}

type Service struct {
	httpServer *fasthttp.Server
	router     *router.Router
	ringBuffer fmt.Stringer
	logger     *logrus.Logger
	cfg        interface{} // for GET /_/info
}

func (o *Service) SetConfig(cfg interface{}) {
	o.cfg = cfg
}

func (o *Service) AddHandler(method, path string, handler RequestHandler) {
	o.router.Handle(method, path, handlerWrapper(handler))
}

func (o *Service) GetLogger() *logrus.Logger {
	return o.logger
}

func (o *Service) Start(bindAddr string) error {
	go func() {
		o.logger.Info("HTTP: сервис запустился")
		if err := o.httpServer.ListenAndServe(bindAddr); err != nil {
			log.Fatal("HTTP:", err)
		}
	}()

	return nil
}

func (o *Service) Shutdown() error {
	if err := o.httpServer.Shutdown(); err != nil {
		return errors.Wrap(err, "Shutdown")
	}

	return nil
}

func (o *Service) requestWrapper(ctx *fasthttp.RequestCtx) {
	ref := string(ctx.Request.Header.Peek("Origin"))
	if ref == "" {
		ref = "*"
	}

	ctx.Response.Header.Set("Access-Control-Allow-Origin", ref)
	ctx.Response.Header.Set("Access-Control-Allow-Headers", "api-key,token")
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

	o.router.Handler(ctx)
}
