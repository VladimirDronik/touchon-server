package messages

import (
	"encoding/json"

	"github.com/pkg/errors"
)

var Publisher string

type MessageType string

const (
	MessageTypeEvent   MessageType = "event"
	MessageTypeCommand MessageType = "command"
)

type QoS int

const (
	QoSNotGuaranteed QoS = iota // Нет никаких гарантий
	QoSMinimumOne               // Гарантировано минимум один раз
	QoSGuaranteedOne            // Гарантировано один раз
)

type TargetType string

const (
	TargetTypeNotMatters TargetType = "not_matters" // Используется, например, в определении ошибки, относящейся к любому типу
	TargetTypeObject     TargetType = "object"
	TargetTypeItem       TargetType = "item"
)

var TargetTypes = map[TargetType]bool{
	TargetTypeNotMatters: true,
	TargetTypeObject:     true,
	TargetTypeItem:       true,
}

type Message interface {
	GetRetained() bool
	GetPublisher() string
	GetTopic() string                   // action_router/object/method
	GetType() MessageType               // event,command
	GetName() string                    // onChange,check
	GetTargetID() int                   // 82
	GetTargetType() TargetType          // object,item
	GetPayload() map[string]interface{} //
	GetQoS() QoS                        //
	GetTopicPublisher() string          // action_router
	GetTopicType() string               // object
	GetTopicAction() string             // method

	SetRetained(bool)
	SetPublisher(string)
	SetTopic(string)
	SetType(MessageType)
	SetName(string)
	SetTargetID(int)
	SetTargetType(TargetType)
	SetPayload(map[string]interface{})
	SetQoS(QoS)

	GetFloatValue(name string) (float32, error)
	GetStringValue(name string) (string, error)
	GetIntValue(name string) (int, error)
	GetBoolValue(name string) (bool, error)

	json.Marshaler
	json.Unmarshaler
}

func NewCommand(method string, targetID int, targetType TargetType, methodArgs map[string]interface{}) (Message, error) {
	m, err := NewMessage(MessageTypeCommand, method, targetID, targetType, methodArgs)
	if err != nil {
		return nil, errors.Wrap(err, "NewCommand")
	}

	return m, nil
}

func NewEvent(name string, targetID int, targetType TargetType, payload map[string]interface{}) (Message, error) {
	m, err := NewMessage(MessageTypeEvent, name, targetID, targetType, payload)
	if err != nil {
		return nil, errors.Wrap(err, "NewEvent")
	}

	return m, nil
}
