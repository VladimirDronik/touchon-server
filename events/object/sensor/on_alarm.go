package sensor

import (
	"github.com/VladimirDronik/touchon-server/event"
	"github.com/VladimirDronik/touchon-server/models"
	"github.com/VladimirDronik/touchon-server/mqtt/messages"
	"github.com/pkg/errors"
)

func init() {
	maker := func() (*event.Event, error) {
		e := &event.Event{
			Code:        "object.sensor.on_alarm",
			Name:        "Данные датчика вышли за пороговые значения",
			Description: "",
			Props:       event.NewProps(),
			TargetType:  messages.TargetTypeObject,
		}

		t := &event.Prop{
			Code: "msg",
			Name: "Сообщение",
			Item: &models.Item{
				Type: models.DataTypeString,
			},
		}

		if err := e.Props.Add(t); err != nil {
			return nil, errors.Wrap(err, "init.maker")
		}

		return e, nil
	}

	if err := event.Register(maker); err != nil {
		panic(err)
	}
}

func NewOnAlarmMessage(topic string, targetID int, msg string) (messages.Message, error) {
	e, err := event.MakeEvent("object.sensor.on_alarm", messages.TargetTypeObject, targetID, map[string]interface{}{"msg": msg})
	if err != nil {
		return nil, errors.Wrap(err, "NewOnAlarmMessage")
	}

	m, err := e.ToMqttMessage(topic)
	if err != nil {
		return nil, errors.Wrap(err, "NewOnAlarmMessage")
	}

	return m, nil
}
