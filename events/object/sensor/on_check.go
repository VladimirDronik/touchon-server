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
			Code:        "object.sensor.on_check",
			Name:        "Данные датчика обновлены",
			Description: "",
			Props:       event.NewProps(),
			TargetType:  messages.TargetTypeObject,
		}

		t := &event.Prop{
			Code: "temperature",
			Name: "Температура",
			Item: &models.Item{
				Type:       models.DataTypeFloat,
				RoundFloat: true,
			},
		}

		h := &event.Prop{
			Code: "humidity",
			Name: "Влажность",
			Item: &models.Item{
				Type:       models.DataTypeFloat,
				RoundFloat: true,
			},
		}

		p := &event.Prop{
			Code: "pressure",
			Name: "Давление",
			Item: &models.Item{
				Type:       models.DataTypeFloat,
				RoundFloat: true,
			},
		}

		i := &event.Prop{
			Code: "illumination",
			Name: "Освещенность",
			Item: &models.Item{
				Type:       models.DataTypeFloat,
				RoundFloat: true,
			},
		}

		c := &event.Prop{
			Code: "current",
			Name: "Ток",
			Item: &models.Item{
				Type:       models.DataTypeFloat,
				RoundFloat: true,
			},
		}

		v := &event.Prop{
			Code: "voltage",
			Name: "Напряжение",
			Item: &models.Item{
				Type:       models.DataTypeFloat,
				RoundFloat: true,
			},
		}

		c2 := &event.Prop{
			Code: "co2",
			Name: "Углекислый газ",
			Item: &models.Item{
				Type:       models.DataTypeFloat,
				RoundFloat: true,
			},
		}

		m := &event.Prop{
			Code: "motion",
			Name: "Движение",
			Item: &models.Item{
				Type:       models.DataTypeFloat,
				RoundFloat: true,
			},
		}

		pr := &event.Prop{
			Code: "presence",
			Name: "Присутствие",
			Item: &models.Item{
				Type:       models.DataTypeFloat,
				RoundFloat: true,
			},
		}

		if err := e.Props.Add(t, h, p, i, c, v, c2, m, pr); err != nil {
			return nil, errors.Wrap(err, "init.maker")
		}

		return e, nil
	}

	// Для регистрации событий надо в service/init.go добавить импорт данного _пакета_!
	if err := event.Register(maker); err != nil {
		panic(err)
	}
}

func NewOnCheckMessage(topic string, targetID int, values map[string]float32) (messages.Message, error) {
	payload := make(map[string]interface{}, len(values))
	for k, v := range values {
		payload[k] = v
	}

	e, err := event.MakeEvent("object.sensor.on_check", messages.TargetTypeObject, targetID, payload)
	if err != nil {
		return nil, errors.Wrap(err, "NewOnCheckMessage")
	}

	m, err := e.ToMqttMessage(topic)
	if err != nil {
		return nil, errors.Wrap(err, "NewOnCheckMessage")
	}

	return m, nil
}
