package controller

import (
	"github.com/pkg/errors"
	"touchon-server/lib/event"
	"touchon-server/lib/interfaces"
	"touchon-server/lib/messages"
)

func NewOnLoad(targetID int) (interfaces.Event, error) {
	msg, err := messages.NewEvent(interfaces.TargetTypeObject, targetID)
	if err != nil {
		return nil, errors.Wrap(err, "NewOnLoad")
	}

	o := &OnLoad{
		Event: &event.EventImpl{
			Message:          msg,
			EventCode:        "object.controller.on_load",
			EventName:        "on_load",
			EventDescription: "Инициализация контроллера после включения питания",
		},
	}

	return o, nil
}

// OnLoad Инициализация контроллера после включения питания
type OnLoad struct {
	interfaces.Event
}

func NewOnUnavailable(targetID int) (interfaces.Event, error) {
	msg, err := messages.NewEvent(interfaces.TargetTypeObject, targetID)
	if err != nil {
		return nil, errors.Wrap(err, "NewOnUnavailable")
	}

	o := &OnUnavailable{
		Event: &event.EventImpl{
			Message:          msg,
			EventCode:        "object.controller.on_unavailable",
			EventName:        "on_unavailable",
			EventDescription: "Контроллер стал недоступен",
		},
	}

	return o, nil
}

// OnUnavailable Контроллер стал недоступен
type OnUnavailable struct {
	interfaces.Event
}
