package script

import (
	"github.com/pkg/errors"
	"touchon-server/lib/event"
	"touchon-server/lib/interfaces"
	"touchon-server/lib/messages"
)

func NewOnComplete(scriptID int, result interface{}) (interfaces.Event, error) {
	msg, err := messages.NewEvent("script.on_complete", interfaces.TargetTypeScript, scriptID)
	if err != nil {
		return nil, errors.Wrap(err, "NewOnComplete")
	}

	o := &OnComplete{
		Event: &event.EventImpl{
			Message:     msg,
			Title:       "on_complete",
			Description: "Скрипт завершил работу",
		},
	}

	o.SetValue("result", result) // Результат работы скрипта

	return o, nil
}

type OnComplete struct {
	interfaces.Event
}
