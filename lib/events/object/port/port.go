package port

import (
	"github.com/pkg/errors"
	"touchon-server/lib/interfaces"
	"touchon-server/lib/messages"
)

func NewOnPress(targetID int) (interfaces.Message, error) {
	msg, err := messages.NewEvent(interfaces.TargetTypeObject, targetID)
	if err != nil {
		return nil, errors.Wrap(err, "NewOnPress")
	}

	o := &OnPress{
		MessageImpl: msg,
	}

	return o, nil
}

// OnPress Порт замкнут
type OnPress struct {
	*messages.MessageImpl
}

func NewOnLongPress(targetID int) (interfaces.Message, error) {
	msg, err := messages.NewEvent(interfaces.TargetTypeObject, targetID)
	if err != nil {
		return nil, errors.Wrap(err, "NewOnLongPress")
	}

	o := &OnLongPress{
		MessageImpl: msg,
	}

	return o, nil
}

// OnLongPress Порт удерживается в замкнутом состоянии
type OnLongPress struct {
	*messages.MessageImpl
}

func NewOnRelease(targetID int) (interfaces.Message, error) {
	msg, err := messages.NewEvent(interfaces.TargetTypeObject, targetID)
	if err != nil {
		return nil, errors.Wrap(err, "NewOnRelease")
	}

	o := &OnRelease{
		MessageImpl: msg,
	}

	return o, nil
}

// OnRelease Порт разомкнут
type OnRelease struct {
	*messages.MessageImpl
}

func NewOnDoubleClick(targetID int) (interfaces.Message, error) {
	msg, err := messages.NewEvent(interfaces.TargetTypeObject, targetID)
	if err != nil {
		return nil, errors.Wrap(err, "NewOnDoubleClick")
	}

	o := &OnDoubleClick{
		MessageImpl: msg,
	}

	return o, nil
}

// OnDoubleClick Двойное замыкание
type OnDoubleClick struct {
	*messages.MessageImpl
}

func NewOnCheck(targetID int, state, value string) (interfaces.Message, error) {
	msg, err := messages.NewEvent(interfaces.TargetTypeObject, targetID)
	if err != nil {
		return nil, errors.Wrap(err, "NewOnCheck")
	}

	o := &OnCheck{
		MessageImpl: msg,
		State:       state,
		Value:       value,
	}

	return o, nil
}

// OnCheck Событие возникает, когда проверяется состояние порта,
// но при этом новое пришедшее состояние порта не различается с тем, что хранится в БД
type OnCheck struct {
	*messages.MessageImpl
	State string `json:"state"` // Состояние
	Value string `json:"value"` // Значение
}
