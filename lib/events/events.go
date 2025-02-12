package events

import (
	"github.com/pkg/errors"
	"touchon-server/lib/interfaces"
	"touchon-server/lib/messages"
)

func NewOnChangeState(targetType interfaces.TargetType, targetID int, state, value string) (interfaces.Message, error) {
	msg, err := messages.NewEvent(targetType, targetID)
	if err != nil {
		return nil, errors.Wrap(err, "NewOnChangeState")
	}

	o := &OnChangeState{
		MessageImpl: msg,
		State:       state,
		Value:       value,
	}

	return o, nil
}

type OnChangeState struct {
	*messages.MessageImpl
	State string `json:"state"` // Состояние
	Value string `json:"value"` // Значение
}

func NewOnError(targetType interfaces.TargetType, targetID int, errMsg string) (interfaces.Message, error) {
	msg, err := messages.NewEvent(targetType, targetID)
	if err != nil {
		return nil, errors.Wrap(err, "NewOnError")
	}

	o := &OnError{
		MessageImpl: msg,
		Error:       errMsg,
	}

	return o, nil
}

type OnError struct {
	*messages.MessageImpl
	Error string `json:"error"` // Текст ошибки
}

func NewNotification(msgType interfaces.NotificationType, msgText string) (interfaces.Message, error) {
	msg, err := messages.NewMessage(interfaces.MessageTypeNotification, "", interfaces.TargetTypeNotMatters, 0)
	if err != nil {
		return nil, errors.Wrap(err, "NewNotification")
	}

	o := &Notification{
		MessageImpl: msg,
		Type:        msgType,
		Text:        msgText,
	}

	return o, nil
}

type Notification struct {
	*messages.MessageImpl
	Type interfaces.NotificationType `json:"type"` // Тип
	Text string                      `json:"text"` // Текст
}
