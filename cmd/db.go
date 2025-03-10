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

	server, err := createObject(model.CategoryServer, "server", nil, true, nil)
	if err != nil {
		return errors.Wrap(err, "createServerObject")
	}

	serverID := server.GetID()

	for i := 0; i < 2; i++ {
		props := map[string]interface{}{
			"connection_string": "rtu:///dev/ttyUSB" + strconv.Itoa(i),
			"speed":             "9600",
			"data_bits":         8,
			"parity":            "0",
			"stop_bits":         "1",
		}

		if _, err := createObject(model.CategoryRS485, "bus", &serverID, true, props); err != nil {
			return errors.Wrap(err, "createServerObject")
		}
	}

	return nil
}

func createObject(objCat model.Category, objType string, parentID *int, enabled bool, props map[string]interface{}) (objects.Object, error) {
	objModel, err := objects.LoadObject(0, objCat, objType, true)
	if err != nil {
		return nil, errors.Wrapf(err, "createObject(%s, %s)", objCat, objType)
	}

	objModel.SetParentID(parentID)
	objModel.SetEnabled(enabled)

	// Выставляем сначала значения по умолчанию
	for _, p := range objModel.GetProps().GetAll().GetValueList() {
		if p.GetValue() == nil && p.DefaultValue != nil {
			if err := p.SetValue(p.DefaultValue); err != nil {
				return nil, errors.Wrapf(err, "createObject(%s, %s)", objCat, objType)
			}
		}
	}

	for k, v := range props {
		if err := objModel.GetProps().Set(k, v); err != nil {
			return nil, errors.Wrapf(err, "createObject(%s, %s)", objCat, objType)
		}
	}

	if err := objModel.GetProps().Check(); err != nil {
		return nil, errors.Wrapf(err, "createObject(%s, %s)", objCat, objType)
	}

	if err := objModel.Save(); err != nil {
		return nil, errors.Wrapf(err, "createObject(%s, %s)", objCat, objType)
	}

	return objModel, nil
}
