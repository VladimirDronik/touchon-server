package http

import (
	"encoding/json"
	"net/http"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
	"touchon-server/internal/g"
	"touchon-server/internal/helpers"
	"touchon-server/internal/model"
	"touchon-server/internal/store"
	"touchon-server/internal/token"
	"touchon-server/lib/http/server"
	"touchon-server/lib/interfaces"
)

func (o *Server) disableAuth(ctx *fasthttp.RequestCtx) *model.Tokens {
	if tokenSecret := o.GetConfig()["token_secret"]; tokenSecret == "disable_auth" {
		tokens, err := helpers.CreateSession(g.DisabledAuthDeviceID)
		if err != nil {
			return nil
		}

		o.setCookie(ctx, "refreshToken", tokens.RefreshToken, true)

		return tokens
	}

	return nil
}

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
	// Disable auth
	if tokens := o.disableAuth(ctx); tokens != nil {
		return tokens, http.StatusOK, nil
	}

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
		if _, err := store.I.Users().GetByDeviceID(deviceID); err != nil {
			return nil, http.StatusUnauthorized, err
		}

	case login != "" && password != "":
		user, err := store.I.Users().GetByLoginAndPassword(login, password)
		if err != nil {
			return nil, http.StatusUnauthorized, err
		}
		deviceID = user.DeviceID

	default:
		return nil, http.StatusUnauthorized, errors.New("credentials not found")
	}

	tokens, err := helpers.CreateSession(deviceID)
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
	// Disable auth
	if tokens := o.disableAuth(ctx); tokens != nil {
		return tokens, http.StatusOK, nil
	}

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

	user, err := store.I.Users().GetByToken(refreshToken)
	if err != nil {
		return nil, http.StatusUnauthorized, err
	}

	if user.DeviceID <= 0 {
		return nil, http.StatusUnauthorized, err
	}

	tokens, err := helpers.CreateSession(user.DeviceID)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	// Обновление refreshToken в httpOnly cookie
	o.setCookie(ctx, "refreshToken", tokens.RefreshToken, true)

	return tokens, http.StatusOK, nil
}

type Middleware func(ctx *fasthttp.RequestCtx, next interfaces.RequestHandler) (interface{}, int, error)

func (o *Server) authMiddleware(ctx *fasthttp.RequestCtx, next interfaces.RequestHandler) (interface{}, int, error) {
	tokenSecret := o.GetConfig()["token_secret"]
	if tokenSecret == "disable_auth" {
		// Disable auth
		ctx.SetUserValue("device_id", g.DisabledAuthDeviceID)
		return next(ctx)
	}

	tkn := string(ctx.Request.Header.Peek("Token"))
	if tkn == "" {
		return nil, http.StatusUnauthorized, errors.New("token not found")
	}

	o.GetLogger().Debugf("authMiddleware: Token header: %s", tkn)

	// Проверяем, не протух ли токен и извлекаем ID юзера
	deviceID, err := token.KeysExtract(tkn, tokenSecret)
	if err != nil {
		return nil, http.StatusUnauthorized, err
	}

	if deviceID <= 0 {
		return nil, http.StatusUnauthorized, err
	}

	ctx.SetUserValue("device_id", deviceID)

	return next(ctx)
}

func (o *Server) addMiddleware(pathPrefix string, middleware Middleware) func(method, path string, handler interfaces.RequestHandler) {
	return func(method, path string, handler interfaces.RequestHandler) {
		// filepath.ToSlash() - for windows
		o.GetRouter().Handle(method, filepath.ToSlash(filepath.Join(pathPrefix, path)), server.JsonHandlerWrapper(func(ctx *fasthttp.RequestCtx) (interface{}, int, error) {
			return middleware(ctx, handler)
		}))
	}
}

func (o *Server) addRawMiddleware(pathPrefix string, middleware Middleware) func(method, path string, handler fasthttp.RequestHandler) {
	return func(method, path string, handler fasthttp.RequestHandler) {
		// filepath.ToSlash() - for windows
		o.GetRouter().Handle(method, filepath.ToSlash(filepath.Join(pathPrefix, path)), func(ctx *fasthttp.RequestCtx) {
			_, status, err := middleware(ctx, func(ctx *fasthttp.RequestCtx) (interface{}, int, error) {
				handler(ctx)
				return nil, 0, nil
			})

			if err != nil {
				ctx.Error(err.Error(), status)
			}
		})
	}
}
