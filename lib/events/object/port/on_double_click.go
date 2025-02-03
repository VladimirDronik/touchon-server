package port

import (
	"github.com/pkg/errors"
	"touchon-server/lib/event"
	"touchon-server/lib/mqtt/messages"
)

func init() {
	maker := func() (*event.Event, error) {
		e := &event.Event{
			Code:        "object.port.on_double_click",
			Name:        "on_double_click",
			Description: "Порт замкнут дважды",
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

func NewOnDoubleClickMessage(topic string, targetID int) (messages.Message, error) {
	e, err := event.MakeEvent("object.port.on_double_click", messages.TargetTypeObject, targetID, nil)
	if err != nil {
		return nil, errors.Wrap(err, "NewOnDoubleClickMessage")
	}

	m, err := e.ToMqttMessage(topic)
	if err != nil {
		return nil, errors.Wrap(err, "NewOnDoubleClickMessage")
	}

	return m, nil
}
