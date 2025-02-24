package interfaces

import "encoding/json"

type Event interface {
	Message
	GetTitle() string
	GetDescription() string
	CheckEvent() error

	json.Marshaler
}
