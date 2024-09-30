package events

import (
	"github.com/VladimirDronik/touchon-server/mqtt/messages"
	"github.com/pkg/errors"
)

const topicEvent = "object_manager/object/event"

func NewEvent(name string, targetID int, targetType messages.TargetType, payload map[string]interface{}) (messages.Message, error) {
	m, err := messages.NewMessage(messages.MessageTypeEvent, name, targetID, targetType, payload)
	if err != nil {
		return nil, errors.Wrap(err, "NewEvent")
	}

	m.SetTopic(topicEvent)

	return m, nil
}

//func MakeObjectEvent(eventPayload interface{}) messages.Message {
//	payload, err := json.Marshal(eventPayload)
//	if err != nil {
//		return messages.Message{
//			Topic:   topicError,
//			Content: "Internal Error: 94",
//		}
//	}
//
//	return messages.Message{
//		Topic:   topicEvent,
//		Content: string(payload),
//	}
//}
