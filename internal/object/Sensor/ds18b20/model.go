package ds18b20

import (
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
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
		return nil, errors.Wrap(err, "ds18b20.MakeModel")
	}

	obj := &SensorModel{}
	obj.SensorModel = baseObj.(*Sensor.SensorModel)

	obj.SetType("ds18b20")
	obj.SetName("DS18B20 Датчик температуры")
	obj.SetTags("ds18b20", string(SensorValue.TypeTemperature))

	iface, err := obj.GetProps().Get("interface")
	if err != nil {
		return nil, errors.Wrap(err, "ds18b20.MakeModel")
	}

	// Делаем св-во редактируемым
	iface.Editable = objects.NewCondition()
	if err := iface.SetValue("1W"); err != nil {
		return nil, errors.Wrap(err, "ds18b20.MakeModel")
	}

	temp, err := SensorValue.Make(SensorValue.TypeTemperature)
	if err != nil {
		return nil, errors.Wrap(err, "ds18b20.MakeModel")
	}

	obj.GetChildren().Add(temp)
	obj.SetGetValuesFunc(obj.getValues)

	return obj, nil
}

type SensorModel struct {
	*Sensor.SensorModel
}

func (o *SensorModel) getValues(timeout time.Duration) (map[SensorValue.Type]float32, error) {
	ifaceType, err := o.GetProps().GetStringValue("interface")
	if err != nil {
		return nil, errors.Wrap(err, "getValues")
	}

	addr, err := o.GetProps().GetStringValue("address")
	if err != nil {
		return nil, errors.Wrap(err, "getValues")
	}

	portAndAddress := strings.Split(addr, ";")

	portObjectID, err := strconv.Atoi(portAndAddress[0])
	if err != nil {
		return nil, errors.Wrap(err, "getValues")
	}

	portObj, err := objects.LoadPort(portObjectID, false)
	if err != nil {
		return nil, errors.Wrap(err, "getValues")
	}

	var value string

	switch ifaceType {

	// если датчик размещен на порту контроллера как отдельный датчик 1-wire
	case "1W":
		params := map[string]string{"f": "s"}

		if value, err = portObj.GetPortState("get", params, timeout); err != nil {
			return nil, errors.Wrap(err, "getValues")
		}

	// если датчик размещен на шине 1-wire-bus
	case "1WBUS":
		address := portAndAddress[1]
		if address == "" {
			return nil, errors.Wrap(errors.New("не указан адрес термометра"), "getValues")
		}

		thermometersString, err := portObj.GetPortState("get", nil, timeout)
		if err != nil {
			return nil, errors.Wrap(err, "getValues")
		}

		thermometers := strings.Split(thermometersString, ";")
		for _, thermometer := range thermometers {
			addrAndValue := strings.Split(thermometer, ":")
			if len(addrAndValue) != 2 {
				continue
			}

			if addrAndValue[0] == address {
				value = addrAndValue[1]
				break
			}
		}

	default:
		return nil, errors.Wrap(errors.Errorf("unsupported interface %q", ifaceType), "getValues")
	}

	if value == "na" {
		return nil, errors.Wrap(errors.New("sensor is faulty"), "getValues")
	}

	v, err := strconv.ParseFloat(value, 32)
	if err != nil {
		return nil, errors.Wrap(err, "getValues")
	}

	return map[SensorValue.Type]float32{
		SensorValue.TypeTemperature: float32(v),
	}, nil
}
