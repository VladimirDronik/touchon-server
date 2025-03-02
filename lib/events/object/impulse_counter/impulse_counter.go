package impulse_counter

import (
	"github.com/pkg/errors"
	"touchon-server/lib/event"
	"touchon-server/lib/interfaces"
	"touchon-server/lib/messages"
)

func NewOnThreshold(targetID, countImpulse int) (interfaces.Event, error) {
	msg, err := messages.NewEvent("object.impulse_counter.on_threshold", interfaces.TargetTypeObject, targetID)
	if err != nil {
		return nil, errors.Wrap(err, "NewOnThreshold")
	}

	o := &OnThreshold{
		Event: &event.EventImpl{
			Message:     msg,
			Title:       "on_threshold",
			Description: "Достигнуто пороговое значение счетчика",
		},
	}

	o.SetValue("count_impulse", countImpulse)

	return o, nil
}

// OnThreshold Достигнуто пороговое значение счетчика
type OnThreshold struct {
	interfaces.Event
}

func NewOnCheck(targetID, countImpulse int) (interfaces.Event, error) {
	msg, err := messages.NewEvent("object.impulse_counter.on_check", interfaces.TargetTypeObject, targetID)
	if err != nil {
		return nil, errors.Wrap(err, "NewOnCheck")
	}

	o := &OnCheck{
		Event: &event.EventImpl{
			Message:     msg,
			Title:       "on_check",
			Description: "Произведено снятие количества импульсов со счетчика",
		},
	}

	o.SetValue("count_impulse", countImpulse)

	return o, nil
}

// OnCheck Произведено снятие количества импульсов со счетчика
type OnCheck struct {
	interfaces.Event
}
