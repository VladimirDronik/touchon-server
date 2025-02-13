package relay

import (
	"github.com/pkg/errors"
	"touchon-server/lib/event"
	"touchon-server/lib/interfaces"
	"touchon-server/lib/messages"
)

func NewOnCheck(targetID int, state, value string) (interfaces.Event, error) {
	msg, err := messages.NewEvent(interfaces.TargetTypeObject, targetID)
	if err != nil {
		return nil, errors.Wrap(err, "NewOnCheck")
	}

	o := &OnCheck{
		Event: &event.EventImpl{
			Message:          msg,
			EventCode:        "object.relay.on_check",
			EventName:        "on_check",
			EventDescription: "Проверка состояния реле",
		},
		State: state,
		Value: value,
	}

	return o, nil
}

// OnCheck Проверка состояния реле
type OnCheck struct {
	interfaces.Event
	State string `json:"state,omitempty"` // Состояние
	Value string `json:"value,omitempty"` // Значение
}

func NewOnStateOn(targetID int) (interfaces.Event, error) {
	msg, err := messages.NewEvent(interfaces.TargetTypeObject, targetID)
	if err != nil {
		return nil, errors.Wrap(err, "NewOnStateOn")
	}

	o := &OnStateOn{
		Event: &event.EventImpl{
			Message:          msg,
			EventCode:        "object.relay.on_state_on",
			EventName:        "on_state_on",
			EventDescription: "Включение реле",
		},
	}

	return o, nil
}

// OnStateOn Включение реле
type OnStateOn struct {
	interfaces.Event
}

func NewOnStateOff(targetID int) (interfaces.Event, error) {
	msg, err := messages.NewEvent(interfaces.TargetTypeObject, targetID)
	if err != nil {
		return nil, errors.Wrap(err, "NewOnStateOff")
	}

	o := &OnStateOff{
		Event: &event.EventImpl{
			Message:          msg,
			EventCode:        "object.relay.on_state_off",
			EventName:        "on_state_off",
			EventDescription: "Выключение реле",
		},
	}

	return o, nil
}

// OnStateOff Выключение реле
type OnStateOff struct {
	interfaces.Event
}
