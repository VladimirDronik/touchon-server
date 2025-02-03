package sqlstore

import (
	"github.com/pkg/errors"
	"touchon-server/internal/model"
)

type DeviceRepository struct {
	store *Store
}

func (o *DeviceRepository) GetControllerByName(name string) (*model.StoreObject, error) {
	obj := &model.StoreObject{}

	err := o.store.db.
		Where("name = ?", name).
		Find(obj).Error

	if err != nil {
		return nil, errors.Wrap(err, "GetControllerByName")
	}

	if obj.ID == 0 {
		return nil, errors.Wrap(errNotFound, "GetControllerByName")
	}

	return obj, nil
}
