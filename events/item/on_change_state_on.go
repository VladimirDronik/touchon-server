package item

import (
	"github.com/VladimirDronik/touchon-server/event"
	"github.com/VladimirDronik/touchon-server/mqtt/messages"
	"github.com/pkg/errors"
)

func init() {
	maker := func() (*event.Event, error) {
		e := &event.Event{
			Code:        "item.on_change_state_on",
			Name:        "Состояние элемента 'Вкл'",
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

func NewOnChangeStateOnMessage(topic string, targetID int) (messages.Message, error) {
	e, err := event.MakeEvent("item.on_change_state_on", messages.TargetTypeItem, targetID, nil)
	if err != nil {
		return nil, errors.Wrap(err, "NewOnChangeStateOnMessage")
	}

	m, err := e.ToMqttMessage(topic)
	if err != nil {
		return nil, errors.Wrap(err, "NewOnChangeStateOnMessage")
	}

	return m, nil
}
