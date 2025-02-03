// Кондиционеры с поддержкой управления через modbus-шлюзы Onokom

package Conditioner

import (
	"github.com/pkg/errors"
	"touchon-server/internal/model"
	"touchon-server/internal/object/Onokom/ModbusGW"
	"touchon-server/internal/objects"
)

func init() {
	// Регистрируем для каждого поддерживаемого Modbus-шлюза свой тип кондиционера.
	for gwModelCode, name := range ModbusGW.SupportedGateways {
		_ = objects.Register(func() (objects.Object, error) {
			return MakeModel(gwModelCode, name)
		})
	}
}

func MakeModel(gwModelCode, name string) (objects.Object, error) {
	baseObj, err := ModbusGW.MakeModel(gwModelCode)
	if err != nil {
		return nil, errors.Wrap(err, "Conditioner.MakeModel")
	}

	obj := &DeviceModel{}
	obj.GatewayModel = baseObj.(*ModbusGW.GatewayModel)

	obj.SetCategory(model.CategoryConditioner)
	obj.SetType("onokom/" + gwModelCode)
	obj.SetName("Кондиционер (Onokom/" + name + ")")
	obj.SetTags(string(model.CategoryConditioner), "onokom", "modbus", gwModelCode)

	return obj, nil
}

type DeviceModel struct {
	*ModbusGW.GatewayModel
}
