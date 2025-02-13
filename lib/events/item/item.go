package item

import (
	"github.com/pkg/errors"
	"touchon-server/lib/event"
	"touchon-server/lib/interfaces"
	"touchon-server/lib/messages"
)

func NewOnChangeStateOn(itemID int) (interfaces.Event, error) {
	msg, err := messages.NewEvent(interfaces.TargetTypeItem, itemID)
	if err != nil {
		return nil, errors.Wrap(err, "NewOnChangeStateOn")
	}

	o := &OnChangeStateOn{
		Event: &event.EventImpl{
			Message:          msg,
			EventCode:        "item.on_change_state_on",
			EventName:        "on_change_state_on",
			EventDescription: "Состояние элемента 'Вкл'",
		},
	}

	return o, nil
}

type OnChangeStateOn struct {
	interfaces.Event
}

func NewOnChangeStateOff(itemID int) (interfaces.Event, error) {
	msg, err := messages.NewEvent(interfaces.TargetTypeItem, itemID)
	if err != nil {
		return nil, errors.Wrap(err, "NewOnChangeStateOff")
	}

	o := &OnChangeStateOff{
		Event: &event.EventImpl{
			Message:          msg,
			EventCode:        "item.on_change_state_off",
			EventName:        "on_change_state_off",
			EventDescription: "Состояние элемента 'Выкл'",
		},
	}

	return o, nil
}

type OnChangeStateOff struct {
	interfaces.Event
}
