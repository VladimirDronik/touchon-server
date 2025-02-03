package sqlstore

import (
	"github.com/pkg/errors"
	"touchon-server/internal/model"
)

type Lights struct {
	store *Store
}

// GetLight Отдает данные источника света
func (o *Lights) GetLight(objectID int) (*model.ViewItem, error) {
	var r *model.ViewItem
	if err := o.store.db.First(&r, "id = ?", objectID).Error; err != nil {
		return nil, errors.Wrap(err, "GetLight")
	}

	var params *model.LightParams
	if err := o.store.db.Table("light_params").First(&params, "view_item_id = ?", objectID).Error; err != nil {
		return nil, errors.Wrap(err, "GetLight")
	}

	r.ParamsOutput = params

	return r, nil
}

// SetLightHSVColor указать HSV цвет
func (o *Lights) SetLightHSVColor(itemID int, hue int, saturation float32, brightness float32) error {
	params, err := o.GetParams(itemID)
	if err != nil {
		return errors.Wrap(err, "SetLightHSVColor")
	}

	params.Hue = hue
	params.Saturation = saturation
	params.Brightness = brightness
	//params.Cct = nil

	if err := o.store.db.Save(params).Error; err != nil {
		return errors.Wrap(err, "SetLightHSVColor")
	}

	return nil
}

// SetLightCCTColor указать CCT цвет
func (o *Lights) SetLightCCTColor(itemID int, cct int) error {
	if err := o.SetFieldValue(itemID, "cct", cct); err != nil {
		return errors.Wrap(err, "SetLightHSVColor")
	}

	return nil
}

// SetBrightness указать яркость
func (o *Lights) SetBrightness(itemID int, brightness float32) error {
	if err := o.SetFieldValue(itemID, "brightness", brightness); err != nil {
		return errors.Wrap(err, "SetBrightness")
	}

	return nil
}

func (o *Lights) GetParams(itemID int) (*model.LightParams, error) {
	r := &model.LightParams{}

	if err := o.store.db.Table("light_params").First(r, "view_item_id = ?", itemID).Error; err != nil {
		return nil, errors.Wrap(err, "GetParams")
	}

	return r, nil
}

func (o *Lights) SetFieldValue(itemID int, field string, value interface{}) error {
	err := o.store.db.
		Table("light_params").
		Where("view_item_id = ?", itemID).
		Update(field, value).
		Error

	if err != nil {
		return errors.Wrap(err, "SetFieldValue")
	}

	return nil
}
