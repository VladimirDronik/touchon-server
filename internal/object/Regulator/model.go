package regulator

import (
	"time"

	"github.com/pkg/errors"
	"touchon-server/internal/g"
	"touchon-server/internal/model"
	"touchon-server/internal/object/Sensor"
	"touchon-server/internal/objects"
	"touchon-server/internal/store/memstore"
	"touchon-server/lib/events/object/regulator"
	"touchon-server/lib/helpers"
	"touchon-server/lib/interfaces"
	"touchon-server/lib/models"
)

type Type = string // Тип регулятора

const (
	TypeSimple  Type = "simple"  // Простой регулятор, управление одним выходом
	TypeComplex Type = "complex" // Управление двумя выходами
	TypePID     Type = "pid"     // Управление одним выходом с учетом законов регулирования
)

func init() {
	_ = objects.Register(MakeModel)
}

func MakeModel(withChildren bool) (objects.Object, error) {
	props := []*objects.Prop{
		{
			Code:        "type",
			Name:        "Тип регулятора",
			Description: "",
			Item: &models.Item{
				Type: models.DataTypeEnum,
				Values: map[string]string{
					TypeSimple:  "Simple",
					TypeComplex: "Complex",
					TypePID:     "PID",
				},
				DefaultValue: TypeSimple,
			},
			Required: objects.True(),
			Editable: objects.True(),
			Visible:  objects.True(),
		},
		{
			Code:        "fallback_sensor_value_id",
			Name:        "ID резервного объекта",
			Description: "ID объекта, который будет использоваться при недоступности основного датчика",
			Item: &models.Item{
				Type:         models.DataTypeInt,
				DefaultValue: 0,
			},
			Required: objects.False(),
			Editable: objects.True(),
			Visible:  objects.True(),
		},
		{
			Code:        "sensor_value_ttl",
			Name:        "Время жизни показания датчика (10s, 1m etc)",
			Description: "Время, в течение которого данные датчика будут считаться актуальными",
			Item: &models.Item{
				Type:         models.DataTypeString,
				DefaultValue: "120s",
			},
			Required:   objects.True(),
			Editable:   objects.True(),
			Visible:    objects.True(),
			CheckValue: objects.AboveOrEqual1(),
		},
		{
			Code:        "min_sp",
			Name:        "Минимальное значение уставки",
			Description: "",
			Item: &models.Item{
				Type:         models.DataTypeFloat,
				DefaultValue: 0,
			},
			Required:   objects.True(),
			Editable:   objects.True(),
			Visible:    objects.True(),
			CheckValue: objects.BelowOrEqual("target_sp"),
		},
		{
			Code:        "target_sp",
			Name:        "Заданное значение уставки",
			Description: "",
			Item: &models.Item{
				Type:         models.DataTypeFloat,
				DefaultValue: 0,
			},
			Required:   objects.True(),
			Editable:   objects.True(),
			Visible:    objects.True(),
			CheckValue: objects.BelowOrEqual("max_sp"),
		},
		{
			Code:        "max_sp",
			Name:        "Максимальное значение уставки",
			Description: "",
			Item: &models.Item{
				Type:         models.DataTypeFloat,
				DefaultValue: 0,
			},
			Required: objects.True(),
			Editable: objects.True(),
			Visible:  objects.True(),
		},
		{
			Code:        "below_tolerance",
			Name:        "Нижнее значение гистерезиса",
			Description: "",
			Item: &models.Item{
				Type:         models.DataTypeFloat,
				DefaultValue: 0,
			},
			Required:   objects.True(),
			Editable:   objects.True(),
			Visible:    objects.True(),
			CheckValue: objects.BelowOrEqual("above_tolerance"),
		}, {
			Code:        "above_tolerance",
			Name:        "Верхнее значение гистерезиса",
			Description: "",
			Item: &models.Item{
				Type:         models.DataTypeFloat,
				DefaultValue: 0,
			},
			Required: objects.True(),
			Editable: objects.True(),
			Visible:  objects.True(),
		},
		{
			Code:        "complex_tolerance",
			Name:        "Используется для определения точки перехода TargetSensor через TargetSP",
			Description: "",
			Item: &models.Item{
				Type:         models.DataTypeFloat,
				DefaultValue: 0,
			},
			Required: objects.NewCondition().PropValueEq("type", TypeComplex),
			Editable: objects.True(),
			Visible:  objects.True(),
		},
		{
			Code:        "update_interval",
			Name:        "Интервал (10s, 1m etc)",
			Description: "Интервал опроса датчиков",
			Item: &models.Item{
				Type:         models.DataTypeString,
				DefaultValue: "30s",
			},
			Required: objects.True(),
			Editable: objects.True(),
			Visible:  objects.True(),
		},
	}

	onBelow, err := regulator.NewOnBelow(0, 0)
	if err != nil {
		return nil, errors.Wrap(err, "Regulator.MakeModel")
	}

	onAbove, err := regulator.NewOnAbove(0, 0)
	if err != nil {
		return nil, errors.Wrap(err, "Regulator.MakeModel")
	}

	onStale, err := regulator.NewOnStale(0)
	if err != nil {
		return nil, errors.Wrap(err, "Regulator.MakeModel")
	}

	onComplexBelow1, err := regulator.NewOnComplexBelow1(0, 0)
	if err != nil {
		return nil, errors.Wrap(err, "Regulator.MakeModel")
	}

	onComplexBelow2, err := regulator.NewOnComplexBelow2(0, 0)
	if err != nil {
		return nil, errors.Wrap(err, "Regulator.MakeModel")
	}

	onComplexAbove1, err := regulator.NewOnComplexAbove1(0, 0)
	if err != nil {
		return nil, errors.Wrap(err, "Regulator.MakeModel")
	}

	onComplexAbove2, err := regulator.NewOnComplexAbove2(0, 0)
	if err != nil {
		return nil, errors.Wrap(err, "Regulator.MakeModel")
	}

	impl, err := objects.NewObjectModelImpl(
		model.CategoryRegulator,
		"regulator",
		0,
		"Регулятор",
		props,
		nil,
		[]interfaces.Event{onBelow, onAbove, onStale, onComplexBelow1, onComplexBelow2, onComplexAbove1, onComplexAbove2},
		nil,
		[]string{"regulator"},
	)
	if err != nil {
		return nil, errors.Wrap(err, "Regulator.MakeModel")
	}

	o := &RegulatorModel{ObjectModelImpl: impl}

	return o, nil
}

type RegulatorModel struct {
	*objects.ObjectModelImpl
}

func (o *RegulatorModel) check() {
	defer func() {
		if o.GetTimer() != nil {
			o.GetTimer().Reset()
		}
	}()

	g.Logger.Debugf("Regulator (%d) check", o.GetID())

	regType, err := o.GetProps().GetStringValue("type")
	if err != nil {
		g.Logger.Error(errors.Wrap(err, "RegulatorModel.check"))
		return
	}

	fallbackSensorValueID, err := o.GetProps().GetIntValue("fallback_sensor_value_id")
	if err != nil {
		g.Logger.Error(errors.Wrap(err, "RegulatorModel.check"))
		return
	}

	sensorValueTTLs, err := o.GetProps().GetStringValue("sensor_value_ttl")
	if err != nil {
		g.Logger.Error(errors.Wrap(err, "RegulatorModel.check"))
		return
	}

	sensorValueTTL, err := time.ParseDuration(sensorValueTTLs)
	if err != nil {
		g.Logger.Error(errors.Wrap(err, "RegulatorModel.check"))
		return
	}

	//minSP, err := o.GetProps().GetFloatValue("min_sp")
	//if err != nil {
	//g.Logger.Error(errors.Wrap(err, "RegulatorModel.check"))
	//return
	//}
	//
	//maxSP, err := o.GetProps().GetFloatValue("max_sp")
	//if err != nil {
	//g.Logger.Error(errors.Wrap(err, "RegulatorModel.check"))
	//return
	//}

	targetSP, err := o.GetProps().GetFloatValue("target_sp")
	if err != nil {
		g.Logger.Error(errors.Wrap(err, "RegulatorModel.check"))
		return
	}

	aboveTolerance, err := o.GetProps().GetFloatValue("above_tolerance")
	if err != nil {
		g.Logger.Error(errors.Wrap(err, "RegulatorModel.check"))
		return
	}

	belowTolerance, err := o.GetProps().GetFloatValue("below_tolerance")
	if err != nil {
		g.Logger.Error(errors.Wrap(err, "RegulatorModel.check"))
		return
	}

	complexTolerance, err := o.GetProps().GetFloatValue("complex_tolerance")
	if err != nil {
		g.Logger.Error(errors.Wrap(err, "RegulatorModel.check"))
		return
	}

	parentID := o.GetParentID()
	if parentID == nil {
		g.Logger.Error(errors.Wrap(errors.Errorf("parent_id of %d is nil", o.GetID()), "RegulatorModel.check"))
		return
	}

	// Запрос данных TargetSensor
	sensorValue, valueUpdatedAt, err := o.requestSensorValue(*parentID)
	if err != nil {
		if fallbackSensorValueID < 1 {
			g.Logger.Error(errors.Wrap(err, "RegulatorModel.check"))
			return
		}

		sensorValue, valueUpdatedAt, err = o.requestSensorValue(fallbackSensorValueID)
		if err != nil {
			g.Logger.Error(errors.Wrap(err, "RegulatorModel.check"))
			return
		}
	}

	// Если значение датчика(ов) давно не обновлялось, генерируем событие и уходим
	if time.Now().Sub(*valueUpdatedAt) > sensorValueTTL {
		msg, err := regulator.NewOnStale(*parentID)
		if err != nil {
			g.Logger.Error(errors.Wrap(err, "RegulatorModel.check"))
			return
		}

		if err := g.Msgs.Send(msg); err != nil {
			g.Logger.Error(errors.Wrap(err, "RegulatorModel.check"))
		}

		return
	}

	var msg interfaces.Message

	switch regType {
	case TypeSimple:
		switch {
		case sensorValue < (targetSP - belowTolerance):
			msg, err = regulator.NewOnBelow(*parentID, sensorValue)
		case sensorValue > (targetSP + aboveTolerance):
			msg, err = regulator.NewOnAbove(*parentID, sensorValue)
		}

	case TypeComplex:
		switch {
		case sensorValue < (targetSP - complexTolerance - belowTolerance):
			msg, err = regulator.NewOnComplexBelow1(*parentID, sensorValue)
		case sensorValue < (targetSP - complexTolerance + aboveTolerance):
			msg, err = regulator.NewOnComplexAbove1(*parentID, sensorValue)
		case sensorValue > (targetSP + complexTolerance - belowTolerance):
			msg, err = regulator.NewOnComplexBelow2(*parentID, sensorValue)
		case sensorValue > (targetSP + complexTolerance + aboveTolerance):
			msg, err = regulator.NewOnComplexAbove2(*parentID, sensorValue)
		}

	case TypePID:
		err = errors.New("regulator(TypePID) logic not implemented")

	default:
		g.Logger.Error(errors.Wrap(errors.Errorf("unexpected regulator type %q", regType), "RegulatorModel.check"))
		return
	}

	if err != nil {
		g.Logger.Error(errors.Wrap(err, "RegulatorModel.check"))
		return
	}

	if err := g.Msgs.Send(msg); err != nil {
		g.Logger.Error(errors.Wrap(err, "RegulatorModel.check"))
	}
}

func (o *RegulatorModel) requestSensorValue(sensorValueID int) (float32, *time.Time, error) {
	sensorValueObj, err := memstore.I.GetObject(sensorValueID)
	if err != nil {
		return 0, nil, errors.Wrap(err, "requestSensorValue")
	}

	// Проверяем, что датчик включен

	if !sensorValueObj.GetEnabled() {
		return 0, nil, errors.Wrap(errors.Errorf("sensor value %d is disabled", sensorValueID), "requestSensorValue")
	}

	sensorID := sensorValueObj.GetParentID()
	if sensorID == nil {
		return 0, nil, errors.Wrap(errors.Errorf("sensor value %d parent_id is nil", sensorValueID), "requestSensorValue")
	}

	sensorObj, err := memstore.I.GetObject(*sensorID)
	if err != nil {
		return 0, nil, errors.Wrap(err, "requestSensorValue")
	}

	if !sensorObj.GetEnabled() {
		return 0, nil, errors.Wrap(errors.Errorf("sensor %d is disabled", *sensorID), "requestSensorValue")
	}

	valueUpdatedAtS, err := sensorValueObj.GetProps().GetStringValue("value_updated_at")
	if err != nil {
		return 0, nil, errors.Wrap(err, "requestSensorValue")
	}

	valueUpdatedAt, err := time.Parse(Sensor.ValueUpdateAtFormat, valueUpdatedAtS)
	if err != nil {
		return 0, nil, errors.Wrap(err, "requestSensorValue")
	}

	value, err := sensorValueObj.GetProps().GetFloatValue("value")
	if err != nil {
		return 0, nil, errors.Wrap(err, "requestSensorValue")
	}

	return helpers.Round(value), &valueUpdatedAt, nil
}

func (o *RegulatorModel) Start() error {
	if err := o.ObjectModelImpl.Start(); err != nil {
		return errors.Wrap(err, "RegulatorModel.Start")
	}

	updateIntervalS, err := o.GetProps().GetStringValue("update_interval")
	if err != nil {
		return errors.Wrap(err, "RegulatorModel.Start")
	}

	updateInterval, err := time.ParseDuration(updateIntervalS)
	if err != nil {
		return errors.Wrap(err, "RegulatorModel.Start")
	}

	o.SetTimer(updateInterval, o.check)
	o.GetTimer().Start()

	return nil
}

func (o *RegulatorModel) Shutdown() error {
	if err := o.ObjectModelImpl.Shutdown(); err != nil {
		return errors.Wrap(err, "RegulatorModel.Shutdown")
	}

	if o.GetTimer() != nil {
		o.GetTimer().Stop()
	}

	return nil
}
