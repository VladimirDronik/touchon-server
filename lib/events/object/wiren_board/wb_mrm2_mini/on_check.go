package wb_mrm2_mini

import (
	"github.com/pkg/errors"
	"touchon-server/lib/event"
	"touchon-server/lib/interfaces"
	"touchon-server/lib/messages"
)

func NewOnCheck(targetID int, values map[string]bool) (interfaces.Event, error) {
	msg, err := messages.NewEvent(interfaces.TargetTypeObject, targetID)
	if err != nil {
		return nil, errors.Wrap(err, "NewOnCheck")
	}

	o := &OnCheck{
		Event: &event.EventImpl{
			Message:          msg,
			EventCode:        "object.wiren_board.wb_mrm2_mini.on_check",
			EventName:        "on_check",
			EventDescription: "Получено состояние выходов",
		},
	}

	if v, ok := values["k1"]; ok {
		o.K1 = v
	}

	if v, ok := values["k2"]; ok {
		o.K2 = v
	}

	return o, nil
}

// OnCheck Получено состояние выходов
type OnCheck struct {
	interfaces.Event
	K1 bool `json:"k1,omitempty"` // Выход 1
	K2 bool `json:"k2,omitempty"` // Выход 2
}
