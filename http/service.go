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

func New(cfg map[string]string, ringBuffer fmt.Stringer, logger *logrus.Logger) (*Service, error) {
	if logger == nil {
		return nil, errors.Wrap(errors.New("logger is nil"), "http.New")
	}

	o := &Service{
		httpServer: &fasthttp.Server{
			ReadTimeout:          5 * time.Second,
			WriteTimeout:         5 * time.Second,
			NoDefaultContentType: true,
		},
		cfg:        cfg,
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
	o.router.GET("/_/info", JsonHandlerWrapper(o.handleGetInfo))
	o.router.GET("/_/log", o.handleGetLog)

	o.httpServer.Handler = RequestWrapper(o.router.Handler)

	return o, nil
}

type Service struct {
	httpServer *fasthttp.Server
	router     *router.Router
	ringBuffer fmt.Stringer
	logger     *logrus.Logger
	cfg        map[string]string
}

func (o *Service) AddHandler(method, path string, handler RequestHandler) {
	o.router.Handle(method, path, JsonHandlerWrapper(handler))
}

func (o *Service) GetServer() *fasthttp.Server {
	return o.httpServer
}

func (o *Service) GetRouter() *router.Router {
	return o.router
}

func (o *Service) GetLogger() *logrus.Logger {
	return o.logger
}

func (o *Service) Start(bindAddr string) error {
	go func() {
		o.logger.Infof("HTTP: сервис запустился %q", bindAddr)
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

// RequestWrapper добавляет CORS заголовки и content type
func RequestWrapper(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
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

		if next != nil {
			next(ctx)
		}
	}
}
