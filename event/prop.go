package event

import (
	"encoding/json"

	"github.com/VladimirDronik/touchon-server/models"
	"github.com/pkg/errors"
)

type Prop struct {
	Code        string `json:"code"` // unique
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	*models.Item
}

func (o *Prop) MarshalJSON() ([]byte, error) {
	type jsonProp struct {
		Code        string `json:"code"` // unique
		Name        string `json:"name"`
		Description string `json:"description,omitempty"`
		*models.Item
		Value interface{} `json:"value,omitempty"`
	}

	v := &jsonProp{
		Code:        o.Code,
		Name:        o.Name,
		Description: o.Description,
		Item:        o.Item,
		Value:       o.GetValue(),
	}

	return json.Marshal(v)
}

func (o *Prop) UnmarshalJSON(data []byte) error {
	type jsonProp struct {
		Code        string `json:"code"` // unique
		Name        string `json:"name"`
		Description string `json:"description,omitempty"`
		*models.Item
		Value interface{} `json:"value,omitempty"`
	}

	v := &jsonProp{}
	if err := json.Unmarshal(data, v); err != nil {
		return err
	}

	o.Code = v.Code
	o.Name = v.Name
	o.Description = v.Description
	o.Item = v.Item
	if err := o.SetValue(v.Value); err != nil {
		return errors.Wrapf(err, "Prop.UnmarshalJSON(%s)", o.Code)
	}

	return nil
}

func (o *Prop) Check() error {
	switch {
	case o.Code == "":
		return errors.New("prop Code is empty")
	case o.Name == "":
		return errors.New("prop name is empty")
	}

	if err := o.Item.Check(); err != nil {
		return err
	}

	return nil
}
