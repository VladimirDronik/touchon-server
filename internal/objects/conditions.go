package objects

import (
	"log"

	"github.com/pkg/errors"
)

type Condition interface {
	Check(props *Props) bool
}

// ----------------------------------------------------------------------------------

func True() *trueCondition {
	return &trueCondition{}
}

type trueCondition struct{}

func (o *trueCondition) Check(props *Props) bool { return true }

// ----------------------------------------------------------------------------------

func False() *falseCondition {
	return &falseCondition{}
}

type falseCondition struct{}

func (o *falseCondition) Check(props *Props) bool { return false }

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

func NewCondition() *baseCondition {
	return &baseCondition{}
}

type baseCondition struct {
	items []*ConditionItem
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

func (o *baseCondition) Check(props *Props) bool {
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
