package events

import (
	"github.com/VladimirDronik/touchon-server/event"
	"github.com/VladimirDronik/touchon-server/models"
	"github.com/VladimirDronik/touchon-server/mqtt/messages"
	"github.com/pkg/errors"
)

type NotifyType string

const (
	NotifyTypeNonCritical NotifyType = ""
	NotifyTypeCritical    NotifyType = "critical"
)

func init() {
	maker := func() (*event.Event, error) {
		e := &event.Event{
			Code:        "on_notify",
			Name:        "on_notify",
			Description: "Уведомление",
			Props:       event.NewProps(),
			TargetType:  messages.TargetTypeNotMatters,
		}

		m := &event.Prop{
			Code: "msg",
			Name: "Текст",
			Item: &models.Item{
				Type: models.DataTypeString,
			},
		}

		t := &event.Prop{
			Code: "type",
			Name: "Тип",
			Item: &models.Item{
				Type: models.DataTypeString,
			},
		}

		if err := e.Props.Add(m, t); err != nil {
			return nil, errors.Wrap(err, "init.maker")
		}

		return e, nil
	}

	// Для регистрации событий надо в service/init.go добавить импорт данного _пакета_!
	if err := event.Register(maker); err != nil {
		panic(err)
	}
}

func NewOnNotifyMessage(topic string, msg string, notifyType NotifyType) (messages.Message, error) {
	e, err := event.MakeEvent("on_notify", messages.TargetTypeNotMatters, 0, map[string]interface{}{"msg": msg, "type": string(notifyType)})
	if err != nil {
		return nil, errors.Wrap(err, "NewOnNotifyMessage")
	}

	m, err := e.ToMqttMessage(topic)
	if err != nil {
		return nil, errors.Wrap(err, "NewOnNotifyMessage")
	}

	return m, nil
}
