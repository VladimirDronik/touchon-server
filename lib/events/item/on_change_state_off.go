package item

import (
	"github.com/pkg/errors"
	"touchon-server/lib/event"
	"touchon-server/lib/mqtt/messages"
)

func init() {
	maker := func() (*event.Event, error) {
		e := &event.Event{
			Code:        "item.on_change_state_off",
			Name:        "Состояние элемента 'Выкл'",
			Description: "",
			Props:       event.NewProps(),
			TargetType:  messages.TargetTypeItem,
		}

		return e, nil
	}

	// Для регистрации событий надо в service/init.go добавить импорт данного _пакета_!
	if err := event.Register(maker); err != nil {
		panic(err)
	}
}

func NewOnChangeStateOffMessage(topic string, targetID int) (messages.Message, error) {
	e, err := event.MakeEvent("item.on_change_state_off", messages.TargetTypeItem, targetID, nil)
	if err != nil {
		return nil, errors.Wrap(err, "NewOnChangeStateOffMessage")
	}

	m, err := e.ToMqttMessage(topic)
	if err != nil {
		return nil, errors.Wrap(err, "NewOnChangeStateOffMessage")
	}

	return m, nil
}
