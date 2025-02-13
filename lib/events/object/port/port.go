package port

import (
	"github.com/pkg/errors"
	"touchon-server/lib/event"
	"touchon-server/lib/interfaces"
	"touchon-server/lib/messages"
)

func NewOnPress(targetID int) (interfaces.Event, error) {
	msg, err := messages.NewEvent("object.port.on_press", interfaces.TargetTypeObject, targetID)
	if err != nil {
		return nil, errors.Wrap(err, "NewOnPress")
	}

	o := &OnPress{
		Event: &event.EventImpl{
			Message:     msg,
			Title:       "on_press",
			Description: "Порт замкнут",
		},
	}

	return o, nil
}

// OnPress Порт замкнут
type OnPress struct {
	interfaces.Event
}

func NewOnLongPress(targetID int) (interfaces.Event, error) {
	msg, err := messages.NewEvent("object.port.on_long_press", interfaces.TargetTypeObject, targetID)
	if err != nil {
		return nil, errors.Wrap(err, "NewOnLongPress")
	}

	o := &OnLongPress{
		Event: &event.EventImpl{
			Message:     msg,
			Title:       "on_long_press",
			Description: "Порт удерживается в замкнутом состоянии",
		},
	}

	return o, nil
}

// OnLongPress Порт удерживается в замкнутом состоянии
type OnLongPress struct {
	interfaces.Event
}

func NewOnRelease(targetID int) (interfaces.Event, error) {
	msg, err := messages.NewEvent("object.port.on_release", interfaces.TargetTypeObject, targetID)
	if err != nil {
		return nil, errors.Wrap(err, "NewOnRelease")
	}

	o := &OnRelease{
		Event: &event.EventImpl{
			Message:     msg,
			Title:       "on_release",
			Description: "Порт разомкнут",
		},
	}

	return o, nil
}

// OnRelease Порт разомкнут
type OnRelease struct {
	interfaces.Event
}

func NewOnDoubleClick(targetID int) (interfaces.Event, error) {
	msg, err := messages.NewEvent("object.port.on_double_click", interfaces.TargetTypeObject, targetID)
	if err != nil {
		return nil, errors.Wrap(err, "NewOnDoubleClick")
	}

	o := &OnDoubleClick{
		Event: &event.EventImpl{
			Message:     msg,
			Title:       "on_double_click",
			Description: "Порт замкнут дважды",
		},
	}

	return o, nil
}

// OnDoubleClick Двойное замыкание
type OnDoubleClick struct {
	interfaces.Event
}

func NewOnCheck(targetID int, state, value string) (interfaces.Event, error) {
	msg, err := messages.NewEvent("object.port.on_check", interfaces.TargetTypeObject, targetID)
	if err != nil {
		return nil, errors.Wrap(err, "NewOnCheck")
	}

	o := &OnCheck{
		Event: &event.EventImpl{
			Message:     msg,
			Title:       "on_check",
			Description: "Событие возникает, когда проверяется состояние порта, но при этом новое пришедшее состояние порта не различается с тем, что хранится в БД",
		},
	}

	o.SetValue("state", state) // Состояние
	o.SetValue("value", value) // Значение

	return o, nil
}

// OnCheck Событие возникает, когда проверяется состояние порта,
// но при этом новое пришедшее состояние порта не различается с тем, что хранится в БД
type OnCheck struct {
	interfaces.Event
}
