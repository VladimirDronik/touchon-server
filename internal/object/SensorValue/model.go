package SensorValue

import (
	"github.com/pkg/errors"
	"touchon-server/internal/model"
	regulator "touchon-server/internal/object/Regulator"
	"touchon-server/internal/objects"
	"touchon-server/internal/store"
	"touchon-server/lib/models"
)

var (
	ErrSensorAlarmValue = errors.New("alarm-sensor")
	ErrSensorErrorValue = errors.New("error-sensor")
)

func MakeModel() (objects.Object, error) {
	props := []*objects.Prop{
		{
			Code:        "value",
			Name:        "Значение датчика",
			Description: "",
			Item: &models.Item{
				Type:         models.DataTypeFloat,
				RoundFloat:   true,
				DefaultValue: 0,
			},
			Required: objects.False(),
			Editable: objects.False(),
			Visible:  objects.True(),
		},
		{
			Code:        "value_updated_at",
			Name:        "Дата обновления значения",
			Description: "",
			Item: &models.Item{
				Type: models.DataTypeString,
			},
			Required: objects.False(),
			Editable: objects.False(),
			Visible:  objects.True(),
		},
		{
			Code:        "write_graph",
			Name:        "Вести график",
			Description: "Записывать значения датчика в таблицу графиков",
			Item: &models.Item{
				Type:         models.DataTypeBool,
				DefaultValue: true,
			},
			Required: objects.True(),
			Editable: objects.True(),
			Visible:  objects.True(),
		},
		{
			Code:        "unit",
			Name:        "Единица измерения",
			Description: "",
			Item: &models.Item{
				Type:         models.DataTypeString,
				DefaultValue: "<ед. изм.>",
			},
			Required: objects.True(),
			Editable: objects.True(),
			Visible:  objects.True(),
		},
		{
			Code:        "min_error_value",
			Name:        "Минимальное значение аварии",
			Description: "",
			Item: &models.Item{
				Type:         models.DataTypeFloat,
				DefaultValue: 0,
			},
			Required:   objects.True(),
			Editable:   objects.True(),
			Visible:    objects.True(),
			CheckValue: objects.BelowOrEqual("min_threshold"),
		}, {
			Code:        "min_threshold",
			Name:        "Минимальное значение",
			Description: "",
			Item: &models.Item{
				Type:         models.DataTypeFloat,
				DefaultValue: 0,
			},
			Required:   objects.True(),
			Editable:   objects.True(),
			Visible:    objects.True(),
			CheckValue: objects.BelowOrEqual("max_threshold"),
		},
		{
			Code:        "max_threshold",
			Name:        "Максимальное значение",
			Description: "",
			Item: &models.Item{
				Type:         models.DataTypeFloat,
				DefaultValue: 0,
			},
			Required:   objects.True(),
			Editable:   objects.True(),
			Visible:    objects.True(),
			CheckValue: objects.BelowOrEqual("max_error_value"),
		},
		{
			Code:        "max_error_value",
			Name:        "Максимальное значение аварии",
			Description: "",
			Item: &models.Item{
				Type:         models.DataTypeFloat,
				DefaultValue: 0,
			},
			Required: objects.True(),
			Editable: objects.True(),
			Visible:  objects.True(),
		},
	}

	impl, err := objects.NewObjectModelImpl(
		model.CategorySensorValue,
		"",
		true,
		"",
		props,
		nil,
		nil,
		nil,
		[]string{"sensor_value"},
	)
	if err != nil {
		return nil, errors.Wrap(err, "SensorValue.MakeModel")
	}

	reg, err := regulator.MakeModel()
	if err != nil {
		return nil, errors.Wrap(err, "SensorValue.MakeModel")
	}

	regProps := map[string]interface{}{
		"type":             string(regulator.TypeSimple),
		"sensor_value_ttl": 30,
	}

	for k, v := range regProps {
		if err := reg.GetProps().Set(k, v); err != nil {
			return nil, errors.Wrap(err, "SensorValue.MakeModel")
		}
	}

	impl.GetChildren().Add(reg)

	return &SensorValueModel{ObjectModelImpl: impl}, nil
}

type SensorValueModel struct {
	*objects.ObjectModelImpl
}

// CheckValue проверяет значение датчика на пороговое значение и на значение, которое находится за пределами восприимчивости датчика
func (o *SensorValueModel) Validate() error {
	value, err := o.getFloatValue("value")
	if err != nil {
		return errors.Wrap(err, "Validate")
	}

	minThreshold, err := o.getFloatValue("min_threshold")
	if err != nil {
		return errors.Wrap(err, "Validate")
	}

	maxThreshold, err := o.getFloatValue("max_threshold")
	if err != nil {
		return errors.Wrap(err, "Validate")
	}

	minErrorValue, err := o.getFloatValue("min_error_value")
	if err != nil {
		return errors.Wrap(err, "Validate")
	}

	maxErrorValue, err := o.getFloatValue("max_error_value")
	if err != nil {
		return errors.Wrap(err, "Validate")
	}

	switch {
	// проверка на аварийные значения, если датчик выходит за аварийные значения то генерим событие
	case value < minErrorValue || value > maxErrorValue:
		return errors.Wrap(ErrSensorErrorValue, "Validate")

	// проверка на пороговые значения, и проверку на исправность датчика
	case value < minThreshold || value > maxThreshold:
		return errors.Wrap(ErrSensorAlarmValue, "Validate")
	}

	return nil
}

func (o *SensorValueModel) getFloatValue(code string) (float32, error) {
	p, err := o.GetProps().Get(code)
	if err != nil {
		return 0, errors.Wrap(err, "getFloatValue")
	}

	value, err := p.GetFloatValue()
	if err != nil {
		return 0, errors.Wrap(err, "getFloatValue")
	}

	return value, nil
}

func (o *SensorValueModel) DeleteChildren() error {
	for _, child := range o.GetChildren().GetAll() {
		if err := child.DeleteChildren(); err != nil {
			return errors.Wrap(err, "DeleteChildren")
		}

		// Удаляем только регулятор
		if !(child.GetCategory() == model.CategoryRegulator && child.GetType() == "regulator") {
			continue
		}

		if err := store.I.ObjectRepository().DelObject(child.GetID()); err != nil {
			return errors.Wrap(err, "DeleteChildren")
		}
	}

	return nil
}
