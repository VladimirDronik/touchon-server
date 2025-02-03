package relay

import (
	"github.com/pkg/errors"
	"touchon-server/lib/event"
	"touchon-server/lib/mqtt/messages"
)

func init() {
	maker := func() (*event.Event, error) {
		e := &event.Event{
			Code:        "object.relay.on_check",
			Name:        "on_check",
			Description: "Проверка состояния реле",
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

func NewCheckMessage(topic string, targetType messages.TargetType, targetID int, state, value string) (messages.Message, error) {
	e, err := event.MakeEvent("on_change_state", targetType, targetID, map[string]interface{}{"state": state, "value": value})
	if err != nil {
		return nil, errors.Wrap(err, "NewCheckMessage")
	}

	m, err := e.ToMqttMessage(topic)
	if err != nil {
		return nil, errors.Wrap(err, "NewCheckMessage")
	}

	return m, nil
}
