package regulator

import (
	"github.com/pkg/errors"
	"touchon-server/lib/event"
	"touchon-server/lib/interfaces"
	"touchon-server/lib/messages"
)

func NewOnAbove(targetID int, value float32) (interfaces.Event, error) {
	msg, err := messages.NewEvent(interfaces.TargetTypeObject, targetID)
	if err != nil {
		return nil, errors.Wrap(err, "NewOnAbove")
	}

	o := &OnAbove{
		Event: &event.EventImpl{
			Message:          msg,
			EventCode:        "object.regulator.on_above",
			EventName:        "on_above",
			EventDescription: "Текущее значение больше заданного",
		},
		Value: value,
	}

	return o, nil
}

// OnAbove Текущее значение больше заданного
type OnAbove struct {
	interfaces.Event
	Value float32 `json:"value"` // Значение
}

func NewOnBelow(targetID int, value float32) (interfaces.Event, error) {
	msg, err := messages.NewEvent(interfaces.TargetTypeObject, targetID)
	if err != nil {
		return nil, errors.Wrap(err, "NewOnBelow")
	}

	o := &OnBelow{
		Event: &event.EventImpl{
			Message:          msg,
			EventCode:        "object.regulator.on_below",
			EventName:        "on_below",
			EventDescription: "Текущее значение меньше заданного",
		},
		Value: value,
	}

	return o, nil
}

// OnBelow Текущее значение меньше заданного
type OnBelow struct {
	interfaces.Event
	Value float32 `json:"value"` // Значение
}

func NewOnStale(targetID int) (interfaces.Event, error) {
	msg, err := messages.NewEvent(interfaces.TargetTypeObject, targetID)
	if err != nil {
		return nil, errors.Wrap(err, "NewOnStale")
	}

	o := &OnStale{
		Event: &event.EventImpl{
			Message:          msg,
			EventCode:        "object.regulator.on_stale",
			EventName:        "on_stale",
			EventDescription: "Текущее значение не актуально",
		},
	}

	return o, nil
}

// OnStale Текущее значение не актуально
type OnStale struct {
	interfaces.Event
}

func NewOnComplexBelow1(targetID int, value float32) (interfaces.Event, error) {
	msg, err := messages.NewEvent(interfaces.TargetTypeObject, targetID)
	if err != nil {
		return nil, errors.Wrap(err, "NewOnComplexBelow1")
	}

	o := &OnComplexBelow1{
		Event: &event.EventImpl{
			Message:          msg,
			EventCode:        "object.regulator.on_complex_below_1",
			EventName:        "on_complex_below_1",
			EventDescription: "Текущее значение < (targetSP - complexTolerance - belowTolerance)",
		},
		Value: value,
	}

	return o, nil
}

// OnComplexBelow1 Текущее значение < (targetSP - complexTolerance - belowTolerance)
type OnComplexBelow1 struct {
	interfaces.Event
	Value float32 `json:"value"` // Значение
}

func NewOnComplexBelow2(targetID int, value float32) (interfaces.Event, error) {
	msg, err := messages.NewEvent(interfaces.TargetTypeObject, targetID)
	if err != nil {
		return nil, errors.Wrap(err, "NewOnComplexBelow2")
	}

	o := &OnComplexBelow2{
		Event: &event.EventImpl{
			Message:          msg,
			EventCode:        "object.regulator.on_complex_below_2",
			EventName:        "on_complex_below_2",
			EventDescription: "Текущее значение > (targetSP + complexTolerance - belowTolerance)",
		},
		Value: value,
	}

	return o, nil
}

// OnComplexBelow2 Текущее значение > (targetSP + complexTolerance - belowTolerance)
type OnComplexBelow2 struct {
	interfaces.Event
	Value float32 `json:"value"` // Значение
}

func NewOnComplexAbove1(targetID int, value float32) (interfaces.Event, error) {
	msg, err := messages.NewEvent(interfaces.TargetTypeObject, targetID)
	if err != nil {
		return nil, errors.Wrap(err, "NewOnComplexAbove1")
	}

	o := &OnComplexAbove1{
		Event: &event.EventImpl{
			Message:          msg,
			EventCode:        "object.regulator.on_complex_above_1",
			EventName:        "on_complex_above_1",
			EventDescription: "Текущее значение < (targetSP - complexTolerance + aboveTolerance)",
		},
		Value: value,
	}

	return o, nil
}

// OnComplexAbove1 Текущее значение < (targetSP - complexTolerance + aboveTolerance)
type OnComplexAbove1 struct {
	interfaces.Event
	Value float32 `json:"value"` // Значение
}

func NewOnComplexAbove2(targetID int, value float32) (interfaces.Event, error) {
	msg, err := messages.NewEvent(interfaces.TargetTypeObject, targetID)
	if err != nil {
		return nil, errors.Wrap(err, "NewOnComplexAbove2")
	}

	o := &OnComplexAbove2{
		Event: &event.EventImpl{
			Message:          msg,
			EventCode:        "object.regulator.on_complex_above_2",
			EventName:        "on_complex_above_2",
			EventDescription: "Текущее значение > (targetSP + complexTolerance + aboveTolerance)",
		},
		Value: value,
	}

	return o, nil
}

// OnComplexAbove2 Текущее значение > (targetSP + complexTolerance + aboveTolerance)
type OnComplexAbove2 struct {
	interfaces.Event
	Value float32 `json:"value"` // Значение
}
