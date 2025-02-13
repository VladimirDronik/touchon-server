package event

import (
	"encoding/json"

	"github.com/pkg/errors"
	"touchon-server/lib/interfaces"
)

type EventImpl struct {
	interfaces.Message
	Title       string
	Description string
}

func (o *EventImpl) GetTitle() string {
	return o.Title
}

func (o *EventImpl) GetDescription() string {
	return o.Description
}

func (o *EventImpl) CheckEvent() error {
	switch {
	case o.GetName() == "": // code
		return errors.Wrap(errors.New("Name is empty"), "EventImpl.CheckEvent")
	case o.Title == "": // name
		return errors.Wrap(errors.New("Title is empty"), "EventImpl.CheckEvent")
	}

	return nil
}

func (o *EventImpl) MarshalJSON() ([]byte, error) {
	type R struct {
		Code        string `json:"code"`
		Name        string `json:"name"`
		Description string `json:"description,omitempty"`
	}

	return json.Marshal(&R{o.GetName(), o.Title, o.Description})
}
