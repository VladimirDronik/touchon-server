package interfaces

type Event interface {
	Message
	GetEventCode() string
	GetEventName() string
	GetEventDescription() string
	CheckEvent() error
}
