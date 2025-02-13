package sensor

import (
	"github.com/pkg/errors"
	"touchon-server/lib/event"
	"touchon-server/lib/interfaces"
	"touchon-server/lib/messages"
)

func NewOnAlarm(targetID int, message string) (interfaces.Event, error) {
	msg, err := messages.NewEvent("object.sensor.on_alarm", interfaces.TargetTypeObject, targetID)
	if err != nil {
		return nil, errors.Wrap(err, "NewOnAlarm")
	}

	o := &OnAlarm{
		Event: &event.EventImpl{
			Message:     msg,
			Title:       "on_alarm",
			Description: "Данные датчика вышли за пороговые значения",
		},
	}

	o.SetValue("message", message) // Сообщение

	return o, nil
}

// OnAlarm Данные датчика вышли за пороговые значения
type OnAlarm struct {
	interfaces.Event
}

func NewOnCheck(targetID int, values map[string]float32) (interfaces.Event, error) {
	msg, err := messages.NewEvent("object.sensor.on_check", interfaces.TargetTypeObject, targetID)
	if err != nil {
		return nil, errors.Wrap(err, "NewOnCheck")
	}

	o := &OnCheck{
		Event: &event.EventImpl{
			Message:     msg,
			Title:       "on_check",
			Description: "Данные датчика обновлены",
		},
	}

	for k, v := range values {
		o.SetValue(k, v) // Значения
	}

	return o, nil
}

// OnCheck Данные датчика обновлены
type OnCheck struct {
	interfaces.Event
}

func NewOnMotionOn(targetID int) (interfaces.Event, error) {
	msg, err := messages.NewEvent("object.sensor.on_motion_on", interfaces.TargetTypeObject, targetID)
	if err != nil {
		return nil, errors.Wrap(err, "NewOnMotionOn")
	}

	o := &OnMotionOn{
		Event: &event.EventImpl{
			Message:     msg,
			Title:       "on_motion_on",
			Description: "Движение есть",
		},
	}

	return o, nil
}

// OnMotionOn Движение есть
type OnMotionOn struct {
	interfaces.Event
}

func NewOnMotionOff(targetID int) (interfaces.Event, error) {
	msg, err := messages.NewEvent("object.sensor.on_motion_off", interfaces.TargetTypeObject, targetID)
	if err != nil {
		return nil, errors.Wrap(err, "NewOnMotionOff")
	}

	o := &OnMotionOff{
		Event: &event.EventImpl{
			Message:     msg,
			Title:       "on_motion_off",
			Description: "Движения нет",
		},
	}

	return o, nil
}

// OnMotionOff Движения нет
type OnMotionOff struct {
	interfaces.Event
}

func NewOnPresenceOn(targetID int) (interfaces.Event, error) {
	msg, err := messages.NewEvent("object.sensor.on_presence_on", interfaces.TargetTypeObject, targetID)
	if err != nil {
		return nil, errors.Wrap(err, "NewOnPresenceOn")
	}

	o := &OnPresenceOn{
		Event: &event.EventImpl{
			Message:     msg,
			Title:       "on_presence_on",
			Description: "Присутствие есть",
		},
	}

	return o, nil
}

// OnPresenceOn Присутствие есть
type OnPresenceOn struct {
	interfaces.Event
}

func NewOnPresenceOff(targetID int) (interfaces.Event, error) {
	msg, err := messages.NewEvent("object.sensor.on_presence_off", interfaces.TargetTypeObject, targetID)
	if err != nil {
		return nil, errors.Wrap(err, "NewOnPresenceOff")
	}

	o := &OnPresenceOff{
		Event: &event.EventImpl{
			Message:     msg,
			Title:       "on_presence_off",
			Description: "Присутствия нет",
		},
	}

	return o, nil
}

// OnPresenceOff Присутствия нет
type OnPresenceOff struct {
	interfaces.Event
}
