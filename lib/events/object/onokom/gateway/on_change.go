package gateway

import (
	"github.com/pkg/errors"
	"touchon-server/lib/event"
	"touchon-server/lib/mqtt/messages"
)

func init() {
	maker := func() (*event.Event, error) {
		e := &event.Event{
			Code:        "object.onokom.gateway.on_change",
			Name:        "on_change",
			Description: "Изменение состояния устройства",
			Props:       event.NewProps(),
			TargetType:  messages.TargetTypeObject,
		}

		if err := e.Props.Add(props...); err != nil {
			return nil, errors.Wrap(err, "init.maker")
		}

		return e, nil
	}

	// Для регистрации событий надо в service/init.go добавить импорт данного _пакета_!
	if err := event.Register(maker); err != nil {
		panic(err)
	}
}

func NewOnChangeMessage(topic string, targetID int, values map[string]interface{}) (messages.Message, error) {
	e, err := event.MakeEvent("object.onokom.gateway.on_change", messages.TargetTypeObject, targetID, values)
	if err != nil {
		return nil, errors.Wrap(err, "NewOnChangeMessage")
	}

	m, err := e.ToMqttMessage(topic)
	if err != nil {
		return nil, errors.Wrap(err, "NewOnChangeMessage")
	}

	return m, nil
}
