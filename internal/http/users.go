package http

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
	"touchon-server/internal/helpers"
	"touchon-server/internal/model"
	"touchon-server/internal/store"
)

// deviceTokenRequest представляет тело запроса для привязки токена устройства
type deviceTokenRequest struct {
	DeviceType  string `json:"deviceType" enums:"android,ios"`
	DeviceToken string `json:"deviceToken" default:"token"`
}

type userRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	DeviceID int    `json:"device_id"`
	SendPush bool   `json:"send_push"`
}

// Установка токена устройства для отправки push уведомлений
// @Security TokenAuth
// @Summary Установка токена устройства для отправки push уведомлений
// @Tags Users
// @Description Установка токена устройства для отправки push уведомлений
// @ID LinkDeviceToken
// @Produce json
// @Param body body deviceTokenRequest true "Body"
// @Success      200 {object} Response[any]
// @Failure      400 {object} Response[any]
// @Failure      500 {object} Response[any]
// @Router /private/users/link-token [patch]
func (o *Server) linkDeviceToken(ctx *fasthttp.RequestCtx) (interface{}, int, error) {
	var req deviceTokenRequest

	// Декодируем JSON-тело запроса в структуру DeviceTokenRequest
	if err := json.Unmarshal(ctx.Request.Body(), &req); err != nil {
		return nil, http.StatusBadRequest, err
	}

	deviceID, ok := ctx.UserValue("device_id").(int)
	if !ok {
		return nil, http.StatusBadRequest, errors.New("device_id not found in context")
	}

	if err := store.I.Users().LinkDeviceToken(deviceID, req.DeviceToken, req.DeviceType); err != nil {
		return nil, http.StatusBadRequest, err
	}

	return nil, http.StatusOK, nil
}

// Получение всех юзеров
// @Security TokenAuth
// @Summary Получение всех юзеров
// @Tags Users
// @Description Получение всех юзеров
// @ID GetAllUsers
// @Produce json
// @Success      200 {object} model.User
// @Failure      400 {object} Response[any]
// @Failure      500 {object} Response[any]
// @Router /private/users [get]
func (o *Server) handleGetAllUsers(ctx *fasthttp.RequestCtx) (interface{}, int, error) {
	var users []*model.User

	users, err := store.I.Users().GetAllUsers()
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return users, http.StatusOK, nil
}

// Добавление юзера
// @Security TokenAuth
// @Summary Добавление юзера
// @Tags Users
// @Description Добавление юзера
// @ID Create
// @Produce json
// @Param body body userRequest true "Body"
// @Success      200 {object} Response[any]
// @Failure      400 {object} Response[any]
// @Failure      500 {object} Response[any]
// @Router /private/users [post]
func (o *Server) handleCreateUser(ctx *fasthttp.RequestCtx) (interface{}, int, error) {
	var req userRequest
	user := model.User{}

	if err := json.Unmarshal(ctx.Request.Body(), &req); err != nil {
		return nil, http.StatusBadRequest, err
	}

	user.Login = req.Login
	user.DeviceID = req.DeviceID
	user.Password = req.Password
	user.SendPush = req.SendPush

	userID, err := store.I.Users().Create(&user)

	return userID, http.StatusOK, err
}

// Удаление пользователя
// @Security TokenAuth
// @Summary Удаление пользователя
// @Tags Users
// @Description Удаление пользователя
// @ID DeleteUser
// @Accept json
// @Produce json
// @Param id query int true "ID пользователя"
// @Success      200 {object} Response[any]
// @Failure      400 {object} Response[any]
// @Failure      500 {object} Response[any]
// @Router /private/users [delete]
func (o *Server) handleDeleteUser(ctx *fasthttp.RequestCtx) (interface{}, int, error) {
	userID, err := helpers.GetUintParam(ctx, "id")
	if err != nil {
		return nil, http.StatusBadRequest, err
	}
	if err := store.I.Users().Delete(userID); err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return nil, http.StatusOK, nil
}
