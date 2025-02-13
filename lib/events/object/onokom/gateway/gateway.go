package gateway

import (
	"github.com/pkg/errors"
	"touchon-server/lib/event"
	"touchon-server/lib/interfaces"
	"touchon-server/lib/messages"
)

func NewOnCheck(targetID int, props map[string]interface{}) (interfaces.Event, error) {
	msg, err := messages.NewEvent(interfaces.TargetTypeObject, targetID)
	if err != nil {
		return nil, errors.Wrap(err, "NewOnCheck")
	}

	o := &OnCheck{
		Event: &event.EventImpl{
			Message:          msg,
			EventCode:        "object.onokom.gateway.on_check",
			EventName:        "on_check",
			EventDescription: "Получение состояния устройства",
		},
		Props: props,
	}

	return o, nil
}

// OnCheck Получение состояния устройства
type OnCheck struct {
	interfaces.Event
	Props map[string]interface{} `json:"props,omitempty"` // Свойства кондиционера
}

func NewOnChange(targetID int, props map[string]interface{}) (interfaces.Event, error) {
	msg, err := messages.NewEvent(interfaces.TargetTypeObject, targetID)
	if err != nil {
		return nil, errors.Wrap(err, "NewOnChange")
	}

	o := &OnChange{
		Event: &event.EventImpl{
			Message:          msg,
			EventCode:        "object.onokom.gateway.on_change",
			EventName:        "on_change",
			EventDescription: "Изменение состояния устройства",
		},
		Props: props,
	}

	return o, nil
}

// OnChange Изменение состояния устройства
type OnChange struct {
	interfaces.Event
	Props map[string]interface{} `json:"props,omitempty"` // Свойства кондиционера
}
