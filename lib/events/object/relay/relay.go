package relay

import (
	"github.com/pkg/errors"
	"touchon-server/lib/interfaces"
	"touchon-server/lib/messages"
)

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

// OnCheck Проверка состояния реле
type OnCheck struct {
	*messages.MessageImpl
	State string `json:"state"` // Состояние
	Value string `json:"value"` // Значение
}

func NewOnStateOn(targetID int) (interfaces.Message, error) {
	msg, err := messages.NewEvent(interfaces.TargetTypeObject, targetID)
	if err != nil {
		return nil, errors.Wrap(err, "NewOnStateOn")
	}

	o := &OnStateOn{
		MessageImpl: msg,
	}

	return o, nil
}

// OnStateOn Включение реле
type OnStateOn struct {
	*messages.MessageImpl
}

func NewOnStateOff(targetID int) (interfaces.Message, error) {
	msg, err := messages.NewEvent(interfaces.TargetTypeObject, targetID)
	if err != nil {
		return nil, errors.Wrap(err, "NewOnStateOff")
	}

	o := &OnStateOff{
		MessageImpl: msg,
	}

	return o, nil
}

// OnStateOff Выключение реле
type OnStateOff struct {
	*messages.MessageImpl
}
