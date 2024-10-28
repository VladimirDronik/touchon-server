package controller

import (
	"github.com/VladimirDronik/touchon-server/event"
	"github.com/VladimirDronik/touchon-server/mqtt/messages"
	"github.com/pkg/errors"
)

func init() {
	maker := func() (*event.Event, error) {
		e := &event.Event{
			Code:        "object.controller.on_load",
			Name:        "Инициализация контроллера после включения питания",
			Description: "",
			Props:       event.NewProps(),
			TargetType:  messages.TargetTypeObject,
		}

		return e, nil
	}

	if err := event.Register(maker); err != nil {
		panic(err)
	}
}

func NewOnLoadMessage(topic string, targetID int) (messages.Message, error) {
	e, err := event.MakeEvent("object.controller.on_load", messages.TargetTypeObject, targetID, nil)
	if err != nil {
		return nil, errors.Wrap(err, "NewOnLoadMessage")
	}

	m, err := e.ToMqttMessage(topic)
	if err != nil {
		return nil, errors.Wrap(err, "NewOnLoadMessage")
	}

	return m, nil
}
