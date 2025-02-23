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
