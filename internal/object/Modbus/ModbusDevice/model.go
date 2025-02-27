// Базовый тип для всех modbus-устройств

package ModbusDevice

import (
	"github.com/pkg/errors"
	"touchon-server/internal/model"
	"touchon-server/internal/object/Modbus"
	"touchon-server/internal/objects"
	"touchon-server/internal/store/memstore"
	"touchon-server/lib/models"
)

type ModbusDevice interface {
	objects.Object
	DoAction(deviceAddr int, action Modbus.Action, actionTries int, resultHandler Modbus.ResultHandler, priority int) error
	GetDefaultTries() int
}

func MakeModel(withChildren bool) (objects.Object, error) {
	props := []*objects.Prop{
		{
			Code:        "address",
			Name:        "Адрес устройства",
			Description: "",
			Item: &models.Item{
				Type:         models.DataTypeInt,
				DefaultValue: 0,
			},
			Required:   objects.True(),
			Editable:   objects.True(),
			Visible:    objects.True(),
			CheckValue: objects.Between(1, 0xFFFF),
		},
		{
			Code:        "update_interval",
			Name:        "Интервал опроса (10s, 1m etc)",
			Description: "Интервал опроса устройства",
			Item: &models.Item{
				Type:         models.DataTypeString,
				DefaultValue: "1m",
			},
			Required:   objects.True(),
			Editable:   objects.True(),
			Visible:    objects.True(),
			CheckValue: objects.AboveOrEqual1(),
		},
	}

	impl, err := objects.NewObjectModelImpl(
		model.CategoryModbus,
		"modbus_device",
		0,
		"Устройство Modbus",
		props,
		nil,
		nil,
		nil,
		[]string{"modbus", "modbus_device"},
	)
	if err != nil {
		return nil, errors.Wrap(err, "ModbusDevice.MakeModel")
	}

	o := &ModbusDeviceImpl{Object: impl}

	return o, nil
}

type ModbusDeviceImpl struct {
	objects.Object
	modbus Modbus.Modbus
}

func (o *ModbusDeviceImpl) Start() error {
	if err := o.Object.Start(); err != nil {
		return errors.Wrap(err, "ModbusDeviceImpl.Start")
	}

	parentID := o.GetParentID()
	if parentID == nil {
		return errors.Wrap(errors.Errorf("parent_id of %d is nil", o.GetID()), "ModbusDeviceImpl.Start")
	}

	bus, err := memstore.I.GetObjectUnsafe(*parentID)
	if err != nil {
		return errors.Wrap(err, "ModbusDeviceImpl.Start")
	}

	var ok bool
	o.modbus, ok = bus.(*Modbus.ModbusImpl)
	if !ok {
		return errors.Wrap(errors.Errorf("parent object %d of object %d is not Modbus.ModbusImpl (%T)", o.GetParentID(), o.GetID(), bus), "ModbusDeviceImpl.Start")
	}

	return nil
}

func (o *ModbusDeviceImpl) Shutdown() error {
	if err := o.Object.Shutdown(); err != nil {
		return errors.Wrap(err, "ModbusDeviceImpl.Shutdown")
	}

	return nil
}

func (o *ModbusDeviceImpl) DoAction(deviceAddr int, action Modbus.Action, actionTries int, resultHandler Modbus.ResultHandler, priority int) error {
	if err := o.modbus.DoAction(deviceAddr, action, actionTries, resultHandler, priority); err != nil {
		return errors.Wrap(err, "ModbusDeviceImpl.DoAction")
	}

	return nil
}

func (o *ModbusDeviceImpl) GetDefaultTries() int {
	return o.modbus.GetDefaultTries()
}
