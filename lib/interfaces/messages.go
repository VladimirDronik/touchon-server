package interfaces

type MessageType = string

const (
	MessageTypeEvent        MessageType = "event"
	MessageTypeCommand      MessageType = "command"
	MessageTypeNotification MessageType = "notification"
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

type NotificationType string

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
}

type MessagesService interface {
	Start() error
	Shutdown() error
	Subscribe(msgType MessageType, name string, targetType TargetType, targetID *int, handler MsgHandler) (int, error)
	Unsubscribe(handlerIDs ...int)
	Send(msg ...Message) error
}

type MsgHandler func(msg Message)
