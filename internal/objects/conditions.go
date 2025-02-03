package objects

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/pkg/errors"
	"touchon-server/internal/model"
)

type Condition interface {
	Check(userAccessLevel model.AccessLevel, props *Props) bool
	GetJS() (string, error)
	json.Marshaler
	json.Unmarshaler
}

// ----------------------------------------------------------------------------------

func NewRequired(required bool) *requiredCondition {
	return &requiredCondition{required: required}
}

type requiredCondition struct {
	required bool
}

func (o *requiredCondition) Check(userAccessLevel model.AccessLevel, props *Props) bool {
	return o.required
}

func (o *requiredCondition) GetJS() (string, error) {
	//return fmt.Sprintf("function(userAccessLevel, props) { return %t; }", o.required), nil
	return fmt.Sprintf("return %t;", o.required), nil
}

func (o *requiredCondition) MarshalJSON() ([]byte, error) {
	js, err := o.GetJS()
	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(js)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// Из JS обратно не создаем логику условия
func (o *requiredCondition) UnmarshalJSON([]byte) error {
	return nil
}

// ----------------------------------------------------------------------------------

type Operator string

const (
	OpEq    Operator = "=="
	OpNotEq Operator = "!="
	OpIn    Operator = "in"
)

type ConditionItem struct {
	PropName string
	Operator Operator
	Value    interface{}
}

func (o *ConditionItem) GetJS() (string, error) {
	switch o.Operator {
	case OpEq, OpNotEq:
		if s, ok := o.Value.(string); ok {
			o.Value = `"` + s + `"`
		}

		return fmt.Sprintf(`if (!(props[%q] %s %v)) { return false; }`, o.PropName, o.Operator, o.Value), nil

	case OpIn:
		values, ok := o.Value.([]interface{})
		if !ok {
			return "", errors.Wrap(errors.New("value is not array"), "ConditionItem.Check")
		}

		s := make([]string, 0, len(values))
		for _, value := range values {
			if s, ok := value.(string); ok {
				value = `"` + s + `"`
			}

			s = append(s, fmt.Sprintf("props[%q] == %v", o.PropName, value))
		}

		return fmt.Sprintf(`if (!(%s)) { return false; }`, strings.Join(s, " || ")), nil

	default:
		return "", errors.Wrap(errors.Errorf("unknown operator %q", o.Operator), "ConditionItem.GetJS")
	}
}

func NewCondition() *baseCondition {
	return &baseCondition{userAccessLevel: model.AccessLevelAllowed}
}

type baseCondition struct {
	userAccessLevel model.AccessLevel
	items           []*ConditionItem
}

func (o *baseCondition) AccessLevel(level model.AccessLevel) *baseCondition {
	o.userAccessLevel = level
	return o
}

func (o *baseCondition) PropValueEq(propName string, value interface{}) *baseCondition {
	o.items = append(o.items, &ConditionItem{
		PropName: propName,
		Operator: OpEq,
		Value:    value,
	})
	return o
}

func (o *baseCondition) PropValueNotEq(propName string, value interface{}) *baseCondition {
	o.items = append(o.items, &ConditionItem{
		PropName: propName,
		Operator: OpNotEq,
		Value:    value,
	})
	return o
}

func (o *baseCondition) PropValueIn(propName string, values ...interface{}) *baseCondition {
	o.items = append(o.items, &ConditionItem{
		PropName: propName,
		Operator: OpIn,
		Value:    values,
	})
	return o
}

func (o *baseCondition) GetJS() (string, error) {
	//s := fmt.Sprintf("function(userAccessLevel, props) { if (userAccessLevel < %d) { return false; }", o.userAccessLevel)
	s := fmt.Sprintf("if (userAccessLevel < %d) { return false; }", o.userAccessLevel)

	for _, item := range o.items {
		js, err := item.GetJS()
		if err != nil {
			return "", errors.Wrap(err, "baseCondition.GetJS")
		}

		s += " " + js
	}

	//return s + " return true; }", nil
	return s + " return true;", nil
}

func (o *baseCondition) MarshalJSON() ([]byte, error) {
	js, err := o.GetJS()
	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(js)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// Из JS обратно не создаем логику условия
func (o *baseCondition) UnmarshalJSON([]byte) error {
	return nil
}

func (o *baseCondition) Check(userAccessLevel model.AccessLevel, props *Props) bool {
	if userAccessLevel < o.userAccessLevel {
		return false
	}

	for _, item := range o.items {
		prop, err := props.Get(item.PropName)
		if err != nil {
			return false
		}

		switch item.Operator {
		case OpEq:
			if prop.GetValue() != item.Value {
				return false
			}

		case OpNotEq:
			if prop.GetValue() == item.Value {
				return false
			}

		case OpIn:
			values, ok := item.Value.([]interface{})
			if !ok {
				log.Println(errors.Wrap(errors.New("item value is not array"), "baseCondition.Check"))
				return false
			}

			r := false
			for _, value := range values {
				if prop.GetValue() == value {
					r = true
					break
				}
			}
			if !r {
				return false
			}

		default:
			log.Println(errors.Wrap(errors.Errorf("unknown operator %q", item.Operator), "baseCondition.Check"))
			return false
		}
	}

	return true
}

// required (bool || deps)
// editable (ACL && deps)
// visible (ACL && deps)
// checks: js + go
