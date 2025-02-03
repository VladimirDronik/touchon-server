package http

import (
	"encoding/json"
	"net/http"
	"path/filepath"

	"github.com/VladimirDronik/touchon-server/helpers"
	"github.com/VladimirDronik/touchon-server/http/server"
	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
	"translator/internal/token"
)

// Генерация нового токена
// @Summary Генерация нового токена
// @Tags Auth
// @Description Генерация нового токена
// @ID GetToken
// @Produce json
// @Param api-key header string true "API key" default(c041d36e381a835afce48c91686370c8)
// @Param login query string false "Login" default(web)
// @Param password query string false "Password" default(12345)
// @Success      200 {object} Response[model.Tokens]
// @Failure      400 {object} Response[any]
// @Failure      500 {object} Response[any]
// @Router /token [get]
func (o *Server) getToken(ctx *fasthttp.RequestCtx) (interface{}, int, error) {
	// Сравниваем полученный server_key с тем, что лежит в конфиге
	apiKey := string(ctx.Request.Header.Peek("api-key"))
	if apiKey != o.GetConfig()["server_key"] {
		return nil, http.StatusBadRequest, errServerKeyIncorrect
	}

	deviceID, _ := helpers.GetUintParam(ctx, "user_id")
	login := helpers.GetParam(ctx, "login")
	password := helpers.GetParam(ctx, "password")

	switch {
	case deviceID > 0:
		if _, err := o.store.Users().GetByDeviceID(deviceID); err != nil {
			return nil, http.StatusUnauthorized, err
		}

	case login != "" && password != "":
		user, err := o.store.Users().GetByLoginAndPassword(login, password)
		if err != nil {
			return nil, http.StatusUnauthorized, err
		}
		deviceID = user.DeviceID

	default:
		return nil, http.StatusUnauthorized, errors.New("credentials not found")
	}

	tokens, err := o.createSession(deviceID)
	if err != nil {
		return nil, http.StatusUnauthorized, err
	}

	o.setCookie(ctx, "refreshToken", tokens.RefreshToken, true)

	return tokens, http.StatusOK, nil
}

// Генерация новой пары токенов
// @Summary Генерация новой пары токенов
// @Tags Auth
// @Description Генерация новой пары токенов
// @ID RefreshToken
// @Produce json
// @Param Cookie header string true "Refresh token" default(refreshToken=сюда_вставить_refreshToken)
// @Success      200 {object} Response[model.Tokens]
// @Failure      400 {object} Response[any]
// @Failure      500 {object} Response[any]
// @Router /token [post]
func (o *Server) refreshToken(ctx *fasthttp.RequestCtx) (interface{}, int, error) {
	// Попытка получить refreshToken из httpOnly cookie
	refreshToken := string(ctx.Request.Header.Cookie("refreshToken"))
	if refreshToken == "" {
		if body := ctx.Request.Body(); json.Valid(body) {
			// Если токен в куках не найден, пробуем получить его из тела запроса
			type request struct {
				RefreshToken string `json:"apiRefreshToken"`
			}

			req := &request{}
			if err := json.Unmarshal(body, req); err != nil {
				return nil, http.StatusBadRequest, err
			}

			refreshToken = req.RefreshToken
		}
	}

	// Проверяем наличие refreshToken
	if refreshToken == "" {
		return nil, http.StatusUnauthorized, errors.New("refresh token is missing")
	}

	user, err := o.store.Users().GetByToken(refreshToken)
	if err != nil {
		return nil, http.StatusUnauthorized, err
	}

	tokens, err := o.createSession(user.DeviceID)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	// Обновление refreshToken в httpOnly cookie
	o.setCookie(ctx, "refreshToken", tokens.RefreshToken, true)

	return tokens, http.StatusOK, nil
}

type Middleware func(ctx *fasthttp.RequestCtx, next server.RequestHandler) (interface{}, int, error)

func (o *Server) authMiddleware(ctx *fasthttp.RequestCtx, next server.RequestHandler) (interface{}, int, error) {
	tokenSecret := o.GetConfig()["token_secret"]
	if tokenSecret == "disable_auth" {
		// Disable auth
		ctx.SetUserValue("device_id", 10)
		return next(ctx)
	}

	tkn := string(ctx.Request.Header.Peek("Token"))

	if tkn == "" {
		return nil, http.StatusBadRequest, errors.New("token not found")
	}

	o.GetLogger().Debugf("authMiddleware: Token header: %s", tkn)

	// Проверяем, не протух ли токен и извлекаем ID юзера
	deviceID, err := token.KeysExtract(tkn, tokenSecret)
	if err != nil {
		return nil, http.StatusUnauthorized, err
	}

	ctx.SetUserValue("device_id", deviceID)

	return next(ctx)
}

func (o *Server) addMiddleware(pathPrefix string, middleware Middleware) func(method, path string, handler server.RequestHandler) {
	return func(method, path string, handler server.RequestHandler) {
		//o.GetRouter().Handle(method, filepath.ToSlash(filepath.Join(pathPrefix, path)), http.JsonHandlerWrapper(func(ctx *fasthttp.RequestCtx) (interface{}, int, error) {
		//	return middleware(ctx, handler)
		//}))

		// Backward compatibility (filepath.ToSlash() - for windows)
		o.GetRouter().Handle(method, filepath.ToSlash(filepath.Join(pathPrefix, path)), JsonHandlerWrapper(func(ctx *fasthttp.RequestCtx) (interface{}, int, error) {
			return middleware(ctx, handler)
		}))
	}
}

type MiddlewareRaw func(ctx *fasthttp.RequestCtx, next fasthttp.RequestHandler)

func (o *Server) authMiddlewareRaw(ctx *fasthttp.RequestCtx, next fasthttp.RequestHandler) {
	tokenSecret := o.GetConfig()["token_secret"]
	if tokenSecret == "disable_auth" {
		// Disable auth
		ctx.SetUserValue("device_id", 10)
		next(ctx)
		return
	}

	tkn := string(ctx.Request.Header.Peek("Token"))

	if tkn == "" {
		ctx.Error("token not found", http.StatusBadRequest)
		return
	}

	o.GetLogger().Debugf("authMiddlewareRaw: Token header: %s", tkn)

	// Проверяем, не протух ли токен и извлекаем ID юзера
	deviceID, err := token.KeysExtract(tkn, tokenSecret)
	if err != nil {
		ctx.Error(err.Error(), http.StatusUnauthorized)
		return
	}

	ctx.SetUserValue("device_id", deviceID)

	next(ctx)
}

func (o *Server) addMiddlewareRaw(pathPrefix string, middleware MiddlewareRaw) func(method, path string, handler fasthttp.RequestHandler) {
	return func(method, path string, handler fasthttp.RequestHandler) {
		// filepath.ToSlash() - for windows
		o.GetRouter().Handle(method, filepath.ToSlash(filepath.Join(pathPrefix, path)), func(ctx *fasthttp.RequestCtx) {
			middleware(ctx, handler)
		})
	}
}
