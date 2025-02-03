package generic_input

import (
	"github.com/pkg/errors"
	"touchon-server/lib/event"
	"touchon-server/lib/mqtt/messages"
)

func init() {
	maker := func() (*event.Event, error) {
		e := &event.Event{
			Code:        "object.generic_input.on_click",
			Name:        "on_click",
			Description: "Одиночное замыкание",
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

func NewOnClickMessage(topic string, targetID int) (messages.Message, error) {
	e, err := event.MakeEvent("object.generic_input.on_click", messages.TargetTypeObject, targetID, nil)
	if err != nil {
		return nil, errors.Wrap(err, "NewOnClickMessage")
	}

	m, err := e.ToMqttMessage(topic)
	if err != nil {
		return nil, errors.Wrap(err, "NewOnClickMessage")
	}

	return m, nil
}
