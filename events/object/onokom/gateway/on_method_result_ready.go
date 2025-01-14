package gateway

import (
	"github.com/VladimirDronik/touchon-server/event"
	"github.com/VladimirDronik/touchon-server/models"
	"github.com/VladimirDronik/touchon-server/mqtt/messages"
	"github.com/pkg/errors"
)

func init() {
	maker := func() (*event.Event, error) {
		e := &event.Event{
			Code:        "object.onokom.gateway.on_method_result_ready",
			Name:        "on_method_result_ready",
			Description: "Результат выполнения метода",
			Props:       event.NewProps(),
			TargetType:  messages.TargetTypeObject,
		}

		m := &event.Prop{
			Code: "method",
			Name: "Название метода",
			Item: &models.Item{
				Type: models.DataTypeString,
			},
		}

		r := &event.Prop{
			Code: "result",
			Name: "Результат",
			Item: &models.Item{
				Type: models.DataTypeBool,
			},
		}

		if err := e.Props.Add(m, r); err != nil {
			return nil, errors.Wrap(err, "init.maker")
		}

		return e, nil
	}

	// Для регистрации событий надо в service/init.go добавить импорт данного _пакета_!
	if err := event.Register(maker); err != nil {
		panic(err)
	}
}

// NewOnMethodResultReadyMessage используется для возврата результат метода.
func NewOnMethodResultReadyMessage(topic string, targetID int, methodCodeName string, result bool) (messages.Message, error) {
	payload := map[string]interface{}{
		"method": methodCodeName,
		"result": result,
	}

	e, err := event.MakeEvent("object.onokom.gateway.on_method_result_ready", messages.TargetTypeObject, targetID, payload)
	if err != nil {
		return nil, errors.Wrap(err, "NewOnMethodResultReadyMessage")
	}

	m, err := e.ToMqttMessage(topic)
	if err != nil {
		return nil, errors.Wrap(err, "NewOnMethodResultReadyMessage")
	}

	return m, nil
}
