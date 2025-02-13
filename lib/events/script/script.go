package script

import (
	"github.com/pkg/errors"
	"touchon-server/lib/event"
	"touchon-server/lib/interfaces"
	"touchon-server/lib/messages"
)

func NewOnComplete(scriptID int, result interface{}) (interfaces.Event, error) {
	msg, err := messages.NewEvent(interfaces.TargetTypeScript, scriptID)
	if err != nil {
		return nil, errors.Wrap(err, "NewOnComplete")
	}

	o := &OnComplete{
		Event: &event.EventImpl{
			Message:          msg,
			EventCode:        "script.on_complete",
			EventName:        "on_complete",
			EventDescription: "Скрипт завершил работу",
		},
		Result: result,
	}

	return o, nil
}

type OnComplete struct {
	interfaces.Event
	Result interface{} `json:"result,omitempty"` // Результат работы скрипта
}
