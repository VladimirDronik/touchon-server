package port

import (
	"github.com/VladimirDronik/touchon-server/event"
	"github.com/VladimirDronik/touchon-server/models"
	"github.com/VladimirDronik/touchon-server/mqtt/messages"
	"github.com/pkg/errors"
)

func init() {
	maker := func() (*event.Event, error) {
		e := &event.Event{
			Code:        "object.port.on_check",
			Name:        "Проверка состояние порта",
			Description: "Событие возникает, когда проверяется состояние порта, но при этом новое пришедшее состояние порта не различается с тем, что хранится в БД",
			Props:       event.NewProps(),
			TargetType:  messages.TargetTypeObject,
		}

		s := &event.Prop{
			Code: "state",
			Name: "Состояние",
			Item: &models.Item{
				Type: models.DataTypeString,
			},
		}

		v := &event.Prop{
			Code: "value",
			Name: "Значение",
			Item: &models.Item{
				Type: models.DataTypeString,
			},
		}

		if err := e.Props.Add(s, v); err != nil {
			return nil, errors.Wrap(err, "init.maker")
		}

		return e, nil
	}

	if err := event.Register(maker); err != nil {
		panic(err)
	}
}

func NewOnCheckMessage(topic string, targetID int, state, value string) (messages.Message, error) {
	e, err := event.MakeEvent("object.port.on_check", messages.TargetTypeObject, targetID, map[string]interface{}{"state": state, "value": value})
	if err != nil {
		return nil, errors.Wrap(err, "NewOnCheckMessage")
	}

	m, err := e.ToMqttMessage(topic)
	if err != nil {
		return nil, errors.Wrap(err, "NewOnCheckMessage")
	}

	return m, nil
}
