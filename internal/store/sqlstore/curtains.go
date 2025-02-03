package sqlstore

import (
	"github.com/pkg/errors"
	"translator/internal/model"
)

type Curtains struct {
	store *Store
}

// GetCurtain Отдает данные для штор
func (o *Curtains) GetCurtain(itemID int) (*model.ViewItem, error) {
	r := &model.ViewItem{}
	if err := o.store.db.First(r, "id = ?", itemID).Error; err != nil {
		return nil, errors.Wrap(err, "GetCurtain")
	}

	params, err := o.GetParams(itemID)
	if err != nil {
		return nil, errors.Wrap(err, "GetCurtain")
	}

	r.ParamsOutput = params

	return r, nil
}

// SetCurtainOpenPercent указать температуру для кондиционера
func (o *Curtains) SetCurtainOpenPercent(itemID int, value float32) error {
	params, err := o.GetParams(itemID)
	if err != nil {
		return errors.Wrap(err, "SetCurtainOpenPercent")
	}

	if params.ControlType != "rs485" {
		return errors.Wrap(errors.Errorf("curtain with ID %d must be a rs485 type to have open percent", itemID), "SetCurtainOpenPercent")
	}

	if err := o.SetFieldValue(itemID, "open_percent", value); err != nil {
		return errors.Wrap(err, "SetCurtainOpenPercent")
	}

	return nil
}

func (o *Curtains) SetFieldValue(itemID int, field string, value interface{}) error {
	err := o.store.db.
		Model(&model.CurtainParams{}).
		Where("view_item_id = ?", itemID).
		Update(field, value).
		Error

	if err != nil {
		return errors.Wrap(err, "SetFieldValue")
	}

	return nil
}

func (o *Curtains) GetParams(itemID int) (*model.CurtainParams, error) {
	r := &model.CurtainParams{}

	if err := o.store.db.Table("curtain_params").First(r, "view_item_id = ?", itemID).Error; err != nil {
		return nil, errors.Wrap(err, "GetParams")
	}

	return r, nil
}
