package sqlstore

import (
	"math"

	"github.com/pkg/errors"
	"touchon-server/internal/model"
)

type Devices struct {
	store *Store
}

// GetSensor Отдает данные для датчика
func (o *Devices) GetSensor(itemID int) (*model.Sensor, error) {
	r := &model.Sensor{}
	if err := o.store.db.First(&r, "view_item_id = ?", itemID).Error; err != nil {
		return nil, errors.Wrap(err, "GetSensor")
	}

	return r, nil
}

// GetDimmer Отдает данные для димера
func (o *Devices) GetDimmer(itemID int) (*model.Dimmer, error) {
	r := &model.Dimmer{}
	if err := o.store.db.First(&r, "view_item_id = ?", itemID).Error; err != nil {
		return nil, errors.Wrap(err, "GetDimmer")
	}

	return r, nil
}

func (o *Devices) SetSensorValue(itemID int, value float32) error {
	if err := o.SetFieldValue("sensors", itemID, "current", float32(math.Round(float64(value)*10)/10)); err != nil {
		return errors.Wrap(err, "SetSensorValue")
	}

	return nil
}

func (o *Devices) SetFieldValue(table string, itemID int, field string, value interface{}) error {
	err := o.store.db.
		Table(table).
		Where("view_item_id = ?", itemID).
		Update(field, value).
		Error

	if err != nil {
		return errors.Wrap(err, "SetFieldValue")
	}

	return nil
}

func (o *Devices) SaveSensor(sensor *model.Sensor) error {
	if sensor == nil {
		return errors.Wrap(errors.New("sensor is nil"), "SaveSensor")
	}

	count := int64(0)
	if err := o.store.db.Model(sensor).Where("id = ?", sensor.ID).Count(&count).Error; err != nil {
		return errors.Wrap(err, "SaveSensor")
	}
	objectIsExists := count == 1

	if objectIsExists {
		if err := o.store.db.Updates(sensor).Error; err != nil {
			return errors.Wrap(err, "SaveSensor(update)")
		}
	} else {
		result := o.store.db.Create(&sensor)
		if result.Error != nil {
			return errors.Wrap(result.Error, "SaveItem(create)")
		}
	}

	return nil
}

func (o *Devices) DeleteSensor(itemID int) error {
	if err := o.store.db.Table("sensors").
		Where("view_item_id = ?", itemID).Delete(&model.Sensor{}).
		Error; err != nil {
		return errors.Wrap(err, "UpdateItem")
	}

	return nil
}
