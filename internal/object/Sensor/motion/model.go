package motion

import (
	"strconv"
	"time"

	"github.com/pkg/errors"
	"touchon-server/internal/context"
	"touchon-server/internal/model"
	"touchon-server/internal/object/Sensor"
	"touchon-server/internal/object/SensorValue"
	"touchon-server/internal/objects"
	"touchon-server/lib/events"
	"touchon-server/lib/events/object/sensor"
	"touchon-server/lib/helpers"
	"touchon-server/lib/interfaces"
	msgs "touchon-server/lib/messages"
	"touchon-server/lib/models"
)

func init() {
	_ = objects.Register(MakeModel)
}

func MakeModel() (objects.Object, error) {
	baseObj, err := Sensor.MakeModel()
	if err != nil {
		return nil, errors.Wrap(err, "motion.MakeModel")
	}

	obj := &MotionSensorModel{}
	obj.SensorModel = baseObj.(*Sensor.SensorModel)

	obj.SetType("motion")
	obj.SetName("Датчик движения")
	obj.SetTags("motion", "движение")

	// Датчик получает состояние через сообщения mqtt
	obj.GetProps().Delete("update_interval")

	p := &objects.Prop{
		Code:        "period",
		Name:        "Период (с)",
		Description: "Время, в течение которого система будет считать, что есть движение",
		Item: &models.Item{
			Type:         models.DataTypeInt,
			DefaultValue: 120,
		},
		Required:   objects.NewRequired(true),
		Editable:   objects.NewCondition(),
		Visible:    objects.NewCondition(),
		CheckValue: objects.AboveOrEqual1(),
	}

	mode := &objects.Prop{
		Code:        "mode",
		Name:        "Режим порта",
		Description: "Режим порта для контроллера",
		Item: &models.Item{
			Type:         models.DataTypeString,
			DefaultValue: "P",
		},
		Required:   objects.NewRequired(true),
		Editable:   objects.NewCondition(),
		Visible:    objects.NewCondition(),
		CheckValue: objects.AboveOrEqual1(),
	}

	if err := obj.GetProps().Add(p, mode); err != nil {
		return nil, errors.Wrap(err, "motion.MakeModel")
	}

	if err := obj.GetProps().Set("interface", "MEGA-IN"); err != nil {
		return nil, errors.Wrap(err, "motion.MakeModel")
	}

	m, err := SensorValue.Make(SensorValue.TypeMotion)
	if err != nil {
		return nil, errors.Wrap(err, "motion.MakeModel")
	}

	if err := m.GetProps().Set("value", 0); err != nil {
		return nil, errors.Wrap(err, "motion.MakeModel")
	}

	// Скрываем св-ва значения датчика, не используемые в данном типе датчика
	for _, code := range []string{"unit", "min_error_value", "min_threshold", "max_threshold", "max_error_value"} {
		prop, err := m.GetProps().Get(code)
		if err != nil {
			return nil, errors.Wrap(err, "motion.MakeModel")
		}

		prop.Required = objects.NewRequired(false)
		prop.Editable = objects.NewCondition().AccessLevel(model.AccessLevelDenied)
		prop.Visible = objects.NewCondition().AccessLevel(model.AccessLevelDenied)
		prop.CheckValue = nil
	}

	// Удаляем регулятор
	m.GetChildren().DeleteAll()

	obj.GetChildren().Add(m)

	// Удаляем лишние события
	obj.GetEvents().Delete("object.sensor.on_check", "object.sensor.on_alarm")

	// Добавляем свои события
	onMotionOn, err := sensor.NewOnMotionOn(0)
	if err != nil {
		return nil, errors.Wrap(err, "motion.MakeModel")
	}

	onMotionOff, err := sensor.NewOnMotionOff(0)
	if err != nil {
		return nil, errors.Wrap(err, "motion.MakeModel")
	}

	if err := obj.GetEvents().Add(onMotionOn, onMotionOff); err != nil {
		return nil, errors.Wrap(err, "motion.MakeModel")
	}

	return obj, nil
}

type MotionSensorModel struct {
	*Sensor.SensorModel
	periodTimer *helpers.Timer
}

func (o *MotionSensorModel) Check(args map[string]interface{}) ([]interfaces.Message, error) {
	// Данный датчик сам не проверяет значения, а получает значения от порта
	//TODO:: реализовать метод check всё равно, т.к. пользователь может запросить состояние датчика
	return nil, errors.Wrap(errors.New("method 'check' not supported"), "MotionSensorModel.Check")
}

func (o *MotionSensorModel) Start() error {
	if err := o.SensorModel.Start(); err != nil {
		return errors.Wrap(err, "MotionSensorModel.Start")
	}

	address, err := o.GetProps().GetStringValue("address")
	if err != nil {
		return errors.Wrap(err, "MotionSensorModel.Start")
	}

	portID, err := strconv.Atoi(address)
	if err != nil {
		return errors.Wrap(err, "MotionSensorModel.Start")
	}

	err = o.Subscribe(
		interfaces.MessageTypeEvent,
		"object.port.on_long_press",
		interfaces.TargetTypeObject,
		&portID,
		o.onMotionOnHandler,
	)
	if err != nil {
		return errors.Wrap(err, "MotionSensorModel.Start")
	}

	err = o.Subscribe(
		interfaces.MessageTypeEvent,
		"object.port.on_release",
		interfaces.TargetTypeObject,
		&portID,
		o.onMotionOffHandler,
	)
	if err != nil {
		return errors.Wrap(err, "MotionSensorModel.Start")
	}

	context.Logger.Debugf("MotionSensorModel(%d) started", o.GetID())

	period, err := o.GetProps().GetIntValue("period")
	if err != nil {
		return errors.Wrap(err, "MotionSensorModel.Start")
	}

	o.periodTimer = helpers.NewTimer(time.Duration(period)*time.Second, o.periodTimerHandler)

	// Получаем текущее состояние движения
	state, err := o.getMotionState()
	if err != nil {
		return errors.Wrap(err, "MotionSensorModel.Start")
	}

	// Если есть движение, отправляем сообщение в шину
	if state {
		msg, err := sensor.NewOnMotionOn(o.GetID())
		if err != nil {
			return errors.Wrap(err, "MotionSensorModel.Start")
		}

		if err := msgs.I.Send(msg); err != nil {
			return errors.Wrap(err, "MotionSensorModel.Start")
		}
	}

	return nil
}

func (o *MotionSensorModel) onMotionOnHandler(interfaces.Message) {
	context.Logger.Debugf("MotionSensorModel(%d): onMotionOnHandler()", o.GetID())

	// Запоминаем текущее состояние движения
	currState, err := o.getMotionState()
	if err != nil {
		context.Logger.Error(errors.Wrap(err, "MotionSensorModel.onMotionOnHandler"))
		return
	}

	context.Logger.Debug("MotionSensorModel.onMotionOnHandler: reset periodTimer")
	o.periodTimer.Reset()

	// Обрабатываем только переход OFF -> ON
	if currState {
		return
	}

	if err := o.setMotionState(true, true); err != nil {
		context.Logger.Error(errors.Wrap(err, "MotionSensorModel.onMotionOnHandler"))
		return
	}

	msg, err := sensor.NewOnMotionOn(o.GetID())
	if err != nil {
		context.Logger.Error(errors.Wrap(err, "MotionSensorModel.onMotionOnHandler"))
		return
	}

	if err := msgs.I.Send(msg); err != nil {
		context.Logger.Error(errors.Wrap(err, "MotionSensorModel.onMotionOnHandler"))
	}
}

func (o *MotionSensorModel) onMotionOffHandler(interfaces.Message) {
	if err := o.CheckEnabled(); err != nil {
		context.Logger.Error(errors.Wrap(err, "MotionSensorModel.onMotionOffHandler"))
		return
	}

	context.Logger.Debugf("MotionSensorModel(%d): onMotionOffHandler()", o.GetID())

	if err := o.setMotionState(false, false); err != nil {
		context.Logger.Error(errors.Wrap(err, "MotionSensorModel.onMotionOffHandler"))
		return
	}
}

func (o *MotionSensorModel) periodTimerHandler() {
	if !o.GetEnabled() {
		return
	}

	context.Logger.Debugf("MotionSensorModel(%d): periodTimerHandler()", o.GetID())

	// получаем текущее значение движения
	currState, err := o.getMotionState()
	if err != nil {
		context.Logger.Error(errors.Wrap(err, "MotionSensorModel.periodTimerHandler"))
		return
	}

	// если движение есть - перезапускаем таймер и уходим
	if currState {
		o.periodTimer.Reset()
		return
	}

	// если при срабатывании таймера движения не было - выставляем статус о завершении движения
	if err := o.setMotionState(false, true); err != nil {
		context.Logger.Error(errors.Wrap(err, "MotionSensorModel.periodTimerHandler"))
		return
	}

	msg, err := sensor.NewOnMotionOff(o.GetID())
	if err != nil {
		context.Logger.Error(errors.Wrap(err, "MotionSensorModel.periodTimerHandler"))

		msg, err = events.NewOnError(msg.GetTargetType(), msg.GetTargetID(), err.Error())
		if err != nil {
			context.Logger.Error(errors.Wrap(err, "MotionSensorModel.periodTimerHandler"))
			return
		}
	}

	if err := msgs.I.Send(msg); err != nil {
		context.Logger.Error(errors.Wrap(err, "MotionSensorModel.periodTimerHandler"))
		return
	}
}

func (o *MotionSensorModel) getMotionState() (bool, error) {
	for _, child := range o.GetChildren().GetAll() {
		if child.GetCategory() == model.CategorySensorValue && child.GetType() == SensorValue.TypeMotion {
			v, err := child.GetProps().GetFloatValue("value")
			if err != nil {
				return false, errors.Wrap(err, "getMotionState")
			}

			switch v {
			case 0:
				return false, nil
			case 1:
				return true, nil
			default:
				return false, errors.Wrap(errors.Errorf("unexpected motion value: %f", v), "getMotionState")
			}
		}
	}

	return false, errors.Wrap(errors.New("child SensorValue.Motion not found"), "getMotionState")
}

func (o *MotionSensorModel) setMotionState(state, save bool) error {
	for _, child := range o.GetChildren().GetAll() {
		if child.GetCategory() == model.CategorySensorValue && child.GetType() == SensorValue.TypeMotion {
			value := 0
			if state {
				value = 1
			}

			if err := child.GetProps().Set("value", value); err != nil {
				return errors.Wrap(err, "setMotionState")
			}

			// В базу сохраняем не всегда
			// При начале движения - сохраняем
			// При завершении движения только если таймер завершен
			if save {
				if err := o.SaveSensorValue(child); err != nil {
					return errors.Wrap(err, "setMotionState")
				}
			}

			return nil
		}
	}

	return errors.Wrap(errors.New("child SensorValue.Motion not found"), "setMotionState")
}

func (o *MotionSensorModel) Shutdown() error {
	if err := o.SensorModel.Shutdown(); err != nil {
		return errors.Wrap(err, "MotionSensorModel.Shutdown")
	}

	context.Logger.Debugf("MotionSensorModel(%d) stopped", o.GetID())

	return nil
}
