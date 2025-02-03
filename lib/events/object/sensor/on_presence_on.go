package sensor

import (
	"github.com/pkg/errors"
	"touchon-server/lib/event"
	"touchon-server/lib/mqtt/messages"
)

func init() {
	maker := func() (*event.Event, error) {
		e := &event.Event{
			Code:        "object.sensor.on_presence_on",
			Name:        "on_presence_on",
			Description: "Присутствие есть",
			Props:       event.NewProps(),
			TargetType:  messages.TargetTypeObject,
		}

		return e, nil
	}

	// Для регистрации событий надо в service/init.go добавить импорт данного _пакета_!
	if err := event.Register(maker); err != nil {
		panic(err)
	}
}

func NewOnPresenceOnMessage(topic string, targetID int) (messages.Message, error) {
	e, err := event.MakeEvent("object.sensor.on_presence_on", messages.TargetTypeObject, targetID, nil)
	if err != nil {
		return nil, errors.Wrap(err, "NewOnPresenceOnMessage")
	}

	m, err := e.ToMqttMessage(topic)
	if err != nil {
		return nil, errors.Wrap(err, "NewOnPresenceOnMessage")
	}

	return m, nil
}
