package events

import (
	"github.com/VladimirDronik/touchon-server/mqtt/messages"
	"github.com/pkg/errors"
)

const objectEventTopic = "object_manager/object/event"

func NewEvent(name string, targetID int, targetType messages.TargetType, payload map[string]interface{}) (messages.Message, error) {
	m, err := messages.NewMessage(messages.MessageTypeEvent, name, targetID, targetType, payload)
	if err != nil {
		return nil, errors.Wrap(err, "NewEvent")
	}

	m.SetTopic(objectEventTopic)

	return m, nil
}
