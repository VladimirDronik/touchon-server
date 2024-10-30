package messages

import (
	"encoding/json"
	"fmt"
	"time"

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
	TargetTypeScript     TargetType = "script"
	TargetTypeService    TargetType = "service"
)

var TargetTypes = map[TargetType]bool{
	TargetTypeNotMatters: true,
	TargetTypeObject:     true,
	TargetTypeItem:       true,
	TargetTypeScript:     true,
	TargetTypeService:    true,
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

	GetSentAt() time.Time
	SetSentAt(time.Time)
	GetReceivedAt() time.Time
	SetReceivedAt(time.Time)

	json.Marshaler
	json.Unmarshaler
	fmt.Stringer
}

func NewCommand(method string, targetType TargetType, targetID int, methodArgs map[string]interface{}) (Message, error) {
	m, err := NewMessage(MessageTypeCommand, method, targetType, targetID, methodArgs)
	if err != nil {
		return nil, errors.Wrap(err, "NewCommand")
	}

	return m, nil
}

func NewEvent(name string, targetType TargetType, targetID int, payload map[string]interface{}) (Message, error) {
	m, err := NewMessage(MessageTypeEvent, name, targetType, targetID, payload)
	if err != nil {
		return nil, errors.Wrap(err, "NewEvent")
	}

	return m, nil
}

type NotificationType string

const (
	NotificationTypeNormal   = "normal"
	NotificationTypeCritical = "critical"
)

type Notification struct {
	Type NotificationType
	Text string
}

func NewNotification(notType NotificationType, text string) (*Notification, error) {
	switch {
	case notType != NotificationTypeNormal && notType != NotificationTypeCritical:
		return nil, errors.Wrap(errors.Errorf("unknown notification type %q", notType), "NewNotification")
	case text == "":
		return nil, errors.Wrap(errors.New("notification text is empty"), "NewNotification")
	}

	return &Notification{
		Type: notType,
		Text: text,
	}, nil
}
