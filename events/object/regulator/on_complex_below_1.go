package regulator

import (
	"github.com/VladimirDronik/touchon-server/event"
	"github.com/VladimirDronik/touchon-server/models"
	"github.com/VladimirDronik/touchon-server/mqtt/messages"
	"github.com/pkg/errors"
)

func init() {
	maker := func() (*event.Event, error) {
		e := &event.Event{
			Code:        "object.regulator.on_complex_below_1",
			Name:        "Текущее значение < (targetSP - complexTolerance - belowTolerance)",
			Description: "",
			Props:       event.NewProps(),
			TargetType:  messages.TargetTypeObject,
		}

		v := &event.Prop{
			Code: "value",
			Name: "Значение",
			Item: &models.Item{
				Type:       models.DataTypeFloat,
				RoundFloat: true,
			},
		}

		if err := e.Props.Add(v); err != nil {
			return nil, errors.Wrap(err, "init.maker")
		}

		return e, nil
	}

	// Для регистрации событий надо в service/init.go добавить импорт данного _пакета_!
	if err := event.Register(maker); err != nil {
		panic(err)
	}
}

func NewOnComplexBelow1Message(topic string, targetID int, value float32) (messages.Message, error) {
	e, err := event.MakeEvent("object.regulator.on_complex_below_1", messages.TargetTypeObject, targetID, map[string]interface{}{"value": value})
	if err != nil {
		return nil, errors.Wrap(err, "NewOnComplexBelow1Message")
	}

	m, err := e.ToMqttMessage(topic)
	if err != nil {
		return nil, errors.Wrap(err, "NewOnComplexBelow1Message")
	}

	return m, nil
}
