package regulator

import (
	"strconv"
	"time"

	"github.com/pkg/errors"
	"touchon-server/internal/context"
	"touchon-server/internal/model"
	"touchon-server/internal/objects"
	"touchon-server/internal/store"
	"touchon-server/lib/events/object/regulator"
	"touchon-server/lib/helpers"
	"touchon-server/lib/models"
	mqttClient "touchon-server/lib/mqtt/client"
	"touchon-server/lib/mqtt/messages"
)

type Type string // Тип регулятора

const (
	TypeSimple  Type = "simple"  // Простой регулятор, управление одним выходом
	TypeComplex Type = "complex" // Управление двумя выходами
	TypePID     Type = "pid"     // Управление одним выходом с учетом законов регулирования
)

func init() {
	_ = objects.Register(MakeModel)
}

func MakeModel() (objects.Object, error) {
	props := []*objects.Prop{
		{
			Code:        "enable",
			Name:        "Состояние регулятора",
			Description: "true - автоматический режим, false - регулятор выключен",
			Item: &models.Item{
				Type:         models.DataTypeBool,
				DefaultValue: false,
			},
			Required: objects.NewRequired(true),
			Editable: objects.NewCondition(),
			Visible:  objects.NewCondition(),
		},
		{
			Code:        "type",
			Name:        "Тип регулятора",
			Description: "",
			Item: &models.Item{
				Type: models.DataTypeEnum,
				Values: map[string]string{
					string(TypeSimple):  "Simple",
					string(TypeComplex): "Complex",
					string(TypePID):     "PID",
				},
				DefaultValue: string(TypeSimple),
			},
			Required: objects.NewRequired(true),
			Editable: objects.NewCondition(),
			Visible:  objects.NewCondition(),
		},
		{
			Code:        "fallback_sensor_value_id",
			Name:        "ID резервного объекта",
			Description: "ID объекта, который будет использоваться при недоступности основного датчика",
			Item: &models.Item{
				Type:         models.DataTypeInt,
				DefaultValue: 0,
			},
			Required: objects.NewRequired(false),
			Editable: objects.NewCondition(),
			Visible:  objects.NewCondition(),
		},
		{
			Code:        "sensor_value_ttl",
			Name:        "Время жизни показания датчика",
			Description: "Время в секундах, в течение которого данные датчика будут считаться актуальными",
			Item: &models.Item{
				Type:         models.DataTypeInt,
				DefaultValue: 30,
			},
			Required:   objects.NewRequired(true),
			Editable:   objects.NewCondition(),
			Visible:    objects.NewCondition(),
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
			Required:   objects.NewRequired(true),
			Editable:   objects.NewCondition(),
			Visible:    objects.NewCondition(),
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
			Required:   objects.NewRequired(true),
			Editable:   objects.NewCondition(),
			Visible:    objects.NewCondition(),
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
			Required: objects.NewRequired(true),
			Editable: objects.NewCondition(),
			Visible:  objects.NewCondition(),
		},
		{
			Code:        "below_tolerance",
			Name:        "Нижнее значение гистерезиса",
			Description: "",
			Item: &models.Item{
				Type:         models.DataTypeFloat,
				DefaultValue: 0,
			},
			Required:   objects.NewRequired(true),
			Editable:   objects.NewCondition(),
			Visible:    objects.NewCondition(),
			CheckValue: objects.BelowOrEqual("above_tolerance"),
		}, {
			Code:        "above_tolerance",
			Name:        "Верхнее значение гистерезиса",
			Description: "",
			Item: &models.Item{
				Type:         models.DataTypeFloat,
				DefaultValue: 0,
			},
			Required: objects.NewRequired(true),
			Editable: objects.NewCondition(),
			Visible:  objects.NewCondition(),
		},
		{
			Code:        "complex_tolerance",
			Name:        "Используется для определения точки перехода TargetSensor через TargetSP",
			Description: "",
			Item: &models.Item{
				Type:         models.DataTypeFloat,
				DefaultValue: 0,
			},
			Required: objects.NewCondition().PropValueEq("type", string(TypeComplex)),
			Editable: objects.NewCondition(),
			Visible:  objects.NewCondition(),
		},
	}

	impl, err := objects.NewObjectModelImpl(
		model.CategoryRegulator,
		"regulator",
		false,
		"Регулятор",
		props,
		nil,
		[]string{
			"object.regulator.on_below",
			"object.regulator.on_above",
			"object.regulator.on_stale",
			"object.regulator.on_complex_below_1",
			"object.regulator.on_complex_above_1",
			"object.regulator.on_complex_below_2",
			"object.regulator.on_complex_above_2",
		},
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
	timer *helpers.Timer
}

func (o *RegulatorModel) sensorOnCheckHandler(messages.Message) ([]messages.Message, error) {
	// Получаем значения свойств

	enable, err := o.GetProps().GetBoolValue("enable")
	if err != nil {
		return nil, errors.Wrap(err, "RegulatorModel.sensorOnCheckHandler")
	}

	if !enable {
		return nil, nil
	}

	// Раз попали сюда, значит значение датчика обновилось
	// o.timer будет nil, если не выполняли вызов метода Start
	// Start вызывается в memstore, не вызывается в http эндпоинтах
	if o.timer != nil {
		o.timer.Reset()
	}

	regTypeS, err := o.GetProps().GetStringValue("type")
	if err != nil {
		return nil, errors.Wrap(err, "RegulatorModel.sensorOnCheckHandler")
	}
	regType := Type(regTypeS)

	fallbackSensorValueID, err := o.GetProps().GetIntValue("fallback_sensor_value_id")
	if err != nil {
		return nil, errors.Wrap(err, "RegulatorModel.sensorOnCheckHandler")
	}

	//minSP, err := o.GetProps().GetFloatValue("min_sp")
	//if err != nil {
	//	return nil, errors.Wrap(err, "RegulatorModel.sensorOnCheckHandler")
	//}
	//
	//maxSP, err := o.GetProps().GetFloatValue("max_sp")
	//if err != nil {
	//	return nil, errors.Wrap(err, "RegulatorModel.sensorOnCheckHandler")
	//}

	targetSP, err := o.GetProps().GetFloatValue("target_sp")
	if err != nil {
		return nil, errors.Wrap(err, "RegulatorModel.sensorOnCheckHandler")
	}

	aboveTolerance, err := o.GetProps().GetFloatValue("above_tolerance")
	if err != nil {
		return nil, errors.Wrap(err, "RegulatorModel.sensorOnCheckHandler")
	}

	belowTolerance, err := o.GetProps().GetFloatValue("below_tolerance")
	if err != nil {
		return nil, errors.Wrap(err, "RegulatorModel.sensorOnCheckHandler")
	}

	complexTolerance, err := o.GetProps().GetFloatValue("complex_tolerance")
	if err != nil {
		return nil, errors.Wrap(err, "RegulatorModel.sensorOnCheckHandler")
	}

	parentID := o.GetParentID()
	if parentID == nil {
		return nil, errors.Wrap(errors.Errorf("parent_id of %d is nil", o.GetID()), "RegulatorModel.sensorOnCheckHandler")
	}

	// Запрос данных TargetSensor
	sensorValue, err := o.requestSensorValue(*parentID)
	if err != nil {
		if fallbackSensorValueID < 1 {
			return nil, errors.Wrap(err, "RegulatorModel.sensorOnCheckHandler")
		}

		sensorValue, err = o.requestSensorValue(fallbackSensorValueID)
		if err != nil {
			return nil, errors.Wrap(err, "RegulatorModel.sensorOnCheckHandler")
		}
	}

	var msg messages.Message

	switch regType {
	case TypeSimple:
		switch {
		case sensorValue < (targetSP - belowTolerance): // TODO minSP ???
			msg, err = regulator.NewOnBelowMessage("object_manager/object/event", *parentID, sensorValue)
		case sensorValue > (targetSP + aboveTolerance): // TODO maxSP ???
			msg, err = regulator.NewOnAboveMessage("object_manager/object/event", *parentID, sensorValue)
		}

	case TypeComplex:
		switch {
		case sensorValue < (targetSP - complexTolerance - belowTolerance):
			msg, err = regulator.NewOnComplexBelow1Message("object_manager/object/event", *parentID, sensorValue)
		case sensorValue < (targetSP - complexTolerance + aboveTolerance):
			msg, err = regulator.NewOnComplexAbove1Message("object_manager/object/event", *parentID, sensorValue)
		case sensorValue > (targetSP + complexTolerance - belowTolerance):
			msg, err = regulator.NewOnComplexBelow2Message("object_manager/object/event", *parentID, sensorValue)
		case sensorValue > (targetSP + complexTolerance + aboveTolerance):
			msg, err = regulator.NewOnComplexAbove2Message("object_manager/object/event", *parentID, sensorValue)
		}

	case TypePID:
		err = errors.New("regulator(TypePID) logic not implemented")

	default:
		return nil, errors.Wrap(errors.Errorf("unexpected regulator type %q", regType), "RegulatorModel.sensorOnCheckHandler")
	}

	if err != nil {
		return nil, errors.Wrap(err, "RegulatorModel.sensorOnCheckHandler")
	}

	return []messages.Message{msg}, nil
}

func (o *RegulatorModel) requestSensorValue(sensorValueID int) (float32, error) {
	s, err := store.I.ObjectRepository().GetProp(sensorValueID, "value")
	if err != nil {
		return 0, errors.Wrap(err, "requestSensorValue")
	}

	value, err := strconv.ParseFloat(s, 32)
	if err != nil {
		return 0, errors.Wrap(err, "requestSensorValue")
	}

	return helpers.Round(float32(value)), nil
}

func (o *RegulatorModel) Start() error {
	if err := o.ObjectModelImpl.Start(); err != nil {
		return errors.Wrap(err, "RegulatorModel.Start")
	}

	enable, err := o.GetProps().GetBoolValue("enable")
	if err != nil {
		return errors.Wrap(err, "RegulatorModel.Start")
	}

	if !enable {
		return nil
	}

	parentID := o.GetParentID()
	if parentID == nil {
		return errors.Wrap(errors.Errorf("parent_id of %d is nil", o.GetID()), "RegulatorModel.Start")
	}

	sensorValue, err := store.I.ObjectRepository().GetObject(*parentID)
	if err != nil {
		return errors.Wrap(err, "RegulatorModel.Start")
	}

	err = o.Subscribe(
		"",
		"",
		messages.MessageTypeEvent,
		"object.sensor.on_check",
		messages.TargetTypeObject,
		sensorValue.ParentID,
		o.sensorOnCheckHandler,
	)
	if err != nil {
		return errors.Wrap(err, "RegulatorModel.Start")
	}

	sensorValueTTL, err := o.GetProps().GetIntValue("sensor_value_ttl")
	if err != nil {
		return errors.Wrap(err, "RegulatorModel.Start")
	}

	o.timer = helpers.NewTimer(time.Duration(sensorValueTTL)*time.Second, o.timerHandler)
	o.timer.Start()

	context.Logger.Debugf("Regulator(%d) started", o.GetID())

	return nil
}

func (o *RegulatorModel) timerHandler() {
	enable, err := o.GetProps().GetBoolValue("enable")
	if err != nil {
		context.Logger.Error(errors.Wrap(err, "RegulatorModel.timerHandler"))
		return
	}

	if !enable {
		return
	}

	context.Logger.Debugf("Regulator(%d): timerHandler(%s)", o.GetID(), o.timer.GetDuration())

	parentID := o.GetParentID()
	if parentID == nil {
		context.Logger.Error(errors.Wrap(errors.Errorf("parent_id of %d is nil", o.GetID()), "RegulatorModel.timerHandler"))
		return

	}

	msg, err := regulator.NewOnStaleMessage("object_manager/object/event", *parentID)
	if err != nil {
		context.Logger.Error(errors.Wrap(err, "RegulatorModel.timerHandler"))
		return
	}

	if err := mqttClient.I.Send(msg); err != nil {
		context.Logger.Error(errors.Wrap(err, "RegulatorModel.timerHandler"))
		return
	}

	o.timer.Reset()
}

func (o *RegulatorModel) Shutdown() error {
	if err := o.ObjectModelImpl.Shutdown(); err != nil {
		return errors.Wrap(err, "RegulatorModel.Shutdown")
	}

	enable, err := o.GetProps().GetBoolValue("enable")
	if err != nil {
		return errors.Wrap(err, "RegulatorModel.Shutdown")
	}

	if !enable {
		return nil
	}

	o.timer.Stop()

	context.Logger.Debugf("Regulator(%d) stopped", o.GetID())

	return nil
}
