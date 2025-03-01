package htu31d

import (
	"strconv"
	"time"

	"github.com/pkg/errors"
	"touchon-server/internal/object/Sensor"
	"touchon-server/internal/object/SensorValue"
	"touchon-server/internal/objects"
	"touchon-server/internal/store"
)

func init() {
	_ = objects.Register(MakeModel)
}

func MakeModel(withChildren bool) (objects.Object, error) {
	baseObj, err := Sensor.MakeModel(withChildren)
	if err != nil {
		return nil, errors.Wrap(err, "htu31d.MakeModel")
	}

	obj := &SensorModel{}
	obj.SensorModel = baseObj.(*Sensor.SensorModel)

	obj.SetType("htu31d")
	obj.SetName("HTU31D Датчик температуры и влажности")
	obj.SetTags("htu31d", SensorValue.TypeTemperature, SensorValue.TypeHumidity)
	obj.SetGetValuesFunc(obj.getValues)

	if err := obj.GetProps().Set("interface", "I2C"); err != nil {
		return nil, errors.Wrap(err, "htu31d.MakeModel")
	}

	if !withChildren {
		return obj, nil
	}

	temp, err := SensorValue.Make(SensorValue.TypeTemperature, withChildren)
	if err != nil {
		return nil, errors.Wrap(err, "htu31d.MakeModel")
	}

	hum, err := SensorValue.Make(SensorValue.TypeHumidity, withChildren)
	if err != nil {
		return nil, errors.Wrap(err, "htu31d.MakeModel")
	}

	obj.GetChildren().Add(temp, hum)

	return obj, nil
}

type SensorModel struct {
	*Sensor.SensorModel
}

func (o *SensorModel) getValues(timeout time.Duration) (map[SensorValue.Type]float32, error) {
	sdaPortObjectID, sclPortObjectID, err := o.ParseI2CAddress()
	if err != nil {
		return nil, errors.Wrap(err, "getValues")
	}

	// Для объекта порта SCL ищем свойство с номером порта, чтобы можно было использовать в строке запроса к контроллеру
	sclPort, err := store.I.ObjectRepository().GetProp(sclPortObjectID, "number")
	if err != nil {
		return nil, errors.Wrap(err, "getValues")
	}

	port, err := objects.LoadPort(sdaPortObjectID, false)
	if err != nil {
		return nil, errors.Wrap(err, "getValues")
	}

	i2cParams := map[SensorValue.Type]string{
		SensorValue.TypeHumidity:    "0",
		SensorValue.TypeTemperature: "1",
	}

	params := map[string]string{
		"scl":     sclPort,
		"i2c_dev": o.GetType(),
	}

	r := make(map[SensorValue.Type]float32, 2)
	for valueType, i2cParam := range i2cParams {
		params["i2c_par"] = i2cParam

		value, err := port.GetPortState("", params, timeout)
		if err != nil {
			return nil, errors.Wrap(err, "getValues")
		}

		if value == "na" {
			return nil, errors.Wrap(errors.New("sensor is faulty"), "getValues")
		}

		v, err := strconv.ParseFloat(value, 32)
		if err != nil {
			return nil, errors.Wrap(err, "getValues")
		}

		r[valueType] = float32(v)
	}

	return r, nil
}
