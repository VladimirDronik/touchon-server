package event

import (
	"github.com/pkg/errors"
)

type Event struct {
	Code        string `json:"code"` // unique
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Props       *Props `json:"props"`
}

func (o *Event) Check() error {
	switch {
	case o.Code == "":
		return errors.New("code is empty")
	case o.Name == "":
		return errors.New("name is empty")
	}

	return nil
}
