package events

import (
	"github.com/pkg/errors"
	"touchon-server/mqtt/messages"
)

const topicEvent = "object_manager/object/event"

func NewEvent(name string, objectID int, payload map[string]interface{}) (messages.Message, error) {
	m, err := messages.NewMessage(messages.MessageTypeEvent, name, objectID, payload)
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
