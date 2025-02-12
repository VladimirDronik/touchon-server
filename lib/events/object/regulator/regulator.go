package regulator

import (
	"github.com/pkg/errors"
	"touchon-server/lib/interfaces"
	"touchon-server/lib/messages"
)

func NewOnAbove(targetID int, value float32) (interfaces.Message, error) {
	msg, err := messages.NewEvent(interfaces.TargetTypeObject, targetID)
	if err != nil {
		return nil, errors.Wrap(err, "NewOnAbove")
	}

	o := &OnAbove{
		MessageImpl: msg,
		Value:       value,
	}

	return o, nil
}

// OnAbove Текущее значение больше заданного
type OnAbove struct {
	*messages.MessageImpl
	Value float32 `json:"value"` // Значение
}

func NewOnBelow(targetID int, value float32) (interfaces.Message, error) {
	msg, err := messages.NewEvent(interfaces.TargetTypeObject, targetID)
	if err != nil {
		return nil, errors.Wrap(err, "NewOnBelow")
	}

	o := &OnBelow{
		MessageImpl: msg,
		Value:       value,
	}

	return o, nil
}

// OnBelow Текущее значение меньше заданного
type OnBelow struct {
	*messages.MessageImpl
	Value float32 `json:"value"` // Значение
}

func NewOnStale(targetID int) (interfaces.Message, error) {
	msg, err := messages.NewEvent(interfaces.TargetTypeObject, targetID)
	if err != nil {
		return nil, errors.Wrap(err, "NewOnStale")
	}

	o := &OnStale{
		MessageImpl: msg,
	}

	return o, nil
}

// OnStale Текущее значение не актуально
type OnStale struct {
	*messages.MessageImpl
}

func NewOnComplexBelow1(targetID int, value float32) (interfaces.Message, error) {
	msg, err := messages.NewEvent(interfaces.TargetTypeObject, targetID)
	if err != nil {
		return nil, errors.Wrap(err, "NewOnComplexBelow1")
	}

	o := &OnComplexBelow1{
		MessageImpl: msg,
		Value:       value,
	}

	return o, nil
}

// OnComplexBelow1 Текущее значение < (targetSP - complexTolerance - belowTolerance)
type OnComplexBelow1 struct {
	*messages.MessageImpl
	Value float32 `json:"value"` // Значение
}

func NewOnComplexBelow2(targetID int, value float32) (interfaces.Message, error) {
	msg, err := messages.NewEvent(interfaces.TargetTypeObject, targetID)
	if err != nil {
		return nil, errors.Wrap(err, "NewOnComplexBelow2")
	}

	o := &OnComplexBelow2{
		MessageImpl: msg,
		Value:       value,
	}

	return o, nil
}

// OnComplexBelow2 Текущее значение > (targetSP + complexTolerance - belowTolerance)
type OnComplexBelow2 struct {
	*messages.MessageImpl
	Value float32 `json:"value"` // Значение
}

func NewOnComplexAbove1(targetID int, value float32) (interfaces.Message, error) {
	msg, err := messages.NewEvent(interfaces.TargetTypeObject, targetID)
	if err != nil {
		return nil, errors.Wrap(err, "NewOnComplexAbove1")
	}

	o := &OnComplexAbove1{
		MessageImpl: msg,
		Value:       value,
	}

	return o, nil
}

// OnComplexAbove1 Текущее значение < (targetSP - complexTolerance + aboveTolerance)
type OnComplexAbove1 struct {
	*messages.MessageImpl
	Value float32 `json:"value"` // Значение
}

func NewOnComplexAbove2(targetID int, value float32) (interfaces.Message, error) {
	msg, err := messages.NewEvent(interfaces.TargetTypeObject, targetID)
	if err != nil {
		return nil, errors.Wrap(err, "NewOnComplexAbove2")
	}

	o := &OnComplexAbove2{
		MessageImpl: msg,
		Value:       value,
	}

	return o, nil
}

// OnComplexAbove2 Текущее значение > (targetSP + complexTolerance + aboveTolerance)
type OnComplexAbove2 struct {
	*messages.MessageImpl
	Value float32 `json:"value"` // Значение
}
