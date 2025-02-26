package main

import (
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

func createServerObject() error {
	count, err := store.I.ObjectRepository().GetTotal(map[string]interface{}{"category": model.CategoryServer, "type": "server"}, nil, model.ChildTypeNobody)
	if err != nil {
		return errors.Wrap(err, "createServerObject")
	}

	switch {
	case count > 1:
		return errors.Wrap(errors.New("Кол-во объектов server/server больше 1"), "createServerObject")
	case count == 1:
		return nil
	}

	serverID, err := createObject(model.CategoryServer, "server", nil, true)
	if err != nil {
		return errors.Wrap(err, "createServerObject")
	}

	for i := 0; i < 2; i++ {
		if _, err := createObject(model.CategoryModbus, "modbus", &serverID, false); err != nil {
			return errors.Wrap(err, "createServerObject")
		}
	}

	return nil
}

func createObject(objCat model.Category, objType string, parentID *int, enabled bool) (int, error) {
	objModel, err := objects.LoadObject(0, objCat, objType, model.ChildTypeAll)
	if err != nil {
		return 0, errors.Wrap(err, "createObject")
	}

	objModel.SetParentID(parentID)
	objModel.SetEnabled(enabled)

	// Выставляем сначала значения по умолчанию
	for _, p := range objModel.GetProps().GetAll().GetValueList() {
		if p.GetValue() == nil && p.DefaultValue != nil {
			if err := p.SetValue(p.DefaultValue); err != nil {
				return 0, errors.Wrap(err, "createObject")
			}
		}
	}

	if err := objModel.GetProps().Check(); err != nil {
		return 0, errors.Wrap(err, "createObject")
	}

	if err := objModel.Save(); err != nil {
		return 0, errors.Wrap(err, "createObject")
	}

	return objModel.GetID(), nil
}
