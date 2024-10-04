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
			Code:        "on_change_state",
			Name:        "Ошибка",
			Description: "",
			Props:       event.NewProps(),
			TargetType:  messages.TargetTypeNotMatters,
		}

		s := &event.Prop{
			Code: "state",
			Name: "Состояние",
			Item: &models.Item{
				Type: models.DataTypeString,
			},
		}

		v := &event.Prop{
			Code: "value",
			Name: "Значение",
			Item: &models.Item{
				Type: models.DataTypeString,
			},
		}

		if err := e.Props.Add(s, v); err != nil {
			return nil, errors.Wrap(err, "init.maker")
		}

		return e, nil
	}

	if err := event.Register(maker); err != nil {
		panic(err)
	}
}
