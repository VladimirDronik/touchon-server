package interfaces

import "encoding/json"

type MessageType = string

const (
	MessageTypeEvent   MessageType = "event"
	MessageTypeCommand MessageType = "command"
)

type TargetType = string

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

type NotificationType = string

const (
	NotificationTypeNormal   = ""
	NotificationTypeCritical = "critical"
)

type Message interface {
	GetType() MessageType      // event,command
	GetName() string           // onChange,check
	GetTargetType() TargetType // object,item
	GetTargetID() int          // 82

	SetType(MessageType)
	SetName(string)
	SetTargetType(TargetType)
	SetTargetID(int)

	GetPayload() map[string]interface{}
	SetPayload(map[string]interface{})
	GetValue(k string) interface{}
	SetValue(k string, v interface{})
	GetFloatValue(name string) (float32, error)
	GetStringValue(name string) (string, error)
	GetIntValue(name string) (int, error)
	GetBoolValue(name string) (bool, error)

	json.Marshaler
	json.Unmarshaler
}

type MessageSender interface {
	Send(msg ...Message) error
}

type MessagesService interface {
	Start() error
	Shutdown() error
	Subscribe(msgType MessageType, name string, targetType TargetType, targetID *int, handler MsgHandler) (int, error)
	Unsubscribe(handlerIDs ...int)
	MessageSender
}

type MsgHandler func(svc MessageSender, msg Message)

type Command interface {
	Message
	GetArgs() map[string]interface{}
}
