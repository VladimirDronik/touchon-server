package event

import (
	"github.com/pkg/errors"
	"touchon-server/lib/interfaces"
)

type EventImpl struct {
	interfaces.Message
	EventCode        string `json:"code"` // unique
	EventName        string `json:"name"`
	EventDescription string `json:"description,omitempty"`
}

func (o *EventImpl) GetEventCode() string {
	return o.EventCode
}

func (o *EventImpl) GetEventName() string {
	return o.EventName
}

func (o *EventImpl) GetEventDescription() string {
	return o.EventDescription
}

func (o *EventImpl) CheckEvent() error {
	switch {
	case o.EventCode == "":
		return errors.Wrap(errors.New("EventCode is empty"), "EventImpl.CheckEvent")
	case o.EventName == "":
		return errors.Wrap(errors.New("EventName is empty"), "EventImpl.CheckEvent")
	}

	return nil
}
