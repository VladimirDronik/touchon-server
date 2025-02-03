package scripts

import (
	"github.com/pkg/errors"
	"touchon-server/lib/models"
)

type Param struct {
	Code        string `json:"code"` // unique
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	*models.Item
}

func (o *Param) Check() error {
	switch {
	case o.Code == "":
		return errors.New("param Code is empty")
	case o.Name == "":
		return errors.New("param name is empty")
	}

	if err := o.Item.Check(); err != nil {
		return err
	}

	return nil
}
