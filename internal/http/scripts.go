package http

import (
	"encoding/json"
	"net/http"

	"github.com/valyala/fasthttp"
	"touchon-server/internal/helpers"
	"touchon-server/internal/model"
	"touchon-server/internal/scripts"
	"touchon-server/internal/store"
	"touchon-server/lib/models"
)

// Получение модели сценария
// @Summary Получение модели сценария
// @Tags Scripts
// @Description Получение модели сценария
// @ID GetScriptModel
// @Produce json
// @Success      200 {object} http.Response[scripts.Script]
// @Failure      400 {object} http.Response[any]
// @Failure      500 {object} http.Response[any]
// @Router /scripts/model [get]
func (o *Server) handleGetScriptModel(ctx *fasthttp.RequestCtx) (interface{}, int, error) {
	params := []*scripts.Param{{
		Code:        "timeout",
		Name:        "Таймаут",
		Description: "Сколько будем спать",
		Item: &models.Item{
			Type: models.DataTypeInt,
		},
	}}

	s, err := scripts.NewScript("my_script", "Сценарий", "Пример сценария", params, "")
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	if err := s.Check(); err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return s, http.StatusOK, nil
}

// Получение сценария
// @Summary Получение сценария
// @Tags Scripts
// @Description Получение сценария
// @ID GetScript
// @Produce json
// @Param id path int true "ID сценария"
// @Success      200 {object} http.Response[scripts.Script]
// @Failure      400 {object} http.Response[any]
// @Failure      500 {object} http.Response[any]
// @Router /scripts/{id} [get]
func (o *Server) handleGetScript(ctx *fasthttp.RequestCtx) (interface{}, int, error) {
	id, err := helpers.GetUintPathParam(ctx, "id")
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	row, err := store.I.ScriptRepository().GetScript(id)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return row, http.StatusOK, nil
}

// Создание сценария
// @Summary Создание сценария
// @Tags Scripts
// @Description Создание сценария
// @ID CreateScript
// @Accept json
// @Produce json
// @Param object body scripts.Script true "Сценарий"
// @Success      200 {object} http.Response[scripts.Script]
// @Failure      400 {object} http.Response[any]
// @Failure      500 {object} http.Response[any]
// @Router /scripts [post]
func (o *Server) handleCreateScript(ctx *fasthttp.RequestCtx) (interface{}, int, error) {
	return o.handleSaveScript(ctx)
}

// Обновление сценария
// @Summary Обновление сценария
// @Tags Scripts
// @Description Обновление сценария
// @ID UpdateScript
// @Accept json
// @Produce json
// @Param object body scripts.Script true "Сценарий"
// @Success      200 {object} http.Response[scripts.Script]
// @Failure      400 {object} http.Response[any]
// @Failure      500 {object} http.Response[any]
// @Router /scripts [put]
func (o *Server) handleUpdateScript(ctx *fasthttp.RequestCtx) (interface{}, int, error) {
	return o.handleSaveScript(ctx)
}

func (o *Server) handleSaveScript(ctx *fasthttp.RequestCtx) (interface{}, int, error) {
	s, err := scripts.NewScript("", "", "", nil, "")
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	if err := json.Unmarshal(ctx.Request.Body(), s); err != nil {
		return nil, http.StatusInternalServerError, err
	}

	params, err := json.MarshalIndent(s.Params, "", "  ")
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	script := &model.StoreScript{
		ID:          s.ID,
		Code:        s.Code,
		Name:        s.Name,
		Description: s.Description,
		Params:      params,
		Body:        s.Body,
	}

	if err := store.I.ScriptRepository().SetScript(script); err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return script, http.StatusOK, nil
}

type GetScriptsResponse struct {
	Total int                  `json:"total"`
	List  []*model.StoreScript `json:"list"`
}

// Получение списка сценариев
// @Summary Получение сценариев
// @Tags Scripts
// @Description Получение всех сценариев
// @ID GetScripts
// @Produce json
// @Param code   query string false "Шаблон кода" default(my_)
// @Param name   query string false "Шаблон названия" default(цен)
// @Param offset query string false "Смещение" default(0)
// @Param limit  query string false "Лимит" default(20)
// @Success      200 {object} http.Response[GetScriptsResponse]
// @Failure      400 {object} http.Response[any]
// @Failure      500 {object} http.Response[any]
// @Router /scripts [get]
func (o *Server) handleGetScripts(ctx *fasthttp.RequestCtx) (interface{}, int, error) {
	code := helpers.GetParam(ctx, "code")
	name := helpers.GetParam(ctx, "name")

	offset, err := helpers.GetUintParam(ctx, "offset")
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	limit, err := helpers.GetUintParam(ctx, "limit")
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	if limit == 0 {
		limit = 20
	}

	rows, err := store.I.ScriptRepository().GetScripts(code, name, offset, limit)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	total, err := store.I.ScriptRepository().GetTotal(code, name, offset, limit)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return GetScriptsResponse{
		Total: total,
		List:  rows,
	}, http.StatusOK, nil
}

// Удаление сценария
// @Summary Удаление сценария
// @Tags Scripts
// @Description Удаление сценария
// @ID DeleteScript
// @Produce json
// @Param id path int true "ID сценария"
// @Success      200 {object} http.Response[any]
// @Failure      400 {object} http.Response[any]
// @Failure      500 {object} http.Response[any]
// @Router /scripts/{id} [delete]
func (o *Server) handleDeleteScript(ctx *fasthttp.RequestCtx) (interface{}, int, error) {
	id, err := helpers.GetUintPathParam(ctx, "id")
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	if err := store.I.ScriptRepository().DelScript(id); err != nil {
		return nil, http.StatusBadRequest, err
	}

	return nil, http.StatusOK, nil
}

// Запуск сценария
// @Summary Запуск сценария
// @Tags Scripts
// @Description Запуск сценария
// @ID ExecScript
// @Produce json
// @Param id path int true "ID сценария" default(2)
// @Param timeout query int true "Время сна" default(5)
// @Success      200 {object} http.Response[any]
// @Failure      400 {object} http.Response[any]
// @Failure      500 {object} http.Response[any]
// @Router /scripts/{id}/exec [post]
func (o *Server) handleExecScript(ctx *fasthttp.RequestCtx) (interface{}, int, error) {
	id, err := helpers.GetUintPathParam(ctx, "id")
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	s, err := scripts.I.GetScript(id)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	args := make(map[string]interface{}, 10)
	ctx.Request.URI().QueryArgs().VisitAll(func(k, v []byte) {
		args[string(k)] = string(v)
	})

	r, err := scripts.I.ExecScript(s, args)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return r, http.StatusOK, nil
}
