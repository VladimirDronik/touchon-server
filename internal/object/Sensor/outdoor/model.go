package outdoor

import (
	"time"

	"github.com/pkg/errors"
	"touchon-server/internal/model"
	"touchon-server/internal/object/Sensor"
	"touchon-server/internal/object/SensorValue"
	"touchon-server/internal/objects"
	"touchon-server/lib/interfaces"
)

func init() {
	_ = objects.Register(MakeModel)
}

func MakeModel() (objects.Object, error) {
	baseObj, err := Sensor.MakeModel()
	if err != nil {
		return nil, errors.Wrap(err, "outdoor.MakeModel")
	}

	obj := &SensorModel{}
	obj.SensorModel = baseObj.(*Sensor.SensorModel)

	obj.SetType("outdoor")
	obj.SetName("Outdoor Уличный датчик (BH1750+BME280)")
	obj.SetTags(
		"outdoor", "bh1750", "bme280",
		string(SensorValue.TypeTemperature),
		string(SensorValue.TypePressure),
		string(SensorValue.TypeHumidity),
		string(SensorValue.TypeIllumination),
	)

	if err := obj.GetProps().Set("interface", "I2C"); err != nil {
		return nil, errors.Wrap(err, "outdoor.MakeModel")
	}

	temp, err := SensorValue.Make(SensorValue.TypeTemperature)
	if err != nil {
		return nil, errors.Wrap(err, "outdoor.MakeModel")
	}

	pres, err := SensorValue.Make(SensorValue.TypePressure)
	if err != nil {
		return nil, errors.Wrap(err, "outdoor.MakeModel")
	}

	hum, err := SensorValue.Make(SensorValue.TypeHumidity)
	if err != nil {
		return nil, errors.Wrap(err, "outdoor.MakeModel")
	}

	illum, err := SensorValue.Make(SensorValue.TypeIllumination)
	if err != nil {
		return nil, errors.Wrap(err, "outdoor.MakeModel")
	}

	obj.GetChildren().Add(temp, pres, hum, illum)

	check, err := objects.NewMethod("check", "Опрашивает датчик, обновляет показания датчика в БД", nil, obj.Check)
	if err != nil {
		return nil, errors.Wrap(err, "outdoor.MakeModel")
	}

	obj.GetMethods().Add(check)

	return obj, nil
}

type SensorModel struct {
	*Sensor.SensorModel
}

func (o *SensorModel) getValues(timeout time.Duration) (map[SensorValue.Type]float32, error) {
	addr, err := o.GetProps().GetStringValue("address")
	if err != nil {
		return nil, errors.Wrap(err, "getValues")
	}

	// Temperature, Humidity, Pressure:
	bme280, err := objects.GetObjectModel(model.CategorySensor, "bme280")
	if err != nil {
		return nil, errors.Wrap(err, "getValues")
	}

	// Illumination:
	bh1750, err := objects.GetObjectModel(model.CategorySensor, "bh1750")
	if err != nil {
		return nil, errors.Wrap(err, "getValues")
	}

	for _, obj := range []objects.Object{bme280, bh1750} {
		if err := obj.GetProps().Set("address", addr); err != nil {
			return nil, errors.Wrap(err, "getValues")
		}

		obj.GetChildren().DeleteAll()
	}

	for _, child := range o.GetChildren().GetAll() {
		switch child.GetType() {
		case SensorValue.TypeTemperature, SensorValue.TypePressure, SensorValue.TypeHumidity:
			bme280.GetChildren().Add(child)

		case SensorValue.TypeIllumination:
			bh1750.GetChildren().Add(child)
		}
	}

	r := make(map[SensorValue.Type]float32, 2)
	for _, obj := range []objects.Object{bme280, bh1750} {
		m, err := obj.GetMethods().Get("check")
		if err != nil {
			return nil, errors.Wrap(err, "getValues")
		}

		if _, err := m.Func(nil); err != nil {
			return nil, errors.Wrap(err, "getValues")
		}

		for _, child := range obj.GetChildren().GetAll() {
			v, err := child.GetProps().GetFloatValue("value")
			if err != nil {
				return nil, errors.Wrap(err, "getValues")
			}

			r[SensorValue.Type(child.GetType())] = v
		}
	}

	return r, nil
}

func (o *SensorModel) Check(args map[string]interface{}) ([]interfaces.Message, error) {
	msgs, err := o.SensorModel.Check(o.getValues)
	if err != nil {
		return nil, errors.Wrap(err, "Check")
	}

	return msgs, nil
}
