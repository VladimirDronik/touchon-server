package event

import (
	"github.com/pkg/errors"
	"touchon-server/lib/interfaces"
)

type Event struct {
	Code        string `json:"code"` // unique
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Props       *props `json:"props"`

	TargetID   int                   `json:"target_id,omitempty"`
	TargetType interfaces.TargetType `json:"target_type,omitempty"`
}

func (o *Event) Check() error {
	if _, ok := interfaces.TargetTypes[o.TargetType]; !ok {
		return errors.Wrap(errors.Errorf("unknown target type %q", o.TargetType), "Event.Check")
	}

	switch {
	case o.Code == "":
		return errors.Wrap(errors.New("code is empty"), "Event.Check")
	case o.Name == "":
		return errors.Wrap(errors.New("name is empty"), "Event.Check")
	case o.Props == nil:
		return errors.Wrap(errors.New("props is empty"), "Event.Check")
	}

	if err := o.Props.Check(); err != nil {
		return errors.Wrap(err, "Event.Check")
	}

	return nil
}
