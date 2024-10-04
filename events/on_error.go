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
			Code:        "on_error",
			Name:        "Ошибка",
			Description: "",
			Props:       event.NewProps(),
			TargetType:  messages.TargetTypeNotMatters,
		}

		msg := &event.Prop{
			Code: "msg",
			Name: "Текст ошибки",
			Item: &models.Item{
				Type: models.DataTypeString,
			},
		}

		if err := e.Props.Add(msg); err != nil {
			return nil, errors.Wrap(err, "init.maker")
		}

		return e, nil
	}

	if err := event.Register(maker); err != nil {
		panic(err)
	}
}
