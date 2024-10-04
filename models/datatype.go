package models

import (
	"strconv"

	"github.com/VladimirDronik/touchon-server/helpers"
	"github.com/pkg/errors"
)

type DataType string

const (
	DataTypeString DataType = "string"
	DataTypeEnum   DataType = "enum" // enum<string>
	DataTypeBool   DataType = "bool"
	DataTypeInt    DataType = "int"
	DataTypeFloat  DataType = "float"
)

var DataTypeToGoType = map[DataType]string{
	DataTypeString: "string",
	DataTypeEnum:   "string",
	DataTypeBool:   "bool",
	DataTypeInt:    "int",
	DataTypeFloat:  "float32",
}

type Item struct {
	Type  DataType
	value interface{}
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
	switch o.Type {
	case DataTypeString:
		s, ok := value.(string)
		if !ok {
			return errors.Wrap(errors.Errorf("value is not string (%T)", value), "SetValue")
		}
		o.value = s

	case DataTypeEnum:
		s, ok := value.(string)
		if !ok {
			return errors.Wrap(errors.Errorf("value is not string (%T)", value), "SetValue")
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
				return errors.Wrap(err, "SetValue")
			}
			o.value = b

		case bool:
			o.value = v

		default:
			return errors.Wrap(errors.Errorf("value is not string or bool (%T)", value), "SetValue")
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
			o.value = i

		case int:
			o.value = v

		default:
			return errors.Wrap(errors.Errorf("value is not string or int (%T)", value), "SetValue")
		}

	case DataTypeFloat:
		switch v := value.(type) {
		case string:
			if v == "" {
				return nil
			}
			f, err := strconv.ParseFloat(v, 32)
			if err != nil {
				return errors.Wrap(err, "SetValue")
			}
			o.value = helpers.Round(float32(f))

		case int:
			o.value = helpers.Round(float32(v))
		case float64:
			o.value = helpers.Round(float32(v))
		case float32:
			o.value = helpers.Round(v)

		default:
			return errors.Wrap(errors.Errorf("value is not string, int or float (%T)", value), "SetValue")
		}

	default:
		return errors.Wrap(errors.Errorf("unexpected prop data type %s", o.Type), "SetValue")
	}

	return nil
}
