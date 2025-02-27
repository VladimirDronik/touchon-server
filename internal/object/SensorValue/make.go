package SensorValue

import (
	"github.com/pkg/errors"
	"touchon-server/internal/objects"
)

func init() {
	// Регистрируем все значения датчиков
	for t := range Names {
		f := func(t Type) objects.ObjectMaker {
			return func(withChildren bool) (objects.Object, error) {
				return Make(t, withChildren)
			}
		}
		_ = objects.Register(f(t))
	}
}

type Type = string

const (
	TypeTemperature  = "temperature"  // Температура
	TypeHumidity     = "humidity"     // Влажность
	TypeCO2          = "co2"          // Уровень СО2
	TypeIllumination = "illumination" // Освещенность
	TypePressure     = "pressure"     // Давление
	TypeCurrent      = "current"      // Ток
	TypeVoltage      = "voltage"      // Напряжение
	TypeMotion       = "motion"       // Движение
	TypePresence     = "presence"     // Присутствие
)

var Names = map[Type]string{
	TypeTemperature:  "Температура",
	TypeHumidity:     "Влажность",
	TypeCO2:          "СО2",
	TypeIllumination: "Освещенность",
	TypePressure:     "Давление",
	TypeCurrent:      "Ток",
	TypeVoltage:      "Напряжение",
	TypeMotion:       "Движение",
	TypePresence:     "Присутствие",
}

func Make(valueType Type, withChildren bool) (objects.Object, error) {
	name, ok := Names[valueType]
	if !ok {
		return nil, errors.Wrap(errors.Errorf("unknown value type %q", valueType), "SensorValue.Make")
	}

	obj, err := MakeModel(withChildren)
	if err != nil {
		return nil, errors.Wrap(err, "SensorValue.Make")
	}

	obj.SetType(string(valueType))
	obj.SetName(name)
	obj.SetTags(string(valueType))
	var props map[string]interface{}

	switch valueType {
	case TypeTemperature:
		props = map[string]interface{}{
			"max_error_value": 50,
			"max_threshold":   30,
			"min_error_value": -50,
			"min_threshold":   0,
			"unit":            "℃",
		}

	case TypeHumidity:
		props = map[string]interface{}{
			"max_error_value": 100,
			"max_threshold":   80,
			"min_error_value": 0,
			"min_threshold":   0,
			"unit":            "%",
		}

	case TypeCO2:
		props = map[string]interface{}{
			"max_error_value": 3000,
			"max_threshold":   1800,
			"min_error_value": 300,
			"min_threshold":   350,
			"unit":            "ppm",
		}

	case TypeIllumination:
		props = map[string]interface{}{
			"max_error_value": 5000,
			"max_threshold":   2700,
			"min_error_value": -50,
			"min_threshold":   0,
			"unit":            "lux",
		}

	case TypePressure:
		props = map[string]interface{}{
			"max_error_value": 5000,
			"max_threshold":   2700,
			"min_error_value": -50,
			"min_threshold":   0,
			"unit":            "pa",
		}

	case TypeCurrent:
		props = map[string]interface{}{
			"max_error_value": 5000,
			"max_threshold":   2700,
			"min_error_value": -50,
			"min_threshold":   0,
			"unit":            "u",
		}

	case TypeVoltage:
		props = map[string]interface{}{
			"max_error_value": 1000,
			"max_threshold":   550,
			"min_error_value": -50,
			"min_threshold":   0,
			"unit":            "v",
		}

	case TypeMotion:
		props = map[string]interface{}{
			"max_error_value": 2,
			"max_threshold":   2,
			"min_error_value": -1,
			"min_threshold":   -1,
			"unit":            "",
		}

	case TypePresence:
		props = map[string]interface{}{
			"max_error_value": 2,
			"max_threshold":   2,
			"min_error_value": -1,
			"min_threshold":   -1,
			"unit":            "",
		}

	default:
		return nil, errors.Wrap(errors.Errorf("unknown value type %q", valueType), "SensorValue.Make")
	}

	for k, v := range props {
		if err := obj.GetProps().Set(k, v); err != nil {
			return nil, errors.Wrapf(err, "SensorValue.Make(%s)", valueType)
		}
	}

	return obj, nil
}
