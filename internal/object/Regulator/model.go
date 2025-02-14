package regulator

import (
	"strconv"
	"time"

	"github.com/pkg/errors"
	"touchon-server/internal/g"
	"touchon-server/internal/model"
	"touchon-server/internal/objects"
	"touchon-server/internal/store"
	"touchon-server/lib/events/object/regulator"
	"touchon-server/lib/helpers"
	"touchon-server/lib/interfaces"
	"touchon-server/lib/models"
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
		false,
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
	timer *helpers.Timer
}

func (o *RegulatorModel) sensorOnCheckHandler(svc interfaces.MessageSender, _ interfaces.Message) {
	if !o.GetEnabled() {
		return
	}

	// Раз попали сюда, значит значение датчика обновилось
	// o.timer будет nil, если не выполняли вызов метода Start
	// Start вызывается в memstore, не вызывается в http эндпоинтах
	if o.timer != nil {
		o.timer.Reset()
	}

	regTypeS, err := o.GetProps().GetStringValue("type")
	if err != nil {
		g.Logger.Error(errors.Wrap(err, "RegulatorModel.sensorOnCheckHandler"))
		return
	}
	regType := Type(regTypeS)

	fallbackSensorValueID, err := o.GetProps().GetIntValue("fallback_sensor_value_id")
	if err != nil {
		g.Logger.Error(errors.Wrap(err, "RegulatorModel.sensorOnCheckHandler"))
		return
	}

	//minSP, err := o.GetProps().GetFloatValue("min_sp")
	//if err != nil {
	//context.Logger.Error(errors.Wrap(err, "RegulatorModel.sensorOnCheckHandler"))
	//return
	//}
	//
	//maxSP, err := o.GetProps().GetFloatValue("max_sp")
	//if err != nil {
	//context.Logger.Error(errors.Wrap(err, "RegulatorModel.sensorOnCheckHandler"))
	//return
	//}

	targetSP, err := o.GetProps().GetFloatValue("target_sp")
	if err != nil {
		g.Logger.Error(errors.Wrap(err, "RegulatorModel.sensorOnCheckHandler"))
		return
	}

	aboveTolerance, err := o.GetProps().GetFloatValue("above_tolerance")
	if err != nil {
		g.Logger.Error(errors.Wrap(err, "RegulatorModel.sensorOnCheckHandler"))
		return
	}

	belowTolerance, err := o.GetProps().GetFloatValue("below_tolerance")
	if err != nil {
		g.Logger.Error(errors.Wrap(err, "RegulatorModel.sensorOnCheckHandler"))
		return
	}

	complexTolerance, err := o.GetProps().GetFloatValue("complex_tolerance")
	if err != nil {
		g.Logger.Error(errors.Wrap(err, "RegulatorModel.sensorOnCheckHandler"))
		return
	}

	parentID := o.GetParentID()
	if parentID == nil {
		g.Logger.Error(errors.Wrap(errors.Errorf("parent_id of %d is nil", o.GetID()), "RegulatorModel.sensorOnCheckHandler"))
		return
	}

	// Запрос данных TargetSensor
	sensorValue, err := o.requestSensorValue(*parentID)
	if err != nil {
		if fallbackSensorValueID < 1 {
			g.Logger.Error(errors.Wrap(err, "RegulatorModel.sensorOnCheckHandler"))
			return
		}

		sensorValue, err = o.requestSensorValue(fallbackSensorValueID)
		if err != nil {
			g.Logger.Error(errors.Wrap(err, "RegulatorModel.sensorOnCheckHandler"))
			return
		}
	}

	var msg interfaces.Message

	switch regType {
	case TypeSimple:
		switch {
		case sensorValue < (targetSP - belowTolerance): // TODO minSP ???
			msg, err = regulator.NewOnBelow(*parentID, sensorValue)
		case sensorValue > (targetSP + aboveTolerance): // TODO maxSP ???
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
		g.Logger.Error(errors.Wrap(errors.Errorf("unexpected regulator type %q", regType), "RegulatorModel.sensorOnCheckHandler"))
		return
	}

	if err != nil {
		g.Logger.Error(errors.Wrap(err, "RegulatorModel.sensorOnCheckHandler"))
		return
	}

	if err := svc.Send(msg); err != nil {
		g.Logger.Error(errors.Wrap(err, "RegulatorModel.sensorOnCheckHandler"))
	}
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

	parentID := o.GetParentID()
	if parentID == nil {
		return errors.Wrap(errors.Errorf("parent_id of %d is nil", o.GetID()), "RegulatorModel.Start")
	}

	sensorValue, err := store.I.ObjectRepository().GetObject(*parentID)
	if err != nil {
		return errors.Wrap(err, "RegulatorModel.Start")
	}

	err = o.Subscribe(
		interfaces.MessageTypeEvent,
		"object.sensor.on_check",
		interfaces.TargetTypeObject,
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

	g.Logger.Debugf("Regulator(%d) started", o.GetID())

	return nil
}

func (o *RegulatorModel) timerHandler() {
	if !o.GetEnabled() {
		return
	}

	g.Logger.Debugf("Regulator(%d): timerHandler(%s)", o.GetID(), o.timer.GetDuration())

	parentID := o.GetParentID()
	if parentID == nil {
		g.Logger.Error(errors.Wrap(errors.Errorf("parent_id of %d is nil", o.GetID()), "RegulatorModel.timerHandler"))
		return

	}

	msg, err := regulator.NewOnStale(*parentID)
	if err != nil {
		g.Logger.Error(errors.Wrap(err, "RegulatorModel.timerHandler"))
		return
	}

	if err := g.Msgs.Send(msg); err != nil {
		g.Logger.Error(errors.Wrap(err, "RegulatorModel.timerHandler"))
		return
	}

	o.timer.Reset()
}

func (o *RegulatorModel) Shutdown() error {
	if err := o.ObjectModelImpl.Shutdown(); err != nil {
		return errors.Wrap(err, "RegulatorModel.Shutdown")
	}

	o.timer.Stop()

	g.Logger.Debugf("Regulator(%d) stopped", o.GetID())

	return nil
}
