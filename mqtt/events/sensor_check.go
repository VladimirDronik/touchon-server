package events

import (
	"github.com/VladimirDronik/touchon-server/mqtt/messages"
	"github.com/pkg/errors"
)

func NewOnSensorCheck(targetID int, targetType messages.TargetType, values map[string]float32) (messages.Message, error) {
	payload := make(map[string]interface{}, len(values))
	for k, v := range values {
		payload[k] = v
	}

	impl, err := NewEvent("onCheck", targetID, targetType, payload)
	if err != nil {
		return nil, errors.Wrap(err, "NewOnSensorCheck")
	}

	return &OnSensorCheck{Message: impl}, nil
}

type OnSensorCheck struct {
	messages.Message
}

func (o *OnSensorCheck) GetTemperature() (float32, error) {
	v, err := o.GetFloatValue("temperature")
	if err != nil {
		return 0, errors.Wrap(err, "GetTemperature")
	}

	return v, nil
}

func (o *OnSensorCheck) GetHumidity() (float32, error) {
	v, err := o.GetFloatValue("humidity")
	if err != nil {
		return 0, errors.Wrap(err, "GetHumidity")
	}

	return v, nil
}

func (o *OnSensorCheck) GetPressure() (float32, error) {
	v, err := o.GetFloatValue("pressure")
	if err != nil {
		return 0, errors.Wrap(err, "GetPressure")
	}

	return v, nil
}

func (o *OnSensorCheck) GetIllumination() (float32, error) {
	v, err := o.GetFloatValue("illumination")
	if err != nil {
		return 0, errors.Wrap(err, "GetIllumination")
	}

	return v, nil
}

func (o *OnSensorCheck) GetCurrent() (float32, error) {
	v, err := o.GetFloatValue("current")
	if err != nil {
		return 0, errors.Wrap(err, "GetCurrent")
	}

	return v, nil
}

func (o *OnSensorCheck) GetVoltage() (float32, error) {
	v, err := o.GetFloatValue("voltage")
	if err != nil {
		return 0, errors.Wrap(err, "GetVoltage")
	}

	return v, nil
}

func (o *OnSensorCheck) GetCO2() (float32, error) {
	v, err := o.GetFloatValue("co2")
	if err != nil {
		return 0, errors.Wrap(err, "GetCO2")
	}

	return v, nil
}

func (o *OnSensorCheck) SetTemperature(v float32) {
	o.GetPayload()["temperature"] = v
}

func (o *OnSensorCheck) SetHumidity(v float32) {
	o.GetPayload()["humidity"] = v
}

func (o *OnSensorCheck) SetPressure(v float32) {
	o.GetPayload()["pressure"] = v
}

func (o *OnSensorCheck) SetIllumination(v float32) {
	o.GetPayload()["illumination"] = v
}

func (o *OnSensorCheck) SetCurrent(v float32) {
	o.GetPayload()["current"] = v
}

func (o *OnSensorCheck) SetVoltage(v float32) {
	o.GetPayload()["voltage"] = v
}

func (o *OnSensorCheck) SetCO2(v float32) {
	o.GetPayload()["co2"] = v
}
