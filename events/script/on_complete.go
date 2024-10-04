package events

import (
	"github.com/VladimirDronik/touchon-server/event"
	"github.com/VladimirDronik/touchon-server/models"
	"github.com/VladimirDronik/touchon-server/mqtt/messages"
	"github.com/pkg/errors"
)

func init() {
	maker := func() (*event.Event, error) {
		e := &event.Event{
			Code:        "script.on_complete",
			Name:        "Скрипт завершил работу",
			Description: "",
			Props:       event.NewProps(),
			TargetType:  messages.TargetTypeScript,
		}

		r := &event.Prop{
			Code: "result",
			Name: "Результат работы скрипта",
			Item: &models.Item{
				Type: models.DataTypeInterface,
			},
		}

		if err := e.Props.Add(r); err != nil {
			return nil, errors.Wrap(err, "init.maker")
		}

		return e, nil
	}

	if err := event.Register(maker); err != nil {
		panic(err)
	}
}
