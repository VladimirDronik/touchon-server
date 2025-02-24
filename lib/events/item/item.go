package item

import (
	"github.com/pkg/errors"
	"touchon-server/lib/event"
	"touchon-server/lib/interfaces"
	"touchon-server/lib/messages"
)

func NewOnChangeStateOn(itemID int) (interfaces.Event, error) {
	msg, err := messages.NewEvent("item.on_change_state_on", interfaces.TargetTypeItem, itemID)
	if err != nil {
		return nil, errors.Wrap(err, "NewOnChangeStateOn")
	}

	o := &OnChangeStateOn{
		Event: &event.EventImpl{
			Message:     msg,
			Title:       "on_change_state_on",
			Description: "Состояние элемента 'Вкл'",
		},
	}

	return o, nil
}

type OnChangeStateOn struct {
	interfaces.Event
}

func NewOnChangeStateOff(itemID int) (interfaces.Event, error) {
	msg, err := messages.NewEvent("item.on_change_state_off", interfaces.TargetTypeItem, itemID)
	if err != nil {
		return nil, errors.Wrap(err, "NewOnChangeStateOff")
	}

	o := &OnChangeStateOff{
		Event: &event.EventImpl{
			Message:     msg,
			Title:       "on_change_state_off",
			Description: "Состояние элемента 'Выкл'",
		},
	}

	return o, nil
}

type OnChangeStateOff struct {
	interfaces.Event
}
