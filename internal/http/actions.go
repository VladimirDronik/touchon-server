package http

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
	"touchon-server/internal/helpers"
	"touchon-server/internal/store"
	"touchon-server/lib/interfaces"
)

// Получение количества действий для событий
// @Security TokenAuth
// @Summary Получение количества действий для событий
// @Tags EventActions
// @Description Получение количества действий для событий
// @ID GetEventsActionsCount
// @Produce json
// @Param target_type query interfaces.TargetType true "Тип сущности" Enums(object,item,script) default(item)
// @Param target_id query int true "ID сущности" default(1)
// @Success      200 {object} http.Response[map[string]int]
// @Failure      400 {object} http.Response[any]
// @Failure      500 {object} http.Response[any]
// @Router /events/actions/count [get]
func (o *Server) handleGetEventsActionsCount(ctx *fasthttp.RequestCtx) (interface{}, int, error) {
	targetType := helpers.GetParam(ctx, "target_type")

	if _, ok := interfaces.TargetTypes[targetType]; !ok {
		return nil, http.StatusBadRequest, errors.Errorf("unknown target type %q", targetType)
	}

	targetID, err := helpers.GetUintParam(ctx, "target_id")
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	events, err := store.I.EventsRepo().GetEvents(targetType, targetID)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	eventIDs := make([]int, 0, len(events))
	for _, ev := range events {
		eventIDs = append(eventIDs, ev.ID)
	}

	actionsCountMap, err := store.I.EventActionsRepo().GetActionsCount(eventIDs...)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	r := make(map[string]int, len(actionsCountMap))
	for _, ev := range events {
		r[ev.EventName] = actionsCountMap[ev.ID]
	}

	return r, http.StatusOK, nil
}

// Получение действий для событий
// @Security TokenAuth
// @Summary Получение действий для событий
// @Tags EventActions
// @Description Получение действий для событий
// @ID GetEventsActions
// @Produce json
// @Param target_type query interfaces.TargetType true "Тип сущности" Enums(object,item,script) default(item)
// @Param target_id query int true "ID сущности" default(1)
// @Success      200 {object} http.Response[map[string][]interfaces.EventAction]
// @Failure      400 {object} http.Response[any]
// @Failure      500 {object} http.Response[any]
// @Router /events/actions [get]
func (o *Server) handleGetEventsActions(ctx *fasthttp.RequestCtx) (interface{}, int, error) {
	targetType := helpers.GetParam(ctx, "target_type")

	if _, ok := interfaces.TargetTypes[targetType]; !ok {
		return nil, http.StatusBadRequest, errors.Errorf("unknown target type %q", targetType)
	}

	targetID, err := helpers.GetUintParam(ctx, "target_id")
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	events, err := store.I.EventsRepo().GetEvents(targetType, targetID)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	eventIDs := make([]int, 0, len(events))
	for _, ev := range events {
		eventIDs = append(eventIDs, ev.ID)
	}

	actionsMap, err := store.I.EventActionsRepo().GetActions(eventIDs...)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	r := make(map[string][]*interfaces.EventAction, len(actionsMap))
	for _, ev := range events {
		// Затираем EventID, чтобы не использовать ID вместо event_name
		for _, act := range actionsMap[ev.ID] {
			act.EventID = 0
		}

		r[ev.EventName] = actionsMap[ev.ID]
	}

	return r, http.StatusOK, nil
}

// Создание действия
// @Security TokenAuth
// @Summary Создание действия
// @Tags EventActions
// @Description Создание действия
// @ID CreateEventAction
// @Accept json
// @Produce json
// @Param target_type query interfaces.TargetType true "Тип сущности" default(item)
// @Param target_id query int true "ID сущности" default(1)
// @Param event_name query string true "Название события" default(on_test)
// @Param object body interfaces.EventAction true "Действие"
// @Success      200 {object} http.Response[any]
// @Failure      400 {object} http.Response[any]
// @Failure      500 {object} http.Response[any]
// @Router /events/actions [post]
func (o *Server) handleCreateEventAction(ctx *fasthttp.RequestCtx) (interface{}, int, error) {
	targetType := helpers.GetParam(ctx, "target_type")

	if _, ok := interfaces.TargetTypes[targetType]; !ok {
		return nil, http.StatusBadRequest, errors.Errorf("unknown target type %q", targetType)
	}

	targetID, err := helpers.GetUintParam(ctx, "target_id")
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	eventName := helpers.GetParam(ctx, "event_name")
	// TODO
	//if _, err := event.GetMaker(eventName); err != nil {
	//	return nil, http.StatusBadRequest, err
	//}

	act := &interfaces.EventAction{}
	if err := json.Unmarshal(ctx.Request.Body(), act); err != nil {
		return nil, http.StatusInternalServerError, err
	}

	if err := o.CreateEventAction(targetType, targetID, eventName, act); err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return nil, http.StatusOK, nil
}

func (o *Server) CreateEventAction(targetType string, targetID int, eventName string, act *interfaces.EventAction) error {
	if _, ok := interfaces.TargetTypes[targetType]; !ok {
		return errors.Wrap(errors.Errorf("unknown target type %q", targetType), "CreateEventAction")
	}

	// TODO
	//if _, err := event.GetMaker(eventName); err != nil {
	//	return nil, http.StatusBadRequest, err
	//}

	event, err := store.I.EventsRepo().GetEvent(targetType, targetID, eventName)
	if err != nil {
		if !errors.Is(err, store.ErrNotFound) {
			return errors.Wrap(err, "createEventAction")
		}

		event = &interfaces.AREvent{
			TargetType: targetType,
			TargetID:   targetID,
			EventName:  eventName,
			Enabled:    true,
		}

		if err := store.I.EventsRepo().SaveEvent(event); err != nil {
			return errors.Wrap(err, "createEventAction")
		}
	}

	act.EventID = event.ID

	if err := store.I.EventActionsRepo().SaveAction(act); err != nil {
		return errors.Wrap(err, "createEventAction")
	}

	return nil
}

// Обновление действия
// @Security TokenAuth
// @Summary Обновление действия
// @Tags EventActions
// @Description Обновление действия
// @ID UpdateEventAction
// @Accept json
// @Produce json
// @Param object body interfaces.EventAction true "Действие"
// @Success      200 {object} http.Response[any]
// @Failure      400 {object} http.Response[any]
// @Failure      500 {object} http.Response[any]
// @Router /events/actions [put]
func (o *Server) handleUpdateEventAction(ctx *fasthttp.RequestCtx) (interface{}, int, error) {
	act := &interfaces.EventAction{}
	if err := json.Unmarshal(ctx.Request.Body(), act); err != nil {
		return nil, http.StatusInternalServerError, err
	}

	if err := store.I.EventActionsRepo().SaveAction(act); err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return nil, http.StatusOK, nil
}

// Удаление всех действий по фильтру
// @Security TokenAuth
// @Summary Удаление всех действий по фильтру
// @Tags EventActions
// @Description Удаление всех действий по фильтру
// @ID DeleteAllEventActions
// @Produce json
// @Param target_type query interfaces.TargetType true "Тип сущности" default(object)
// @Param target_id query int true "ID сущности" default(1)
// @Success      200 {object} http.Response[any]
// @Failure      400 {object} http.Response[any]
// @Failure      500 {object} http.Response[any]
// @Router /events/all-actions [delete]
func (o *Server) handleDeleteAllEventActions(ctx *fasthttp.RequestCtx) (interface{}, int, error) {
	targetType := helpers.GetParam(ctx, "target_type")

	if _, ok := interfaces.TargetTypes[targetType]; !ok {
		return nil, http.StatusBadRequest, errors.Errorf("unknown target type %q", targetType)
	}

	targetID, err := helpers.GetUintParam(ctx, "target_id")
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	if err := store.I.EventActionsRepo().DeleteActionByObject(targetType, targetID); err != nil {
		return nil, http.StatusBadRequest, err
	}

	return nil, http.StatusOK, nil
}

// Удаление действия
// @Security TokenAuth
// @Summary Удаление действия
// @Tags EventActions
// @Description Удаление действия
// @ID DeleteEventActions
// @Produce json
// @Param id path int true "ID действия"
// @Success      200 {object} http.Response[any]
// @Failure      400 {object} http.Response[any]
// @Failure      500 {object} http.Response[any]
// @Router /events/actions/{id} [delete]
func (o *Server) handleDeleteEventAction(ctx *fasthttp.RequestCtx) (interface{}, int, error) {
	itemID, err := helpers.GetUintPathParam(ctx, "id")
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	if err := store.I.EventActionsRepo().DeleteAction(itemID); err != nil {
		return nil, http.StatusBadRequest, err
	}

	return nil, http.StatusOK, nil
}

// Смена порядка действий для события
// @Security TokenAuth
// @Summary Смена порядка действий для события
// @Tags EventActions
// @Description Смена порядка действий для события
// @ID OrderEventActions
// @Accept json
// @Produce json
// @Param object body []int true "Упорядоченный список идентификаторов действий"
// @Success      200 {object} http.Response[any]
// @Failure      400 {object} http.Response[any]
// @Failure      500 {object} http.Response[any]
// @Router /events/actions/order [put]
func (o *Server) handleOrderEventActions(ctx *fasthttp.RequestCtx) (interface{}, int, error) {
	actIDs := make([]int, 0, 10)
	if err := json.Unmarshal(ctx.Request.Body(), &actIDs); err != nil {
		return nil, http.StatusInternalServerError, err
	}

	if err := store.I.EventActionsRepo().OrderActions(actIDs); err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return nil, http.StatusOK, nil
}

// Удаление события с действиями
// @Security TokenAuth
// @Summary Удаление события с действиями
// @Tags EventActions
// @Description Удаление события с действиями
// @ID DeleteEvent
// @Accept json
// @Produce json
// @Param target_type query interfaces.TargetType true "Тип сущности" default(item)
// @Param target_id query int true "ID сущности" default(1)
// @Param event_name query string true "Название события. all - удаляет все события" default(all)
// @Success      200 {object} http.Response[any]
// @Failure      400 {object} http.Response[any]
// @Failure      500 {object} http.Response[any]
// @Router /events [delete]
func (o *Server) handleDeleteEvent(ctx *fasthttp.RequestCtx) (interface{}, int, error) {
	targetType := helpers.GetParam(ctx, "target_type")

	if _, ok := interfaces.TargetTypes[targetType]; !ok {
		return nil, http.StatusBadRequest, errors.Errorf("unknown target type %q", targetType)
	}

	targetID, err := helpers.GetUintParam(ctx, "target_id")
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	eventName := helpers.GetParam(ctx, "event_name")

	if err := store.I.EventsRepo().DeleteEvent(targetType, targetID, eventName); err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return nil, http.StatusOK, nil
}
