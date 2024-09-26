package events

import (
	"github.com/pkg/errors"
	"touchon-server/mqtt/messages"
)

const topicError = "object_manager/error"

func NewOnError(topic string, objectID int, err error) (messages.Message, error) {
	if err == nil {
		return nil, errors.Wrap(errors.New("err is nil"), "NewOnError")
	}

	payload := map[string]interface{}{"error": err.Error()}

	impl, err := NewEvent("onError", objectID, payload)
	if err != nil {
		return nil, errors.Wrap(err, "NewOnError")
	}

	t := topicError
	if topic != "" {
		t += "/" + topic
	}

	impl.SetTopic(t)

	return &OnError{Message: impl}, nil
}

type OnError struct {
	messages.Message
}

func (o *OnError) GetErrorMessage() (string, error) {
	v, err := o.GetStringValue("error")
	if err != nil {
		return "", errors.Wrap(err, "GetErrorMessage")
	}

	return v, nil
}

func (o *OnError) SetErrorMessage(v string) {
	o.GetPayload()["error"] = v
}

//
//func NewOnError(topic string, err error) messages.Message {
//	if err == nil {
//		return nil
//	}
//
//	if topic != "" {
//		topic = "/" + topic
//	}
//
//	msg := messages.NewEvent()
//
//	return messages.Message{
//		Topic:   topicError + topic,
//		Content: err.Error(),
//	}
//}
