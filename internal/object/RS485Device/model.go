// Базовый тип для всех rs485-устройств

package RS485Device

import (
	"github.com/pkg/errors"
	"touchon-server/internal/model"
	"touchon-server/internal/object/RS485"
	"touchon-server/internal/objects"
	"touchon-server/internal/store/memstore"
	"touchon-server/lib/models"
)

type RS485Device interface {
	objects.Object
	DoAction(deviceAddr int, action RS485.Action, actionTries int, resultHandler RS485.ResultHandler, priority int) error
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
		model.CategoryRS485,
		"rs485_device",
		0,
		"Устройство RS485",
		props,
		nil,
		nil,
		nil,
		[]string{"rs485", "rs485_device"},
	)
	if err != nil {
		return nil, errors.Wrap(err, "RS485Device.MakeModel")
	}

	o := &RS485DeviceImpl{Object: impl}

	return o, nil
}

type RS485DeviceImpl struct {
	objects.Object
	bus RS485.RS485
}

func (o *RS485DeviceImpl) Start() error {
	if err := o.Object.Start(); err != nil {
		return errors.Wrap(err, "RS485DeviceImpl.Start")
	}

	parentID := o.GetParentID()
	if parentID == nil {
		return errors.Wrap(errors.Errorf("parent_id of %d is nil", o.GetID()), "RS485DeviceImpl.Start")
	}

	bus, err := memstore.I.GetObjectUnsafe(*parentID)
	if err != nil {
		return errors.Wrap(err, "RS485DeviceImpl.Start")
	}

	var ok bool
	o.bus, ok = bus.(*RS485.RS485Impl)
	if !ok {
		return errors.Wrap(errors.Errorf("parent object %d of object %d is not RS485.RS485Impl (%T)", o.GetParentID(), o.GetID(), bus), "RS485DeviceImpl.Start")
	}

	return nil
}

func (o *RS485DeviceImpl) Shutdown() error {
	if err := o.Object.Shutdown(); err != nil {
		return errors.Wrap(err, "RS485DeviceImpl.Shutdown")
	}

	return nil
}

func (o *RS485DeviceImpl) DoAction(deviceAddr int, action RS485.Action, actionTries int, resultHandler RS485.ResultHandler, priority int) error {
	if err := o.bus.DoAction(deviceAddr, action, actionTries, resultHandler, priority); err != nil {
		return errors.Wrap(err, "RS485DeviceImpl.DoAction")
	}

	return nil
}

func (o *RS485DeviceImpl) GetDefaultTries() int {
	return o.bus.GetDefaultTries()
}
