package models

import (
	"fmt"
	"strconv"

	"github.com/VladimirDronik/touchon-server/helpers"
	"github.com/pkg/errors"
)

type DataType string

const (
	DataTypeString    DataType = "string"
	DataTypeEnum      DataType = "enum" // enum<string>
	DataTypeBool      DataType = "bool"
	DataTypeInt       DataType = "int"
	DataTypeFloat     DataType = "float"
	DataTypeInterface DataType = "interface"
)

var DataTypeToGoType = map[DataType]string{
	DataTypeString:    "string",
	DataTypeEnum:      "string",
	DataTypeBool:      "bool",
	DataTypeInt:       "int",
	DataTypeFloat:     "float32",
	DataTypeInterface: "interface",
}

type Item struct {
	Type         DataType          `json:"type"`                  //
	Values       map[string]string `json:"values,omitempty"`      // Для DataTypeEnum
	RoundFloat   bool              `json:"round_float,omitempty"` // Для DataTypeFloat. Округлять вещественные числа до десятых долей
	DefaultValue interface{}       `json:"-"`
	value        interface{}       //
}

func (o *Item) GetValue() interface{} {
	return o.value
}

// GetStringValue Метод-хэлпер для получения строкового значения
func (o *Item) GetStringValue() (string, error) {
	if v, ok := o.value.(string); ok {
		return v, nil
	}

	return "", errors.Wrap(errors.Errorf("value is not string (%T)", o.value), "GetStringValue")
}

// GetBoolValue Метод-хэлпер для получения логического значения
func (o *Item) GetBoolValue() (bool, error) {
	if v, ok := o.value.(bool); ok {
		return v, nil
	}

	return false, errors.Wrap(errors.Errorf("value is not bool (%T)", o.value), "GetBoolValue")
}

// GetEnumValue Метод-хэлпер для получения значения-перечисления
func (o *Item) GetEnumValue() (string, error) {
	if v, ok := o.value.(string); ok {
		return v, nil
	}

	return "", errors.Wrap(errors.Errorf("value is not string (%T)", o.value), "GetEnumValue")
}

func (o *Item) GetIntValue() (int, error) {
	switch v := o.value.(type) {
	case float64:
		return int(v), nil
	case int:
		return v, nil
	case int64:
		return int(v), nil
	case string:
		if v, err := strconv.Atoi(v); err == nil {
			return v, nil
		}
	}

	return 0, errors.Wrap(errors.Errorf("value is not int (%T)", o.value), "GetIntValue")
}

func (o *Item) GetFloatValue() (float32, error) {
	switch v := o.value.(type) {
	case float32:
		return v, nil
	case float64:
		return float32(v), nil
	case int:
		return float32(v), nil
	default:
		return 0, errors.Wrap(errors.Errorf("value is not float32, float64 or int (%T)", o.value), "GetFloatValue")
	}
}

func (o *Item) SetValue(value interface{}) error {
	if value == nil {
		return nil
	}

	switch o.Type {
	case DataTypeString:
		s, ok := value.(string)
		if !ok {
			return errors.Wrap(errors.Errorf("value is not string (%T)", value), "Item.SetValue")
		}
		o.value = s

	case DataTypeEnum:
		s, ok := value.(string)
		if !ok {
			s = fmt.Sprintf("%v", value)
		}
		if _, ok := o.Values[s]; !ok {
			return errors.Wrap(errors.Errorf("value %q not found in enum values %v", s, o.Values), "Item.SetValue")
		}
		o.value = s

	case DataTypeBool:
		switch v := value.(type) {
		case string:
			if v == "" {
				return nil
			}
			b, err := strconv.ParseBool(v)
			if err != nil {
				return errors.Wrap(err, "Item.SetValue")
			}
			o.value = b

		case bool:
			o.value = v

		default:
			return errors.Wrap(errors.Errorf("value is not string or bool (%T)", value), "Item.SetValue")
		}

	case DataTypeInt:
		switch v := value.(type) {
		case string:
			if v == "" {
				return nil
			}
			i, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				return errors.Wrap(err, "SetValue")
			}
			o.value = int(i)

		case float32:
			o.value = int(v)

		case float64:
			o.value = int(v)

		case int:
			o.value = v

		default:
			return errors.Wrap(errors.Errorf("value is not string or int (%T)", value), "Item.SetValue")
		}

	case DataTypeFloat:
		round := o.noRound
		if o.RoundFloat {
			round = o.round
		}

		switch v := value.(type) {
		case string:
			if v == "" {
				return nil
			}
			f, err := strconv.ParseFloat(v, 32)
			if err != nil {
				return errors.Wrap(err, "Item.SetValue")
			}
			o.value = round(float32(f))

		case int:
			o.value = round(float32(v))
		case float64:
			o.value = round(float32(v))
		case float32:
			o.value = round(v)

		default:
			return errors.Wrap(errors.Errorf("value is not string, int or float (%T)", value), "Item.SetValue")
		}

	case DataTypeInterface:
		o.value = value

	default:
		return errors.Wrap(errors.Errorf("unexpected prop data type %s", o.Type), "Item.SetValue")
	}

	return nil
}

func (o *Item) RemoveValue() {
	o.value = nil
}

func (o *Item) Check() error {
	if _, ok := DataTypeToGoType[o.Type]; !ok {
		return errors.Errorf("unknown type %q", o.Type)
	}

	switch {
	case o.Type == DataTypeEnum && len(o.Values) == 0:
		return errors.New("values is empty")
	case o.Type != DataTypeEnum && len(o.Values) > 0:
		return errors.Errorf("may be type must be %q?", DataTypeEnum)
	case o.Type != DataTypeFloat && o.RoundFloat:
		return errors.Errorf("may be type must be %q?", DataTypeFloat)
	}

	// Проверяем DefaultValue
	if o.DefaultValue != nil {
		currValue := o.GetValue()
		if err := o.SetValue(o.DefaultValue); err != nil {
			return err
		}

		o.SetValueUnsafe(currValue)
	}

	return nil
}

func (o *Item) SetValueUnsafe(value interface{}) {
	o.value = value
}

// StringValue возвращает строковое представление значения свойства
func (o *Item) StringValue() string {
	v := fmt.Sprintf("%v", o.value)
	if v == "<nil>" {
		return ""
	}
	return v
}

func (o *Item) round(v float32) float32 {
	return helpers.Round(v)
}

func (o *Item) noRound(v float32) float32 {
	return v
}
