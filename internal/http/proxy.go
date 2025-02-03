package http

import (
	"net/http"
	"strings"
	"time"

	"github.com/valyala/fasthttp"
	"touchon-server/lib/helpers"
)

// Прокси для сервисов
// @Security TokenAuth
// @Summary Прокси для сервисов
// @Tags Proxy
// @Description Прокси для сервисов
// @ID Proxy
// @Accept */*
// @Produce */*
// @Param service path string true "Service" Enums(action-router,object-manager)
// @Param filepath path string true "Path" default(_/info)
// @Router /proxy/{service}/{filepath} [get]
// @Router /proxy/{service}/{filepath} [post]
// @Router /proxy/{service}/{filepath} [put]
// @Router /proxy/{service}/{filepath} [options]
// @Router /proxy/{service}/{filepath} [patch]
// @Router /proxy/{service}/{filepath} [delete]
func (o *Server) proxy(ctx *fasthttp.RequestCtx) {
	service := helpers.GetPathParam(ctx, "service")
	//path := helpers.GetPathParam(ctx, "filepath")

	addr := o.GetConfig()[strings.ReplaceAll(service, "-", "_")+"_addr"]
	if addr == "" {
		ctx.Error("service addr not found", http.StatusBadRequest)
		return
	}

	ctx.Request.Header.Set("Original-User-Agent", string(ctx.Request.Header.UserAgent()))
	ctx.Request.Header.SetUserAgent(o.GetConfig()["service_name"])
	u := ctx.Request.URI()
	u.SetScheme("http")
	u.SetHost(addr)
	u.SetPath(strings.TrimPrefix(string(u.Path()), "/proxy/"+service))

	if err := o.fasthttpClient.DoTimeout(&ctx.Request, &ctx.Response, 10*time.Second); err != nil {
		ctx.Response.Reset()
		ctx.Error(err.Error(), http.StatusInternalServerError)
		return
	}
}
