package wb_mrm2_mini

import (
	"github.com/pkg/errors"
	"touchon-server/lib/interfaces"
	"touchon-server/lib/messages"
)

func NewOnCheck(targetID int, values map[string]bool) (interfaces.Message, error) {
	msg, err := messages.NewEvent(interfaces.TargetTypeObject, targetID)
	if err != nil {
		return nil, errors.Wrap(err, "NewOnCheck")
	}

	o := &OnCheck{
		MessageImpl: msg,
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
	*messages.MessageImpl
	K1 bool `json:"k1"` // Выход 1
	K2 bool `json:"k2"` // Выход 2
}
