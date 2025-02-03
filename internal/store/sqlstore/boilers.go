package sqlstore

import (
	"sort"

	"github.com/pkg/errors"
	"translator/internal/model"
)

type Boilers struct {
	store *Store
}

// GetBoiler Отдает данные для котла
func (o *Boilers) GetBoiler(boilerID int) (*model.Boiler, error) {
	r := &model.Boiler{}

	if err := o.store.db.Table("boilers").First(r, "id = ?", boilerID).Error; err != nil {
		return nil, errors.Wrap(err, "GetBoiler")
	}

	if err := o.store.db.Table("boiler_presets").Select("temp_out, temp_coolant").Find(&r.Presets, "boiler_id = ?", boilerID).Error; err != nil {
		return nil, errors.Wrap(err, "GetBoiler")
	}

	if err := o.store.db.Table("boiler_properties").Find(&r.Properties, "boiler_id = ?", boilerID).Error; err != nil {
		return nil, errors.Wrap(err, "GetBoiler")
	}

	return r, nil
}

func (o *Boilers) SetFieldValue(boilerID int, field string, value interface{}) error {
	err := o.store.db.
		Model(&model.Boiler{}).
		Where("id = ?", boilerID).
		Update(field, value).
		Error

	if err != nil {
		return errors.Wrap(err, "SetFieldValue")
	}

	return nil
}

// SetOutlineStatus включение или выключение контуров котла
func (o *Boilers) SetOutlineStatus(boilerID int, outline string, status string) error {
	if status != "on" && status != "off" {
		return errors.Wrap(errors.New("wrong status. Can be `on` or `off`"), "SetOutlineStatus")
	}

	if outline != "heating" && outline != "water" {
		return errors.Wrap(errors.New("wrong outline. Can be `heating` or `water`"), "SetOutlineStatus")
	}

	if err := o.SetFieldValue(boilerID, outline+"_status", status); err != nil {
		return errors.Wrap(err, "SetOutlineStatus")
	}

	return nil
}

// SetHeatingMode установка режима обогрева котла
func (o *Boilers) SetHeatingMode(boilerID int, mode string) error {
	if mode != "auto" && mode != "manual" {
		return errors.Wrap(errors.New("wrong mode. Can be `auto` or `manual`"), "SetHeatingMode")
	}

	if err := o.SetFieldValue(boilerID, "heating_mode", mode); err != nil {
		return errors.Wrap(err, "SetOutlineStatus")
	}

	return nil
}

// SetHeatingTemperature установке температуры для ручного обогрева котла
func (o *Boilers) SetHeatingTemperature(boilerID int, value float32) error {
	if err := o.SetFieldValue(boilerID, "heating_optimal_temp", value); err != nil {
		return errors.Wrap(err, "SetOutlineStatus")
	}

	return nil
}

// UpdateBoilerPresets обновляет все предустановки для указанного котла.
// Удаляет все существующие предустановки и вставляет новые.
func (o *Boilers) UpdateBoilerPresets(boilerID int, presets []*model.BoilerPreset) (e error) {
	sort.Slice(presets, func(i, j int) bool {
		return presets[i].TempOut < presets[j].TempOut
	})

	// Начинаем транзакцию
	tx := o.store.db.Begin()
	defer func() {
		if e != nil {
			tx.Rollback()
		}
	}()

	// Удаляем все существующие пресеты для данного котла
	if err := tx.Where("boiler_id = ?", boilerID).Delete(&model.BoilerPreset{}).Error; err != nil {
		return errors.Wrap(err, "UpdateBoilerPresets")
	}

	// Вставляем новые пресеты
	for _, preset := range presets {
		preset.BoilerID = boilerID

		if err := tx.Create(&preset).Error; err != nil {
			return errors.Wrap(err, "UpdateBoilerPresets")
		}
	}

	// Фиксируем транзакцию
	if err := tx.Commit().Error; err != nil {
		return errors.Wrap(err, "UpdateBoilerPresets")
	}

	return nil
}
