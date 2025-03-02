package http

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
	"touchon-server/internal/model"
	"touchon-server/internal/store"
	"touchon-server/lib/interfaces"
)

type Event struct {
	Name    string                    `json:"name"`
	Actions []*interfaces.EventAction `json:"actions"`
}

// Определение типа для сваггера
type wizardCreateItemRequest struct {
	Item   model.ViewItem `json:"item"`
	Events []*Event       `json:"events"`
}

// Мастер создания элемента
// @Security TokenAuth
// @Summary Мастер создания элемента
// @Tags Wizard
// @Description Мастер создания элемента
// @ID WizardCreateItem
// @Accept json
// @Produce json
// @Param object body wizardCreateItemRequest true "Элемент"
// @Success      200 {object} Response[model.ViewItem]
// @Failure      400 {object} Response[any]
// @Failure      500 {object} Response[any]
// @Router /private/wizard/create_item [post]
func (o *Server) handleWizardCreateItem(ctx *fasthttp.RequestCtx) (_ interface{}, _ int, e error) {
	type wizardCreateItemRequest struct {
		Item   json.RawMessage `json:"item"`
		Events []*Event        `json:"events"`
	}

	req := &wizardCreateItemRequest{}
	if err := json.Unmarshal(ctx.Request.Body(), req); err != nil {
		return nil, http.StatusBadRequest, err
	}

	// Проверяем названия событий на валидность
	for _, item := range req.Events {
		// TODO
		//if _, err := event.GetMaker(item.Name); err != nil {
		//	return nil, http.StatusBadRequest, err
		//}

		if len(item.Actions) == 0 {
			return nil, http.StatusBadRequest, errors.Errorf("event %q: actions list is empty", item.Name)
		}
	}

	// Сохраняем объект в базу и memstore
	item, status, err := o.saveItem(req.Item)
	if err != nil {
		return nil, status, err
	}

	// Если не закончили транзакцию, удаляем созданный объект
	defer func() {
		if e != nil {
			if _, _, err := o.deleteItem(item.ID); err != nil {
				e = err
				o.GetLogger().Error(err)
			}
		}
	}()

	// Если событий нет, то уходим
	if len(req.Events) == 0 {
		return item, status, nil
	}

	// Если транзакцию не закончили, удаляем события со всеми действиями
	defer func() {
		if e != nil {
			for _, ev := range req.Events {
				if err := store.I.EventsRepo().DeleteEvent(interfaces.TargetTypeItem, item.ID, ev.Name); err != nil {
					e = err
					o.GetLogger().Error(err)
				}
			}
		}
	}()

	for _, ev := range req.Events {
		for _, act := range ev.Actions {
			if err := o.CreateEventAction(interfaces.TargetTypeItem, item.ID, ev.Name, act); err != nil {
				return nil, http.StatusInternalServerError, err
			}
		}
	}

	//Если был указан управляющий объект, то сохраняем его в таблице events для итема
	if item.ControlObject != 0 {
		event := &model.TrEvent{}

		event.EventName = "on_change_state"
		event.TargetType = "object"
		event.TargetID = item.ControlObject

		eventID, err := store.I.Events().AddEvent(event)
		if err != nil || eventID == 0 {
			return nil, http.StatusInternalServerError, err
		}

		eventAction := &model.EventActions{
			EventID:    eventID,
			TargetType: "item",
			TargetID:   item.ID,
			Type:       "method",
			Name:       "set_state",
			Args:       "{\"param\":\"state\"}",
		}

		_, err = store.I.Events().AddEventAction(eventAction)
		if err != nil {
			return nil, http.StatusInternalServerError, err
		}
	}

	return item, http.StatusOK, nil
}
