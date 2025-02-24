package wb_mrm2_mini

import (
	"github.com/pkg/errors"
	"touchon-server/lib/event"
	"touchon-server/lib/interfaces"
	"touchon-server/lib/messages"
)

func NewOnCheck(targetID int, values map[string]bool) (interfaces.Event, error) {
	msg, err := messages.NewEvent("object.wiren_board.wb_mrm2_mini.on_check", interfaces.TargetTypeObject, targetID)
	if err != nil {
		return nil, errors.Wrap(err, "NewOnCheck")
	}

	o := &OnCheck{
		Event: &event.EventImpl{
			Message:     msg,
			Title:       "on_check",
			Description: "Получено состояние выходов",
		},
	}

	for _, k := range []string{"k1", "k2"} {
		if v, ok := values[k]; ok {
			o.SetValue(k, v)
		}
	}

	return o, nil
}

// OnCheck Получено состояние выходов
type OnCheck struct {
	interfaces.Event
}
