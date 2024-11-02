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
			Name:        "on_change_state",
			Description: "Изменение состояния объекта",
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

	// Для регистрации событий надо в service/init.go добавить импорт данного _пакета_!
	if err := event.Register(maker); err != nil {
		panic(err)
	}
}

func NewOnChangeStateMessage(topic string, targetType messages.TargetType, targetID int, state, value string) (messages.Message, error) {
	e, err := event.MakeEvent("on_change_state", targetType, targetID, map[string]interface{}{"state": state, "value": value})
	if err != nil {
		return nil, errors.Wrap(err, "NewOnChangeStateMessage")
	}

	m, err := e.ToMqttMessage(topic)
	if err != nil {
		return nil, errors.Wrap(err, "NewOnChangeStateMessage")
	}

	return m, nil
}
