package bh1750

import (
	"strconv"
	"time"

	"github.com/pkg/errors"
	"touchon-server/internal/model"
	"touchon-server/internal/object/Sensor"
	"touchon-server/internal/object/SensorValue"
	"touchon-server/internal/objects"
	"touchon-server/internal/store"
	"touchon-server/lib/mqtt/messages"
)

func init() {
	_ = objects.Register(MakeModel)
}

func MakeModel() (objects.Object, error) {
	baseObj, err := Sensor.MakeModel()
	if err != nil {
		return nil, errors.Wrap(err, "bh1750.MakeModel")
	}

	obj := &SensorModel{}
	obj.SensorModel = baseObj.(*Sensor.SensorModel)

	obj.SetType("bh1750")
	obj.SetName("BH1750 Датчик интенсивности света")
	obj.SetTags("bh1750", string(SensorValue.TypeIllumination))

	if err := obj.GetProps().Set("interface", "I2C"); err != nil {
		return nil, errors.Wrap(err, "bh1750.MakeModel")
	}

	illum, err := SensorValue.Make(SensorValue.TypeIllumination)
	if err != nil {
		return nil, errors.Wrap(err, "bh1750.MakeModel")
	}

	obj.GetChildren().Add(illum)

	check, err := objects.NewMethod("check", "Опрашивает датчик, обновляет показания датчика в БД", nil, obj.Check)
	if err != nil {
		return nil, errors.Wrap(err, "bh1750.MakeModel")
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

	params := map[string]string{
		"scl":     sclPort,
		"i2c_dev": o.GetType(),
	}

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

	return map[SensorValue.Type]float32{
		SensorValue.TypeIllumination: float32(v),
	}, nil
}

func (o *SensorModel) Check(args map[string]interface{}) ([]messages.Message, error) {
	msgs, err := o.SensorModel.Check(o.getValues)
	if err != nil {
		return nil, errors.Wrap(err, "Check")
	}

	return msgs, nil
}
