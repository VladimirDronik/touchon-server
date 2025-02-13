package interfaces

import "encoding/json"

type Event interface {
	Message
	GetEventCode() string
	GetEventName() string
	GetEventDescription() string
	CheckEvent() error

	json.Marshaler
}
