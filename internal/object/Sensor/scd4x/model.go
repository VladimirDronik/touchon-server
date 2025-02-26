package scd4x

import (
	"regexp"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"touchon-server/internal/model"
	"touchon-server/internal/object/Sensor"
	"touchon-server/internal/object/SensorValue"
	"touchon-server/internal/objects"
	"touchon-server/internal/store"
)

func init() {
	_ = objects.Register(MakeModel)
}

func MakeModel() (objects.Object, error) {
	baseObj, err := Sensor.MakeModel()
	if err != nil {
		return nil, errors.Wrap(err, "scd4x.MakeModel")
	}

	obj := &SensorModel{}
	obj.SensorModel = baseObj.(*Sensor.SensorModel)

	obj.SetType("scd4x")
	obj.SetName("SCD4x Датчик уровня CO2, температуры и влажности")
	obj.SetTags("scd4x",
		SensorValue.TypeTemperature,
		SensorValue.TypeHumidity,
		SensorValue.TypeCO2,
	)

	if err := obj.GetProps().Set("interface", "I2C"); err != nil {
		return nil, errors.Wrap(err, "scd4x.MakeModel")
	}

	temp, err := SensorValue.Make(SensorValue.TypeTemperature)
	if err != nil {
		return nil, errors.Wrap(err, "scd4x.MakeModel")
	}

	hum, err := SensorValue.Make(SensorValue.TypeHumidity)
	if err != nil {
		return nil, errors.Wrap(err, "scd4x.MakeModel")
	}

	co2, err := SensorValue.Make(SensorValue.TypeCO2)
	if err != nil {
		return nil, errors.Wrap(err, "scd4x.MakeModel")
	}

	obj.GetChildren().Add(temp, hum, co2)
	obj.SetGetValuesFunc(obj.getValues)

	return obj, nil
}

type SensorModel struct {
	*Sensor.SensorModel
}

var valuesPatt = regexp.MustCompile(`^co2:([\d.]+)/temp:([\d.]+)/hum:([\d.]+)$`)

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

	port, err := objects.LoadPort(sdaPortObjectID, model.ChildTypeNobody)
	if err != nil {
		return nil, errors.Wrap(err, "getValues")
	}

	params := map[string]string{
		"scl":     sclPort,
		"i2c_dev": o.GetType(),
		"i2c_par": "0",
	}

	value, err := port.GetPortState("", params, timeout)
	if err != nil {
		return nil, errors.Wrap(err, "getValues")
	}

	if value == "na" {
		return nil, errors.Wrap(errors.New("sensor is faulty"), "getValues")
	}

	values := valuesPatt.FindAllStringSubmatch(value, -1)
	if len(values) != 1 || len(values[0]) != 4 {
		return nil, errors.Wrap(errors.Errorf("bad values format %q", value), "getValues")
	}

	r := make(map[SensorValue.Type]float32, 2)

	for i, t := range []SensorValue.Type{SensorValue.TypeCO2, SensorValue.TypeTemperature, SensorValue.TypeHumidity} {
		value := values[0][i+1]
		if value == "na" {
			return nil, errors.Wrap(errors.New("sensor is faulty"), "getValues")
		}

		v, err := strconv.ParseFloat(value, 32)
		if err != nil {
			return nil, errors.Wrap(err, "getValues")
		}

		r[t] = float32(v)
	}

	return r, nil
}
