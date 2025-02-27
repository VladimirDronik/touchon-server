package bmp280

import (
	"regexp"
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
		return nil, errors.Wrap(err, "bmp280.MakeModel")
	}

	obj := &SensorModel{}
	obj.SensorModel = baseObj.(*Sensor.SensorModel)

	obj.SetType("bmp280")
	obj.SetName("BMP280 Датчик температуры и влажности")
	obj.SetTags("bmp280", SensorValue.TypeTemperature, SensorValue.TypeHumidity)
	obj.SetGetValuesFunc(obj.getValues)

	if err := obj.GetProps().Set("interface", "I2C"); err != nil {
		return nil, errors.Wrap(err, "bmp280.MakeModel")
	}

	if !withChildren {
		return obj, nil
	}

	temp, err := SensorValue.Make(SensorValue.TypeTemperature, withChildren)
	if err != nil {
		return nil, errors.Wrap(err, "bmp280.MakeModel")
	}

	hum, err := SensorValue.Make(SensorValue.TypeHumidity, withChildren)
	if err != nil {
		return nil, errors.Wrap(err, "bmp280.MakeModel")
	}

	obj.GetChildren().Add(temp, hum)

	return obj, nil
}

type SensorModel struct {
	*Sensor.SensorModel
}

var valuesPatt = regexp.MustCompile(`^temp:([\d.]+)/hum:([\d.]+)$`)

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

	params := map[string]string{
		"scl":     sclPort,
		"i2c_dev": "bmx280",
		"i2c_par": "3",
	}

	r := make(map[SensorValue.Type]float32, 2)

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

	for i, t := range []SensorValue.Type{SensorValue.TypeTemperature, SensorValue.TypeHumidity} {
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
