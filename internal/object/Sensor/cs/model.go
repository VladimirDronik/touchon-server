package cs

import (
	"strconv"
	"time"

	"github.com/pkg/errors"
	"touchon-server/internal/model"
	"touchon-server/internal/object/Sensor"
	"touchon-server/internal/object/SensorValue"
	"touchon-server/internal/objects"
)

func init() {
	_ = objects.Register(MakeModel)
}

func MakeModel() (objects.Object, error) {
	baseObj, err := Sensor.MakeModel()
	if err != nil {
		return nil, errors.Wrap(err, "cs.MakeModel")
	}

	obj := &SensorModel{}
	obj.SensorModel = baseObj.(*Sensor.SensorModel)

	obj.SetType("cs")
	obj.SetName("CS Датчик тока")
	obj.SetTags("cs", string(SensorValue.TypeCurrent))

	if err := obj.GetProps().Set("interface", "ADC"); err != nil {
		return nil, errors.Wrap(err, "cs.MakeModel")
	}

	cur, err := SensorValue.Make(SensorValue.TypeCurrent)
	if err != nil {
		return nil, errors.Wrap(err, "cs.MakeModel")
	}

	obj.GetChildren().Add(cur)
	obj.SetGetValuesFunc(obj.getValues)

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

	portObjectID, err := strconv.Atoi(addr)
	if err != nil {
		return nil, errors.Wrap(err, "getValues")
	}

	port, err := objects.LoadPort(portObjectID, model.ChildTypeNobody)
	if err != nil {
		return nil, errors.Wrap(err, "getValues")
	}

	value, err := port.GetPortState("get", nil, timeout)
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
		SensorValue.TypeCurrent: float32(v),
	}, nil
}
