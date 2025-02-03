package objects

import (
	"encoding/json"

	"github.com/pkg/errors"
	"touchon-server/lib/models"
)

// PropValueChecker Свойства объекта должны уметь проверять свои значения
type PropValueChecker func(prop *Prop, allProps map[string]*Prop) error

type Prop struct {
	Code        string `json:"code"` // unique
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	*models.Item
	Required   Condition        `json:"required"`
	Editable   Condition        `json:"editable"` // Можно ли изменять значение или оно неизменяемо (например, показания температуры)
	Visible    Condition        `json:"visible"`
	CheckValue PropValueChecker `json:"-"` // Поле, которое при реализации должно содержать метод проверки значения
}

func (o *Prop) MarshalJSON() ([]byte, error) {
	type jsonProp struct {
		Code        string `json:"code"` // unique
		Name        string `json:"name"`
		Description string `json:"description,omitempty"`
		*models.Item
		Required Condition   `json:"required"`
		Editable Condition   `json:"editable"` // Можно ли изменять значение или оно неизменяемо (например, показания температуры)
		Visible  Condition   `json:"visible"`
		Value    interface{} `json:"value,omitempty"`
	}

	value := o.GetValue()
	if value == nil {
		value = o.DefaultValue
	}

	v := &jsonProp{
		Code:        o.Code,
		Name:        o.Name,
		Description: o.Description,
		Item:        o.Item,
		Required:    o.Required,
		Editable:    o.Editable,
		Visible:     o.Visible,
		Value:       value,
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
		return errors.Wrapf(err, "Prop.UnmarshalJSON(%s)", o.Code)
	}

	o.Code = v.Code
	o.Name = v.Name
	o.Description = v.Description
	o.Item = &models.Item{}
	o.SetValueUnsafe(v.Value)

	return nil
}

func (o *Prop) CheckDefinition() error {
	switch {
	case o.Code == "":
		return errors.New("prop Code is empty")
	case o.Name == "":
		return errors.New("prop Name is empty")
	case o.Required == nil:
		return errors.New("prop Required is nil")
	case o.Visible == nil:
		return errors.New("prop Visible is nil")
	case o.Editable == nil:
		return errors.New("prop Editable is nil")
	}

	if err := o.Item.Check(); err != nil {
		return errors.Wrapf(err, "Prop.CheckDefinition(%s)", o.Code)
	}

	return nil
}
