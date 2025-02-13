package regulator

import (
	"github.com/pkg/errors"
	"touchon-server/lib/event"
	"touchon-server/lib/interfaces"
	"touchon-server/lib/messages"
)

func NewOnAbove(targetID int, value float32) (interfaces.Event, error) {
	msg, err := messages.NewEvent("object.regulator.on_above", interfaces.TargetTypeObject, targetID)
	if err != nil {
		return nil, errors.Wrap(err, "NewOnAbove")
	}

	o := &OnAbove{
		Event: &event.EventImpl{
			Message:     msg,
			Title:       "on_above",
			Description: "Текущее значение больше заданного",
		},
	}

	o.SetValue("value", value) // Значение

	return o, nil
}

// OnAbove Текущее значение больше заданного
type OnAbove struct {
	interfaces.Event
}

func NewOnBelow(targetID int, value float32) (interfaces.Event, error) {
	msg, err := messages.NewEvent("object.regulator.on_below", interfaces.TargetTypeObject, targetID)
	if err != nil {
		return nil, errors.Wrap(err, "NewOnBelow")
	}

	o := &OnBelow{
		Event: &event.EventImpl{
			Message:     msg,
			Title:       "on_below",
			Description: "Текущее значение меньше заданного",
		},
	}

	o.SetValue("value", value) // Значение

	return o, nil
}

// OnBelow Текущее значение меньше заданного
type OnBelow struct {
	interfaces.Event
}

func NewOnStale(targetID int) (interfaces.Event, error) {
	msg, err := messages.NewEvent("object.regulator.on_stale", interfaces.TargetTypeObject, targetID)
	if err != nil {
		return nil, errors.Wrap(err, "NewOnStale")
	}

	o := &OnStale{
		Event: &event.EventImpl{
			Message:     msg,
			Title:       "on_stale",
			Description: "Текущее значение не актуально",
		},
	}

	return o, nil
}

// OnStale Текущее значение не актуально
type OnStale struct {
	interfaces.Event
}

func NewOnComplexBelow1(targetID int, value float32) (interfaces.Event, error) {
	msg, err := messages.NewEvent("object.regulator.on_complex_below_1", interfaces.TargetTypeObject, targetID)
	if err != nil {
		return nil, errors.Wrap(err, "NewOnComplexBelow1")
	}

	o := &OnComplexBelow1{
		Event: &event.EventImpl{
			Message:     msg,
			Title:       "on_complex_below_1",
			Description: "Текущее значение < (targetSP - complexTolerance - belowTolerance)",
		},
	}

	o.SetValue("value", value) // Значение

	return o, nil
}

// OnComplexBelow1 Текущее значение < (targetSP - complexTolerance - belowTolerance)
type OnComplexBelow1 struct {
	interfaces.Event
}

func NewOnComplexBelow2(targetID int, value float32) (interfaces.Event, error) {
	msg, err := messages.NewEvent("object.regulator.on_complex_below_2", interfaces.TargetTypeObject, targetID)
	if err != nil {
		return nil, errors.Wrap(err, "NewOnComplexBelow2")
	}

	o := &OnComplexBelow2{
		Event: &event.EventImpl{
			Message:     msg,
			Title:       "on_complex_below_2",
			Description: "Текущее значение > (targetSP + complexTolerance - belowTolerance)",
		},
	}

	o.SetValue("value", value) // Значение

	return o, nil
}

// OnComplexBelow2 Текущее значение > (targetSP + complexTolerance - belowTolerance)
type OnComplexBelow2 struct {
	interfaces.Event
}

func NewOnComplexAbove1(targetID int, value float32) (interfaces.Event, error) {
	msg, err := messages.NewEvent("object.regulator.on_complex_above_1", interfaces.TargetTypeObject, targetID)
	if err != nil {
		return nil, errors.Wrap(err, "NewOnComplexAbove1")
	}

	o := &OnComplexAbove1{
		Event: &event.EventImpl{
			Message:     msg,
			Title:       "on_complex_above_1",
			Description: "Текущее значение < (targetSP - complexTolerance + aboveTolerance)",
		},
	}

	o.SetValue("value", value) // Значение

	return o, nil
}

// OnComplexAbove1 Текущее значение < (targetSP - complexTolerance + aboveTolerance)
type OnComplexAbove1 struct {
	interfaces.Event
}

func NewOnComplexAbove2(targetID int, value float32) (interfaces.Event, error) {
	msg, err := messages.NewEvent("object.regulator.on_complex_above_2", interfaces.TargetTypeObject, targetID)
	if err != nil {
		return nil, errors.Wrap(err, "NewOnComplexAbove2")
	}

	o := &OnComplexAbove2{
		Event: &event.EventImpl{
			Message:     msg,
			Title:       "on_complex_above_2",
			Description: "Текущее значение > (targetSP + complexTolerance + aboveTolerance)",
		},
	}

	o.SetValue("value", value) // Значение

	return o, nil
}

// OnComplexAbove2 Текущее значение > (targetSP + complexTolerance + aboveTolerance)
type OnComplexAbove2 struct {
	interfaces.Event
}
