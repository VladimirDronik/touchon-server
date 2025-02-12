package controller

import (
	"github.com/pkg/errors"
	"touchon-server/lib/interfaces"
	"touchon-server/lib/messages"
)

func NewOnLoad(targetID int) (interfaces.Message, error) {
	msg, err := messages.NewEvent(interfaces.TargetTypeObject, targetID)
	if err != nil {
		return nil, errors.Wrap(err, "NewOnLoad")
	}

	o := &OnLoad{
		MessageImpl: msg,
	}

	return o, nil
}

// OnLoad Инициализация контроллера после включения питания
type OnLoad struct {
	*messages.MessageImpl
}

func NewOnUnavailable(targetID int) (interfaces.Message, error) {
	msg, err := messages.NewEvent(interfaces.TargetTypeObject, targetID)
	if err != nil {
		return nil, errors.Wrap(err, "NewOnUnavailable")
	}

	o := &OnUnavailable{
		MessageImpl: msg,
	}

	return o, nil
}

// OnUnavailable Контроллер стал недоступен
type OnUnavailable struct {
	*messages.MessageImpl
}
