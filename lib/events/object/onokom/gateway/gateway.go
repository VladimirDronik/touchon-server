package gateway

import (
	"github.com/pkg/errors"
	"touchon-server/lib/event"
	"touchon-server/lib/interfaces"
	"touchon-server/lib/messages"
)

func NewOnCheck(targetID int, props map[string]interface{}) (interfaces.Event, error) {
	msg, err := messages.NewEvent("object.onokom.gateway.on_check", interfaces.TargetTypeObject, targetID)
	if err != nil {
		return nil, errors.Wrap(err, "NewOnCheck")
	}

	o := &OnCheck{
		Event: &event.EventImpl{
			Message:     msg,
			Title:       "on_check",
			Description: "Получение состояния устройства",
		},
	}

	for k, v := range props {
		o.SetValue(k, v)
	}

	return o, nil
}

// OnCheck Получение состояния устройства
type OnCheck struct {
	interfaces.Event
}

func NewOnChange(targetID int, props map[string]interface{}) (interfaces.Event, error) {
	msg, err := messages.NewEvent("object.onokom.gateway.on_change", interfaces.TargetTypeObject, targetID)
	if err != nil {
		return nil, errors.Wrap(err, "NewOnChange")
	}

	o := &OnChange{
		Event: &event.EventImpl{
			Message:     msg,
			Title:       "on_change",
			Description: "Изменение состояния устройства",
		},
	}

	for k, v := range props {
		o.SetValue(k, v)
	}

	return o, nil
}

// OnChange Изменение состояния устройства
type OnChange struct {
	interfaces.Event
}
