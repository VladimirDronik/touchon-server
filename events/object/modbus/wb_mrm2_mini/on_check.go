package wb_mrm2_mini

import (
	"github.com/VladimirDronik/touchon-server/event"
	"github.com/VladimirDronik/touchon-server/models"
	"github.com/VladimirDronik/touchon-server/mqtt/messages"
	"github.com/pkg/errors"
)

func init() {
	maker := func() (*event.Event, error) {
		e := &event.Event{
			Code:        "object.modbus.wb_mrm2_mini.on_check",
			Name:        "on_check",
			Description: "Получено состояние выходов",
			Props:       event.NewProps(),
			TargetType:  messages.TargetTypeObject,
		}

		k1 := &event.Prop{
			Code: "k1",
			Name: "Выход 1",
			Item: &models.Item{
				Type: models.DataTypeBool,
			},
		}

		k2 := &event.Prop{
			Code: "k2",
			Name: "Выход 2",
			Item: &models.Item{
				Type: models.DataTypeBool,
			},
		}

		if err := e.Props.Add(k1, k2); err != nil {
			return nil, errors.Wrap(err, "init.maker")
		}

		return e, nil
	}

	// Для регистрации событий надо в service/init.go добавить импорт данного _пакета_!
	if err := event.Register(maker); err != nil {
		panic(err)
	}
}

// NewOnCheckMessage используется при проверке состояния как всех выходов, так и при проверке отдельных выходов.
func NewOnCheckMessage(topic string, targetID int, values map[string]bool) (messages.Message, error) {
	payload := make(map[string]interface{}, len(values))
	for _, k := range []string{"k1", "k2"} {
		if v, ok := values[k]; ok {
			payload[k] = v
		}
	}

	e, err := event.MakeEvent("object.modbus.wb_mrm2_mini.on_check", messages.TargetTypeObject, targetID, payload)
	if err != nil {
		return nil, errors.Wrap(err, "NewOnCheckMessage")
	}

	m, err := e.ToMqttMessage(topic)
	if err != nil {
		return nil, errors.Wrap(err, "NewOnCheckMessage")
	}

	return m, nil
}
