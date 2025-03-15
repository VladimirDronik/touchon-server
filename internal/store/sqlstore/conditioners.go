package sqlstore

import (
	"touchon-server/internal/g"
	"touchon-server/internal/objects"

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
	//params, err := o.GetParams(itemID)
	//if err != nil {
	//	return errors.Wrap(err, "SetConditionerOperatingMode")
	//}

	//if !slices.Contains(params.OperatingModes, mode) {
	//	return errors.Wrap(errors.Errorf("unknown mode %q", mode), "SetConditionerOperatingMode")
	//}

	if err := o.SetFieldValue(itemID, "operating_mode", mode); err != nil {
		return errors.Wrap(err, "SetConditionerOperatingMode")
	}

	return nil
}

// SetConditionerDirection указать направление ламелей для кондиционера
func (o *Conditioners) SetConditionerDirection(itemID int, plane string, direction string) error {
	//params, err := o.GetParams(itemID)
	//if err != nil {
	//	return errors.Wrap(err, "SetConditionerDirection")
	//}
	//
	//var directions []string
	//
	//switch plane {
	//case "vertical":
	//	directions = params.VerticalDirections
	//case "horizontal":
	//	directions = params.HorizontalDirections
	//default:
	//	return errors.Wrap(errors.Errorf("unknown plane %q", plane), "SetConditionerDirection")
	//}
	//
	//if !slices.Contains(directions, direction) {
	//	return errors.Wrap(errors.Errorf("unknown direction %q", direction), "SetConditionerDirection")
	//}
	//
	//if err := o.SetFieldValue(itemID, plane+"_direction", direction); err != nil {
	//	return errors.Wrap(err, "SetConditionerDirection")
	//}
	//
	return nil
}

// SetConditionerFanSpeed указать скорость вентилятора для кондиционера
func (o *Conditioners) SetConditionerFanSpeed(itemID int, speed string) error {
	//params, err := o.GetParams(itemID)
	//if err != nil {
	//	return errors.Wrap(err, "SetConditionerFanSpeed")
	//}
	//
	//if !slices.Contains(params.FanSpeeds, speed) {
	//	return errors.Wrap(errors.Errorf("unknown speed %q", speed), "SetConditionerFanSpeed")
	//}
	//
	//if err := o.SetFieldValue(itemID, "fan_speed", speed); err != nil {
	//	return errors.Wrap(err, "SetConditionerFanSpeed")
	//}
	//
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

func (o *Conditioners) GetParams(itemID int) (*model.Conditioner, error) {
	conditioner := &model.Conditioner{}
	var objectID int

	o.store.db.Debug().Table("conditioner").Select("object_id").
		Where("view_item_id = ?", itemID).Find(&objectID)

	conditioner.CondParams.ID = objectID
	conditioner.CondParams.ViewItemID = itemID

	condObj, err := objects.LoadObject(objectID, "", "", false)
	if err != nil {
		return nil, errors.Wrap(err, "GetParams")
	}

	opModes, err := condObj.GetProps().Get("operating_mode")
	if err != nil {
		g.Logger.Error(errors.Wrap(err, "GetParams: operating_mode"))
	}

	fanSpeeds, err := condObj.GetProps().Get("fan_speed")
	if err != nil {
		g.Logger.Error(errors.Wrap(err, "GetParams: fan_speed"))
	}

	horizontalDir, err := condObj.GetProps().Get("horizontal_slats_mode")
	if err != nil {
		g.Logger.Error(errors.Wrap(err, "GetParams: horizontal_slats_mode"))
	}

	verticalDir, err := condObj.GetProps().Get("vertical_slats_mode")
	if err != nil {
		return nil, errors.Wrap(err, "GetParams: vertical_slats_mode")
	}

	conditioner.OperatingModes = opModes.Values
	conditioner.FanSpeeds = fanSpeeds.Values
	conditioner.HorizontalDirections = horizontalDir.Values
	conditioner.VerticalDirections = verticalDir.Values

	conditioner.CondParams.FanSpeed, err = condObj.GetProps().GetIntValue("fan_speed")
	if err != nil {
		g.Logger.Error(errors.Wrap(err, "GetParams: fan_speed"))
	}

	conditioner.CondParams.OperatingMode, err = condObj.GetProps().GetIntValue("operating_mode")
	if err != nil {
		g.Logger.Error(errors.Wrap(err, "GetParams: operating_mode"))
	}

	conditioner.CondParams.HorizontalDirection, err = condObj.GetProps().GetIntValue("horizontal_slats_mode")
	if err != nil {
		g.Logger.Error(errors.Wrap(err, "GetParams: horizontal_slats_mode"))
	}

	conditioner.CondParams.VerticalDirection, err = condObj.GetProps().GetIntValue("vertical_slats_mode")
	if err != nil {
		g.Logger.Error(errors.Wrap(err, "GetParams: vertical_slats_mode"))
	}

	conditioner.CondParams.DisplayBacklight, err = condObj.GetProps().GetBoolValue("display_backlight")
	if err != nil {
		g.Logger.Error(errors.Wrap(err, "GetParams: display_backlight"))
	}

	conditioner.CondParams.EcoMode, err = condObj.GetProps().GetBoolValue("eco_mode")
	if err != nil {
		g.Logger.Error(errors.Wrap(err, "GetParams: eco_mode"))
	}

	conditioner.CondParams.OutsideTemp, err = condObj.GetProps().GetFloatValue("external_temperature")
	if err != nil {
		g.Logger.Error(errors.Wrap(err, "GetParams: external_temperature"))
	}

	conditioner.CondParams.InsideTemp, err = condObj.GetProps().GetFloatValue("internal_temperature")
	if err != nil {
		g.Logger.Error(errors.Wrap(err, "GetParams: internal_temperature"))
	}

	conditioner.CondParams.Ionisation, err = condObj.GetProps().GetBoolValue("ionisation")
	if err != nil {
		g.Logger.Error(errors.Wrap(err, "GetParams: ionisation"))
	}

	conditioner.CondParams.PowerStatus, err = condObj.GetProps().GetBoolValue("power_status")
	if err != nil {
		g.Logger.Error(errors.Wrap(err, "GetParams: power_status"))
	}

	conditioner.CondParams.SelfCleaning, err = condObj.GetProps().GetBoolValue("self_cleaning")
	if err != nil {
		g.Logger.Error(errors.Wrap(err, "GetParams: self_cleaning"))
	}

	conditioner.CondParams.SilentMode, err = condObj.GetProps().GetBoolValue("silent_mode")
	if err != nil {
		g.Logger.Error(errors.Wrap(err, "GetParams: silent_mode"))
	}

	conditioner.CondParams.SleepMode, err = condObj.GetProps().GetBoolValue("sleep_mode")
	if err != nil {
		g.Logger.Error(errors.Wrap(err, "GetParams: sleep_mode"))
	}

	conditioner.CondParams.Sound, err = condObj.GetProps().GetBoolValue("sounds")
	if err != nil {
		g.Logger.Error(errors.Wrap(err, "GetParams: sounds"))
	}

	conditioner.CondParams.TargetTemp, err = condObj.GetProps().GetFloatValue("target_temperature")
	if err != nil {
		g.Logger.Error(errors.Wrap(err, "GetParams: target_temperature"))
	}

	conditioner.CondParams.TurboMode, err = condObj.GetProps().GetBoolValue("turbo_mode")
	if err != nil {
		g.Logger.Error(errors.Wrap(err, "GetParams: turbo_mode"))
	}

	return conditioner, nil
}
