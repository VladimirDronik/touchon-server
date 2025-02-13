package sensor

import (
	"github.com/pkg/errors"
	"touchon-server/lib/event"
	"touchon-server/lib/interfaces"
	"touchon-server/lib/messages"
)

func NewOnAlarm(targetID int, msgText string) (interfaces.Event, error) {
	msg, err := messages.NewEvent(interfaces.TargetTypeObject, targetID)
	if err != nil {
		return nil, errors.Wrap(err, "NewOnAlarm")
	}

	o := &OnAlarm{
		Event: &event.EventImpl{
			Message:          msg,
			EventCode:        "object.sensor.on_alarm",
			EventName:        "on_alarm",
			EventDescription: "Данные датчика вышли за пороговые значения",
		},
		MsgText: msgText,
	}

	return o, nil
}

// OnAlarm Данные датчика вышли за пороговые значения
type OnAlarm struct {
	interfaces.Event
	MsgText string `json:"message,omitempty"` // Сообщение
}

func NewOnCheck(targetID int, values map[string]float32) (interfaces.Event, error) {
	msg, err := messages.NewEvent(interfaces.TargetTypeObject, targetID)
	if err != nil {
		return nil, errors.Wrap(err, "NewOnCheck")
	}

	o := &OnCheck{
		Event: &event.EventImpl{
			Message:          msg,
			EventCode:        "object.sensor.on_check",
			EventName:        "on_check",
			EventDescription: "Данные датчика обновлены",
		},
		Values: values,
	}

	return o, nil
}

// OnCheck Данные датчика обновлены
type OnCheck struct {
	interfaces.Event
	Values map[string]float32 `json:"values,omitempty"` // Значения
}

func NewOnMotionOn(targetID int) (interfaces.Event, error) {
	msg, err := messages.NewEvent(interfaces.TargetTypeObject, targetID)
	if err != nil {
		return nil, errors.Wrap(err, "NewOnMotionOn")
	}

	o := &OnMotionOn{
		Event: &event.EventImpl{
			Message:          msg,
			EventCode:        "object.sensor.on_motion_on",
			EventName:        "on_motion_on",
			EventDescription: "Движение есть",
		},
	}

	return o, nil
}

// OnMotionOn Движение есть
type OnMotionOn struct {
	interfaces.Event
}

func NewOnMotionOff(targetID int) (interfaces.Event, error) {
	msg, err := messages.NewEvent(interfaces.TargetTypeObject, targetID)
	if err != nil {
		return nil, errors.Wrap(err, "NewOnMotionOff")
	}

	o := &OnMotionOff{
		Event: &event.EventImpl{
			Message:          msg,
			EventCode:        "object.sensor.on_motion_off",
			EventName:        "on_motion_off",
			EventDescription: "Движения нет",
		},
	}

	return o, nil
}

// OnMotionOff Движения нет
type OnMotionOff struct {
	interfaces.Event
}

func NewOnPresenceOn(targetID int) (interfaces.Event, error) {
	msg, err := messages.NewEvent(interfaces.TargetTypeObject, targetID)
	if err != nil {
		return nil, errors.Wrap(err, "NewOnPresenceOn")
	}

	o := &OnPresenceOn{
		Event: &event.EventImpl{
			Message:          msg,
			EventCode:        "object.sensor.on_presence_on",
			EventName:        "on_presence_on",
			EventDescription: "Присутствие есть",
		},
	}

	return o, nil
}

// OnPresenceOn Присутствие есть
type OnPresenceOn struct {
	interfaces.Event
}

func NewOnPresenceOff(targetID int) (interfaces.Event, error) {
	msg, err := messages.NewEvent(interfaces.TargetTypeObject, targetID)
	if err != nil {
		return nil, errors.Wrap(err, "NewOnPresenceOff")
	}

	o := &OnPresenceOff{
		Event: &event.EventImpl{
			Message:          msg,
			EventCode:        "object.sensor.on_presence_off",
			EventName:        "on_presence_off",
			EventDescription: "Присутствия нет",
		},
	}

	return o, nil
}

// OnPresenceOff Присутствия нет
type OnPresenceOff struct {
	interfaces.Event
}
