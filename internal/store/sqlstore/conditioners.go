package sqlstore

import (
	"encoding/json"
	"slices"

	"github.com/pkg/errors"
	"touchon-server/internal/model"
)

type Conditioners struct {
	store *Store
}

// GetConditioner Отдает данные кондиционера
func (o *Conditioners) GetConditioner(conditionerID int) (*model.ViewItem, error) {
	r := &model.ViewItem{}
	if err := o.store.db.First(r, "id = ?", conditionerID).Error; err != nil {
		return nil, errors.Wrap(err, "GetConditioner")
	}

	params, err := o.GetParams(conditionerID)
	if err != nil {
		return nil, errors.Wrap(err, "GetConditioner")
	}

	r.ParamsOutput = params

	return r, nil
}

// SetConditionerTemperature указать температуру для кондиционера
func (o *Conditioners) SetConditionerTemperature(itemID int, value float32) error {
	if err := o.SetFieldValue(itemID, "optimal_temp", value); err != nil {
		return errors.Wrap(err, "SetConditionerTemperature")
	}

	return nil
}

// SetConditionerMode указать режим для кондиционера
func (o *Conditioners) SetConditionerMode(itemID int, mode string, value bool) error {
	switch mode {
	case "eco_mode", "silent_mode", "turbo_mode", "sleep_mode":
	default:
		return errors.Wrap(errors.Errorf("unknown mode %q", mode), "SetConditionerMode")
	}

	if err := o.SetFieldValue(itemID, mode, value); err != nil {
		return errors.Wrap(err, "SetConditionerMode")
	}

	return nil
}

// SetConditionerOperatingMode указать режим работы для кондиционера
func (o *Conditioners) SetConditionerOperatingMode(itemID int, mode string) error {
	params, err := o.GetParams(itemID)
	if err != nil {
		return errors.Wrap(err, "SetConditionerOperatingMode")
	}

	if !slices.Contains(params.OperatingModes, mode) {
		return errors.Wrap(errors.Errorf("unknown mode %q", mode), "SetConditionerOperatingMode")
	}

	if err := o.SetFieldValue(itemID, "operating_mode", mode); err != nil {
		return errors.Wrap(err, "SetConditionerOperatingMode")
	}

	return nil
}

// SetConditionerDirection указать направление ламелей для кондиционера
func (o *Conditioners) SetConditionerDirection(itemID int, plane string, direction string) error {
	params, err := o.GetParams(itemID)
	if err != nil {
		return errors.Wrap(err, "SetConditionerDirection")
	}

	var directions []string

	switch plane {
	case "vertical":
		directions = params.VerticalDirections
	case "horizontal":
		directions = params.HorizontalDirections
	default:
		return errors.Wrap(errors.Errorf("unknown plane %q", plane), "SetConditionerDirection")
	}

	if !slices.Contains(directions, direction) {
		return errors.Wrap(errors.Errorf("unknown direction %q", direction), "SetConditionerDirection")
	}

	if err := o.SetFieldValue(itemID, plane+"_direction", direction); err != nil {
		return errors.Wrap(err, "SetConditionerDirection")
	}

	return nil
}

// SetConditionerFanSpeed указать скорость вентилятора для кондиционера
func (o *Conditioners) SetConditionerFanSpeed(itemID int, speed string) error {
	params, err := o.GetParams(itemID)
	if err != nil {
		return errors.Wrap(err, "SetConditionerFanSpeed")
	}

	if !slices.Contains(params.FanSpeeds, speed) {
		return errors.Wrap(errors.Errorf("unknown speed %q", speed), "SetConditionerFanSpeed")
	}

	if err := o.SetFieldValue(itemID, "fan_speed", speed); err != nil {
		return errors.Wrap(err, "SetConditionerFanSpeed")
	}

	return nil
}

// SetConditionerExtraMode указать доп режим для кондиционера
func (o *Conditioners) SetConditionerExtraMode(itemID int, mode string, value bool) error {
	switch mode {
	case "ionisation", "self_cleaning", "anti_mold", "sound", "on_duty_heating", "soft_top":
	default:
		return errors.Wrap(errors.Errorf("unknown extra mode %q", mode), "SetConditionerExtraMode")
	}

	if err := o.SetFieldValue(itemID, mode, value); err != nil {
		return errors.Wrap(err, "SetConditionerExtraMode")
	}

	return nil
}

func (o *Conditioners) SetFieldValue(conditionerID int, field string, value interface{}) error {
	err := o.store.db.
		Table("conditioner_params").
		Where("view_item_id = ?", conditionerID).
		Update(field, value).
		Error

	if err != nil {
		return errors.Wrap(err, "SetFieldValue")
	}

	return nil
}

func (o *Conditioners) GetParams(itemID int) (*model.ConditionerParams, error) {
	row := &model.StoreConditionerParams{}

	if err := o.store.db.Table("conditioner_params").First(row, "view_item_id = ?", itemID).Error; err != nil {
		return nil, errors.Wrap(err, "GetParams")
	}

	r := &model.ConditionerParams{CondParams: row.CondParams}

	if err := json.Unmarshal([]byte(row.OperatingModes), &r.OperatingModes); err != nil {
		return nil, errors.Wrap(err, "GetParams")
	}

	if err := json.Unmarshal([]byte(row.FanSpeeds), &r.FanSpeeds); err != nil {
		return nil, errors.Wrap(err, "GetParams")
	}

	if err := json.Unmarshal([]byte(row.VerticalDirections), &r.VerticalDirections); err != nil {
		return nil, errors.Wrap(err, "GetParams")
	}

	if err := json.Unmarshal([]byte(row.HorizontalDirections), &r.HorizontalDirections); err != nil {
		return nil, errors.Wrap(err, "GetParams")
	}

	return r, nil
}
