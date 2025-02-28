package main

import (
	"strconv"

	"github.com/pkg/errors"
	"touchon-server/internal/model"
	"touchon-server/internal/objects"
	"touchon-server/internal/store"
)

func prepareDB() error {
	if err := createServerObject(); err != nil {
		return errors.Wrap(err, "prepareDB")
	}

	return nil
}

// Создаем единственный экземпляр типа server/server и две шины RS485
func createServerObject() error {
	count, err := store.I.ObjectRepository().GetTotal(map[string]interface{}{"category": model.CategoryServer, "type": "server"}, nil)
	if err != nil {
		return errors.Wrap(err, "createServerObject")
	}

	switch {
	case count > 1:
		return errors.Wrap(errors.New("Кол-во объектов server/server больше 1"), "createServerObject")
	case count == 1:
		return nil
	}

	server, err := createObject(model.CategoryServer, "server", nil, true)
	if err != nil {
		return errors.Wrap(err, "createServerObject")
	}

	serverID := server.GetID()

	bus0, err := createObject(model.CategoryRS485, "bus", &serverID, true)
	if err != nil {
		return errors.Wrap(err, "createServerObject")
	}

	bus1, err := createObject(model.CategoryRS485, "bus", &serverID, true)
	if err != nil {
		return errors.Wrap(err, "createServerObject")
	}

	for i, bus := range []objects.Object{bus0, bus1} {
		props := bus.GetProps()

		if err := props.Set("connection_string", "rtu:///dev/ttyUSB"+strconv.Itoa(i)); err != nil {
			return errors.Wrap(err, "createServerObject")
		}

		if err := props.Set("speed", "9600"); err != nil {
			return errors.Wrap(err, "createServerObject")
		}

		if err := props.Set("data_bits", 8); err != nil {
			return errors.Wrap(err, "createServerObject")
		}

		if err := props.Set("parity", "0"); err != nil {
			return errors.Wrap(err, "createServerObject")
		}

		if err := props.Set("stop_bits", "1"); err != nil {
			return errors.Wrap(err, "createServerObject")
		}

		if err := bus.Save(); err != nil {
			return errors.Wrap(err, "createServerObject")
		}
	}

	return nil
}

func createObject(objCat model.Category, objType string, parentID *int, enabled bool) (objects.Object, error) {
	objModel, err := objects.LoadObject(0, objCat, objType, true)
	if err != nil {
		return nil, errors.Wrap(err, "createObject")
	}

	objModel.SetParentID(parentID)
	objModel.SetEnabled(enabled)

	// Выставляем сначала значения по умолчанию
	for _, p := range objModel.GetProps().GetAll().GetValueList() {
		if p.GetValue() == nil && p.DefaultValue != nil {
			if err := p.SetValue(p.DefaultValue); err != nil {
				return nil, errors.Wrap(err, "createObject")
			}
		}
	}

	if err := objModel.GetProps().Check(); err != nil {
		return nil, errors.Wrap(err, "createObject")
	}

	if err := objModel.Save(); err != nil {
		return nil, errors.Wrap(err, "createObject")
	}

	return objModel, nil
}
