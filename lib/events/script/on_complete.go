package script

import (
	"github.com/pkg/errors"
	"touchon-server/lib/event"
	"touchon-server/lib/models"
	"touchon-server/lib/mqtt/messages"
)

func init() {
	maker := func() (*event.Event, error) {
		e := &event.Event{
			Code:        "script.on_complete",
			Name:        "on_complete",
			Description: "Скрипт завершил работу",
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

	// Для регистрации событий надо в service/init.go добавить импорт данного _пакета_!
	if err := event.Register(maker); err != nil {
		panic(err)
	}
}

func NewOnCompleteMessage(topic string, targetID int, result interface{}) (messages.Message, error) {
	e, err := event.MakeEvent("script.on_complete", messages.TargetTypeObject, targetID, map[string]interface{}{"result": result})
	if err != nil {
		return nil, errors.Wrap(err, "NewOnCompleteMessage")
	}

	m, err := e.ToMqttMessage(topic)
	if err != nil {
		return nil, errors.Wrap(err, "NewOnCompleteMessage")
	}

	return m, nil
}
