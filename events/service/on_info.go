package service

import (
	"github.com/VladimirDronik/touchon-server/event"
	"github.com/VladimirDronik/touchon-server/info"
	"github.com/VladimirDronik/touchon-server/models"
	"github.com/VladimirDronik/touchon-server/mqtt/messages"
	"github.com/pkg/errors"
)

func init() {
	maker := func() (*event.Event, error) {
		e := &event.Event{
			Code:        "service.on_info",
			Name:        "Информация о сервисе",
			Description: "",
			Props:       event.NewProps(),
			TargetType:  messages.TargetTypeService,
		}

		msg := &event.Prop{
			Code: "info",
			Name: "Информация",
			Item: &models.Item{
				Type: models.DataTypeInterface,
			},
		}

		if err := e.Props.Add(msg); err != nil {
			return nil, errors.Wrap(err, "init.maker")
		}

		return e, nil
	}

	if err := event.Register(maker); err != nil {
		panic(err)
	}
}

func NewOnInfoMessage(topic string) (messages.Message, error) {
	nfo, err := info.GetInfo()
	if err != nil {
		return nil, errors.Wrap(err, "NewOnInfoMessage")
	}

	e, err := event.MakeEvent("service.on_info", messages.TargetTypeService, 0, map[string]interface{}{"info": nfo})
	if err != nil {
		return nil, errors.Wrap(err, "NewOnInfoMessage")
	}

	m, err := e.ToMqttMessage(topic)
	if err != nil {
		return nil, errors.Wrap(err, "NewOnInfoMessage")
	}

	return m, nil
}
