package events

import (
	"github.com/pkg/errors"
	"touchon-server/lib/event"
	"touchon-server/lib/interfaces"
	"touchon-server/lib/messages"
)

func NewOnChangeState(targetType interfaces.TargetType, targetID int, state, value string) (interfaces.Event, error) {
	msg, err := messages.NewEvent("on_change_state", targetType, targetID)
	if err != nil {
		return nil, errors.Wrap(err, "NewOnChangeState")
	}

	o := &OnChangeState{
		Event: &event.EventImpl{
			Message:     msg,
			Title:       "on_change_state",
			Description: "Изменение состояния",
		},
	}

	o.SetState(state)

	return o, nil
}

type OnChangeState struct {
	interfaces.Event
	State string `json:"state,omitempty"` // Состояние
	Value string `json:"value,omitempty"` // Значение
}

func NewOnError(targetType interfaces.TargetType, targetID int, errMsg string) (interfaces.Event, error) {
	msg, err := messages.NewEvent("on_error", targetType, targetID)
	if err != nil {
		return nil, errors.Wrap(err, "NewOnError")
	}

	o := &OnError{
		Event: &event.EventImpl{
			Message:     msg,
			Title:       "on_error",
			Description: "Ошибка",
		},
		Error: errMsg,
	}

	return o, nil
}

type OnError struct {
	interfaces.Event
	Error string `json:"error,omitempty"` // Текст ошибки
}

func NewNotification(msgType interfaces.NotificationType, msgText string) (interfaces.Event, error) {
	msg, err := messages.NewEvent("on_notify", interfaces.TargetTypeNotMatters, 0)
	if err != nil {
		return nil, errors.Wrap(err, "NewNotification")
	}

	o := &Notification{
		Event: &event.EventImpl{
			Message:     msg,
			Title:       "on_notify",
			Description: "Уведомление",
		},
	}

	o.SetValue("type", msgType) // Тип
	o.SetValue("text", msgText) // Текст

	return o, nil
}

type Notification struct {
	interfaces.Event
}
