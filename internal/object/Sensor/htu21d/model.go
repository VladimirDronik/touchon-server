package htu21d

import (
	"strconv"
	"time"

	"github.com/pkg/errors"
	"touchon-server/internal/g"
	"touchon-server/internal/model"
	"touchon-server/internal/object/Sensor"
	"touchon-server/internal/object/SensorValue"
	"touchon-server/internal/objects"
	"touchon-server/internal/store"
	"touchon-server/lib/interfaces"
)

func init() {
	_ = objects.Register(MakeModel)
}

func MakeModel() (objects.Object, error) {
	baseObj, err := Sensor.MakeModel()
	if err != nil {
		return nil, errors.Wrap(err, "htu21d.MakeModel")
	}

	obj := &SensorModel{}
	obj.SensorModel = baseObj.(*Sensor.SensorModel)

	obj.SetType("htu21d")
	obj.SetName("HTU21D Датчик температуры и влажности")
	obj.SetTags("htu21d", string(SensorValue.TypeTemperature), string(SensorValue.TypeHumidity))

	if err := obj.GetProps().Set("interface", "I2C"); err != nil {
		return nil, errors.Wrap(err, "htu21d.MakeModel")
	}

	temp, err := SensorValue.Make(SensorValue.TypeTemperature)
	if err != nil {
		return nil, errors.Wrap(err, "htu21d.MakeModel")
	}

	hum, err := SensorValue.Make(SensorValue.TypeHumidity)
	if err != nil {
		return nil, errors.Wrap(err, "htu21d.MakeModel")
	}

	obj.GetChildren().Add(temp, hum)

	check, err := objects.NewMethod("check", "Опрашивает датчик, обновляет показания датчика в БД", nil, obj.Check)
	if err != nil {
		return nil, errors.Wrap(err, "htu21d.MakeModel")
	}

	obj.GetMethods().Add(check)

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

	port, err := objects.LoadPort(sdaPortObjectID, model.ChildTypeNobody)
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

func (o *SensorModel) Check(args map[string]interface{}) ([]interfaces.Message, error) {
	msgs, err := o.SensorModel.Check(o.getValues)
	if err != nil {
		return nil, errors.Wrap(err, "Check")
	}

	return msgs, nil
}

func (o *SensorModel) Start() error {
	if err := o.SensorModel.Start(); err != nil {
		return errors.Wrap(err, "htu21d.SensorModel.Start")
	}

	updateInterval, err := o.GetProps().GetIntValue("update_interval")
	if err != nil {
		return errors.Wrap(err, "htu21d.SensorModel.Start")
	}

	o.SetTimer(time.Duration(updateInterval)*time.Second, o.check)
	o.GetTimer().Start()

	g.Logger.Debugf("htu21d(%d) started", o.GetID())

	return nil
}

func (o *SensorModel) Shutdown() error {
	if err := o.SensorModel.Shutdown(); err != nil {
		return errors.Wrap(err, "htu21d.SensorModel.Shutdown")
	}

	g.Logger.Debugf("htu21d(%d) stopped", o.GetID())

	return nil
}

func (o *SensorModel) check() {
	g.Logger.Debugf("htu21d(%d) check", o.GetID())

	// TODO....

	o.GetTimer().Reset()
}
