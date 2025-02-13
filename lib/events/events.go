package events

import (
	"github.com/pkg/errors"
	"touchon-server/lib/event"
	"touchon-server/lib/interfaces"
	"touchon-server/lib/messages"
)

func NewOnChangeState(targetType interfaces.TargetType, targetID int, state, value string) (interfaces.Event, error) {
	msg, err := messages.NewEvent(targetType, targetID)
	if err != nil {
		return nil, errors.Wrap(err, "NewOnChangeState")
	}

	o := &OnChangeState{
		Event: &event.EventImpl{
			Message:          msg,
			EventCode:        "on_change_state",
			EventName:        "on_change_state",
			EventDescription: "Изменение состояния",
		},
		State: state,
		Value: value,
	}

	return o, nil
}

type OnChangeState struct {
	interfaces.Event
	State string `json:"state,omitempty"` // Состояние
	Value string `json:"value,omitempty"` // Значение
}

func NewOnError(targetType interfaces.TargetType, targetID int, errMsg string) (interfaces.Event, error) {
	msg, err := messages.NewEvent(targetType, targetID)
	if err != nil {
		return nil, errors.Wrap(err, "NewOnError")
	}

	o := &OnError{
		Event: &event.EventImpl{
			Message:          msg,
			EventCode:        "on_error",
			EventName:        "on_error",
			EventDescription: "Ошибка",
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
	msg, err := messages.NewMessage(interfaces.MessageTypeNotification, "", interfaces.TargetTypeNotMatters, 0)
	if err != nil {
		return nil, errors.Wrap(err, "NewNotification")
	}

	o := &Notification{
		Event: &event.EventImpl{
			Message:          msg,
			EventCode:        "on_notify",
			EventName:        "on_notify",
			EventDescription: "Уведомление",
		},
		Type: msgType,
		Text: msgText,
	}

	return o, nil
}

type Notification struct {
	interfaces.Event
	Type interfaces.NotificationType `json:"type"`           // Тип
	Text string                      `json:"text,omitempty"` // Текст
}
