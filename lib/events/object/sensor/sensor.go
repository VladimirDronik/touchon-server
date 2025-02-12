package sensor

import (
	"github.com/pkg/errors"
	"touchon-server/lib/interfaces"
	"touchon-server/lib/messages"
)

func NewOnAlarm(targetID int, msgText string) (interfaces.Message, error) {
	msg, err := messages.NewEvent(interfaces.TargetTypeObject, targetID)
	if err != nil {
		return nil, errors.Wrap(err, "NewOnAlarm")
	}

	o := &OnAlarm{
		MessageImpl: msg,
		Message:     msgText,
	}

	return o, nil
}

// OnAlarm Данные датчика вышли за пороговые значения
type OnAlarm struct {
	*messages.MessageImpl
	Message string `json:"message"` // Сообщение
}

func NewOnCheck(targetID int, values map[string]float32) (interfaces.Message, error) {
	msg, err := messages.NewEvent(interfaces.TargetTypeObject, targetID)
	if err != nil {
		return nil, errors.Wrap(err, "NewOnCheck")
	}

	o := &OnCheck{
		MessageImpl: msg,
		Values:      values,
	}

	return o, nil
}

// OnCheck Данные датчика обновлены
type OnCheck struct {
	*messages.MessageImpl
	Values map[string]float32 `json:"values"` // Значения
}

func NewOnMotionOn(targetID int) (interfaces.Message, error) {
	msg, err := messages.NewEvent(interfaces.TargetTypeObject, targetID)
	if err != nil {
		return nil, errors.Wrap(err, "NewOnMotionOn")
	}

	o := &OnMotionOn{
		MessageImpl: msg,
	}

	return o, nil
}

// OnMotionOn Движение есть
type OnMotionOn struct {
	*messages.MessageImpl
}

func NewOnMotionOff(targetID int) (interfaces.Message, error) {
	msg, err := messages.NewEvent(interfaces.TargetTypeObject, targetID)
	if err != nil {
		return nil, errors.Wrap(err, "NewOnMotionOff")
	}

	o := &OnMotionOff{
		MessageImpl: msg,
	}

	return o, nil
}

// OnMotionOff Движения нет
type OnMotionOff struct {
	*messages.MessageImpl
}

func NewOnPresenceOn(targetID int) (interfaces.Message, error) {
	msg, err := messages.NewEvent(interfaces.TargetTypeObject, targetID)
	if err != nil {
		return nil, errors.Wrap(err, "NewOnPresenceOn")
	}

	o := &OnPresenceOn{
		MessageImpl: msg,
	}

	return o, nil
}

// OnPresenceOn Присутствие есть
type OnPresenceOn struct {
	*messages.MessageImpl
}

func NewOnPresenceOff(targetID int) (interfaces.Message, error) {
	msg, err := messages.NewEvent(interfaces.TargetTypeObject, targetID)
	if err != nil {
		return nil, errors.Wrap(err, "NewOnPresenceOff")
	}

	o := &OnPresenceOff{
		MessageImpl: msg,
	}

	return o, nil
}

// OnPresenceOff Присутствия нет
type OnPresenceOff struct {
	*messages.MessageImpl
}
