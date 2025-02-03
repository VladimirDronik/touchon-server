package regulator

import (
	"github.com/pkg/errors"
	"touchon-server/lib/event"
	"touchon-server/lib/models"
	"touchon-server/lib/mqtt/messages"
)

func init() {
	maker := func() (*event.Event, error) {
		e := &event.Event{
			Code:        "object.regulator.on_complex_below_2",
			Name:        "Текущее значение > (targetSP + complexTolerance - belowTolerance)",
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

func NewOnComplexBelow2Message(topic string, targetID int, value float32) (messages.Message, error) {
	e, err := event.MakeEvent("object.regulator.on_complex_below_2", messages.TargetTypeObject, targetID, map[string]interface{}{"value": value})
	if err != nil {
		return nil, errors.Wrap(err, "NewOnComplexBelow2Message")
	}

	m, err := e.ToMqttMessage(topic)
	if err != nil {
		return nil, errors.Wrap(err, "NewOnComplexBelow2Message")
	}

	return m, nil
}
