package relay

import (
	"github.com/pkg/errors"
	"touchon-server/lib/event"
	"touchon-server/lib/interfaces"
	"touchon-server/lib/messages"
)

func NewOnCheck(targetID int, state, value string) (interfaces.Event, error) {
	msg, err := messages.NewEvent("object.relay.on_check", interfaces.TargetTypeObject, targetID)
	if err != nil {
		return nil, errors.Wrap(err, "NewOnCheck")
	}

	o := &OnCheck{
		Event: &event.EventImpl{
			Message:     msg,
			Title:       "on_check",
			Description: "Проверка состояния реле",
		},
	}

	o.SetValue("state", state) // Состояние
	o.SetValue("value", value) // Значение

	return o, nil
}

// OnCheck Проверка состояния реле
type OnCheck struct {
	interfaces.Event
}

func NewOnStateOn(targetID int) (interfaces.Event, error) {
	msg, err := messages.NewEvent("object.relay.on_state_on", interfaces.TargetTypeObject, targetID)
	if err != nil {
		return nil, errors.Wrap(err, "NewOnStateOn")
	}

	o := &OnStateOn{
		Event: &event.EventImpl{
			Message:     msg,
			Title:       "on_state_on",
			Description: "Включение реле",
		},
	}

	return o, nil
}

// OnStateOn Включение реле
type OnStateOn struct {
	interfaces.Event
}

func NewOnStateOff(targetID int) (interfaces.Event, error) {
	msg, err := messages.NewEvent("object.relay.on_state_off", interfaces.TargetTypeObject, targetID)
	if err != nil {
		return nil, errors.Wrap(err, "NewOnStateOff")
	}

	o := &OnStateOff{
		Event: &event.EventImpl{
			Message:     msg,
			Title:       "on_state_off",
			Description: "Выключение реле",
		},
	}

	return o, nil
}

// OnStateOff Выключение реле
type OnStateOff struct {
	interfaces.Event
}
