package http

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
	"touchon-server/internal/model"
	"touchon-server/internal/store"
	"touchon-server/lib/helpers"
	"touchon-server/lib/interfaces"
)

// Создание действия по расписанию
// @Summary Создание действия по расписанию
// @Tags CronTasks
// @Description Создание действия по расписанию
// @ID CreateCronAction
// @Accept json
// @Produce json
// @Param object body model.CronTask  true "Действия по расписанию"
// @Success      200 {object} http.Response[model.CronTask]
// @Failure      400 {object} http.Response[any]
// @Failure      500 {object} http.Response[any]
// @Router /cron/task [post]
func (o *Server) handleCreateTask(ctx *fasthttp.RequestCtx) (interface{}, int, error) {
	var req *model.CronTask

	err := json.Unmarshal(ctx.PostBody(), &req)
	if err != nil {
		return nil, fasthttp.StatusBadRequest, err
	}

	taskID, err := store.I.CronRepo().CreateTask(req)
	if err != nil {
		return nil, fasthttp.StatusInternalServerError, err
	}

	for _, action := range req.Actions {
		action.TaskID = taskID

		if err := store.I.CronRepo().CreateTaskAction(action); err != nil {
			return nil, fasthttp.StatusInternalServerError, err
		}
	}

	return taskID, http.StatusOK, nil
}

// Изменение задачи CRON
// @Summary Изменение задачи CRON
// @Tags CronTasks
// @Description Изменение задачи CRON
// @ID UpdateCronTask
// @Accept json
// @Produce json
// @Param object body model.CronTask  true "Действия по расписанию"
// @Success      200 {object} http.Response[model.CronTask]
// @Failure      400 {object} http.Response[any]
// @Failure      500 {object} http.Response[any]
// @Router /cron/task [put]
func (o *Server) handleUpdateTask(ctx *fasthttp.RequestCtx) (interface{}, int, error) {
	var req *model.CronTask

	err := json.Unmarshal(ctx.PostBody(), &req)
	if err != nil {
		return nil, fasthttp.StatusBadRequest, err
	}

	err = store.I.CronRepo().UpdateTask(req)
	if err != nil {
		return nil, fasthttp.StatusInternalServerError, err
	}

	return nil, http.StatusOK, nil
}

// Удаление задание расписания
// @Summary Удаление задание расписания
// @Tags CronTasks
// @Description Удаление задание расписания
// @ID DeleteCronTask
// @Produce json
// @Param target_type query interfaces.TargetType true "Тип сущности" default(item)
// @Param target_id query int true "ID сущности" default(1)
// @Success      200 {object} http.Response[any]
// @Failure      400 {object} http.Response[any]
// @Failure      500 {object} http.Response[any]
// @Router /cron/task [delete]
func (o *Server) handleDeleteTask(ctx *fasthttp.RequestCtx) (interface{}, int, error) {
	targetType := helpers.GetParam(ctx, "target_type")

	if _, ok := interfaces.TargetTypes[targetType]; !ok {
		return nil, http.StatusBadRequest, errors.Errorf("unknown target type %q", targetType)
	}

	targetID, err := helpers.GetUintParam(ctx, "target_id")
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	if err := store.I.CronRepo().DeleteTask(targetID, targetType); err != nil {
		return nil, http.StatusBadRequest, err
	}

	return nil, http.StatusOK, nil
}
