package Sensor

import (
	"fmt"
	"github.com/pkg/errors"
	"touchon-server/internal/objects"
)

func MakeAlarmMessage(valueObj objects.Object) (string, error) {
	value, err := valueObj.GetProps().GetFloatValue("value")
	if err != nil {
		return "", errors.Wrap(err, "MakeAlarmMessage")
	}

	var unit string
	unit, err = valueObj.GetProps().GetStringValue("unit")
	if err != nil {
		return "", errors.Wrap(err, "MakeAlarmMessage")
	}

	var minThreshold float32
	minThreshold, err = valueObj.GetProps().GetFloatValue("min_threshold")
	if err != nil {
		return "", errors.Wrap(err, "MakeAlarmMessage")
	}

	var maxThreshold float32
	maxThreshold, err = valueObj.GetProps().GetFloatValue("max_threshold")
	if err != nil {
		return "", errors.Wrap(err, "MakeAlarmMessage")
	}

	return fmt.Sprintf(
		"%s: %.2f%s не соответствует пороговым значениям от %.2f до %.2f",
		valueObj.GetName(), value, unit, minThreshold, maxThreshold,
	), nil
}
