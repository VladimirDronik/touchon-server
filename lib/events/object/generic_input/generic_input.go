package generic_input

import (
	"github.com/pkg/errors"
	"touchon-server/lib/event"
	"touchon-server/lib/interfaces"
	"touchon-server/lib/messages"
)

func NewOnClick(targetID int) (interfaces.Event, error) {
	msg, err := messages.NewEvent(interfaces.TargetTypeObject, targetID)
	if err != nil {
		return nil, errors.Wrap(err, "NewOnClick")
	}

	o := &OnClick{
		Event: &event.EventImpl{
			Message:          msg,
			EventCode:        "object.generic_input.on_click",
			EventName:        "on_click",
			EventDescription: "Одиночное замыкание",
		},
	}

	return o, nil
}

// OnClick Одиночное замыкание
type OnClick struct {
	interfaces.Event
}

func NewOnDoubleClick(targetID int) (interfaces.Event, error) {
	msg, err := messages.NewEvent(interfaces.TargetTypeObject, targetID)
	if err != nil {
		return nil, errors.Wrap(err, "NewOnDoubleClick")
	}

	o := &OnDoubleClick{
		Event: &event.EventImpl{
			Message:          msg,
			EventCode:        "object.generic_input.on_double_click",
			EventName:        "on_double_click",
			EventDescription: "Двойное замыкание",
		},
	}

	return o, nil
}

// OnDoubleClick Двойное замыкание
type OnDoubleClick struct {
	interfaces.Event
}

func NewOnLongPress(targetID int) (interfaces.Event, error) {
	msg, err := messages.NewEvent(interfaces.TargetTypeObject, targetID)
	if err != nil {
		return nil, errors.Wrap(err, "NewOnLongPress")
	}

	o := &OnLongPress{
		Event: &event.EventImpl{
			Message:          msg,
			EventCode:        "object.generic_input.on_long_press",
			EventName:        "on_long_press",
			EventDescription: "Длительное замыкание",
		},
	}

	return o, nil
}

// OnLongPress Длительное замыкание
type OnLongPress struct {
	interfaces.Event
}
