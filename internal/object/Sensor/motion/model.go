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
	"touchon-server/lib/event"
	"touchon-server/lib/events"
	"touchon-server/lib/events/object/sensor"
	"touchon-server/lib/helpers"
	"touchon-server/lib/models"
	mqttClient "touchon-server/lib/mqtt/client"
	"touchon-server/lib/mqtt/messages"
)

func init() {
	_ = objects.Register(MakeModel)
}

func MakeModel() (objects.Object, error) {
	baseObj, err := Sensor.MakeModel()
	if err != nil {
		return nil, errors.Wrap(err, "motion.MakeModel")
	}

	obj := &SensorModel{}
	obj.SensorModel = baseObj.(*Sensor.SensorModel)

	obj.SetType("motion")
	obj.SetName("Датчик движения")
	obj.SetTags("motion", "движение")

	// Датчик получает состояние через сообщения mqtt
	obj.GetProps().Delete("update_interval")

	e := &objects.Prop{
		Code:        "enable",
		Name:        "Состояние датчика",
		Description: "вкл/выкл",
		Item: &models.Item{
			Type:         models.DataTypeBool,
			DefaultValue: true,
		},
		Required: objects.NewRequired(true),
		Editable: objects.NewCondition(),
		Visible:  objects.NewCondition(),
	}

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

	if err := obj.GetProps().Add(e, p, mode); err != nil {
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
	for _, eventName := range []string{"object.sensor.on_motion_on", "object.sensor.on_motion_off"} {
		ev, err := event.MakeEvent(eventName, messages.TargetTypeObject, 0, nil)
		if err != nil {
			return nil, errors.Wrap(err, "motion.MakeModel")
		}

		if err := obj.GetEvents().Add(ev); err != nil {
			return nil, errors.Wrap(err, "motion.MakeModel")
		}
	}

	return obj, nil
}

type SensorModel struct {
	*Sensor.SensorModel
	periodTimer *helpers.Timer
}

func (o *SensorModel) Check(args map[string]interface{}) ([]messages.Message, error) {
	// Данный датчик сам не проверяет значения, а получает значения от порта
	//TODO:: реализовать метод check всё равно, т.к. пользователь может запросить состояние датчика
	return nil, errors.Wrap(errors.New("method 'check' not supported"), "motion.SensorModel.Check")
}

func (o *SensorModel) Start() error {
	if err := o.SensorModel.Start(); err != nil {
		return errors.Wrap(err, "motion.SensorModel.Start")
	}

	enable, err := o.GetProps().GetBoolValue("enable")
	if err != nil {
		return errors.Wrap(err, "motion.SensorModel.Start")
	}

	if !enable {
		return nil
	}

	address, err := o.GetProps().GetStringValue("address")
	if err != nil {
		return errors.Wrap(err, "motion.SensorModel.Start")
	}

	portID, err := strconv.Atoi(address)
	if err != nil {
		return errors.Wrap(err, "motion.SensorModel.Start")
	}

	err = o.Subscribe(
		"",
		"",
		messages.MessageTypeEvent,
		"object.port.on_long_press",
		messages.TargetTypeObject,
		&portID,
		o.onMotionOnHandler,
	)
	if err != nil {
		return errors.Wrap(err, "motion.SensorModel.Start")
	}

	err = o.Subscribe(
		"",
		"",
		messages.MessageTypeEvent,
		"object.port.on_release",
		messages.TargetTypeObject,
		&portID,
		o.onMotionOffHandler,
	)
	if err != nil {
		return errors.Wrap(err, "motion.SensorModel.Start")
	}

	context.Logger.Debugf("motion.SensorModel(%d) started", o.GetID())

	period, err := o.GetProps().GetIntValue("period")
	if err != nil {
		return errors.Wrap(err, "motion.SensorModel.Start")
	}

	o.periodTimer = helpers.NewTimer(time.Duration(period)*time.Second, o.periodTimerHandler)

	// Получаем текущее состояние движения
	state, err := o.getMotionState()
	if err != nil {
		return errors.Wrap(err, "motion.SensorModel.Start")
	}

	// Если есть движение, отправляем сообщение в шину
	if state {
		msg, err := sensor.NewOnMotionOnMessage("object_manager/object/event", o.GetID())
		if err != nil {
			return errors.Wrap(err, "motion.SensorModel.Start")
		}

		if err := mqttClient.I.Send(msg); err != nil {
			return errors.Wrap(err, "motion.SensorModel.Start")
		}
	}

	return nil
}

func (o *SensorModel) onMotionOnHandler(messages.Message) ([]messages.Message, error) {
	context.Logger.Debugf("motion.SensorModel(%d): onMotionOnHandler()", o.GetID())

	// Запоминаем текущее состояние движения
	currState, err := o.getMotionState()
	if err != nil {
		return nil, errors.Wrap(err, "motion.SensorModel.onMotionOnHandler")
	}

	context.Logger.Debug("motion.SensorModel.onMotionOnHandler: reset periodTimer")
	o.periodTimer.Reset()

	// Обрабатываем только переход OFF -> ON
	if currState {
		return nil, nil
	}

	if err := o.setMotionState(true, true); err != nil {
		return nil, errors.Wrap(err, "motion.SensorModel.onMotionOnHandler")
	}

	msg, err := sensor.NewOnMotionOnMessage("object_manager/object/event", o.GetID())
	if err != nil {
		return nil, errors.Wrap(err, "motion.SensorModel.onMotionOnHandler")
	}

	return []messages.Message{msg}, nil
}

func (o *SensorModel) onMotionOffHandler(messages.Message) ([]messages.Message, error) {
	enable, err := o.GetProps().GetBoolValue("enable")
	if err != nil {
		return nil, errors.Wrap(err, "motion.SensorModel.onMotionOffHandler")
	}

	if !enable {
		return nil, errors.Wrap(err, "motion.SensorModel.onMotionOffHandler")
	}

	context.Logger.Debugf("motion.SensorModel(%d): onMotionOffHandler()", o.GetID())

	if err := o.setMotionState(false, false); err != nil {
		return nil, errors.Wrap(err, "motion.SensorModel.onMotionOffHandler")
	}

	return nil, nil
}

func (o *SensorModel) periodTimerHandler() {
	enable, err := o.GetProps().GetBoolValue("enable")
	if err != nil {
		context.Logger.Error(errors.Wrap(err, "motion.SensorModel.periodTimerHandler"))
		return
	}

	if !enable {
		return
	}

	context.Logger.Debugf("motion.SensorModel(%d): periodTimerHandler()", o.GetID())

	// получаем текущее значение движения
	currState, err := o.getMotionState()
	if err != nil {
		context.Logger.Error(errors.Wrap(err, "motion.SensorModel.periodTimerHandler"))
		return
	}

	// если движение есть - перезапускаем таймер и уходим
	if currState {
		o.periodTimer.Reset()
		return
	}

	// если при срабатывании таймера движения не было - выставляем статус о завершении движения
	if err := o.setMotionState(false, true); err != nil {
		context.Logger.Error(errors.Wrap(err, "motion.SensorModel.periodTimerHandler"))
		return
	}

	msg, err := sensor.NewOnMotionOffMessage("object_manager/object/event", o.GetID())
	if err != nil {
		context.Logger.Error(errors.Wrap(err, "motion.SensorModel.periodTimerHandler"))

		msg, err = events.NewOnErrorMessage("object_manager/error", msg.GetTargetType(), msg.GetTargetID(), err.Error())
		if err != nil {
			context.Logger.Error(errors.Wrap(err, "motion.SensorModel.periodTimerHandler"))
			return
		}
	}

	if err := mqttClient.I.Send(msg); err != nil {
		context.Logger.Error(errors.Wrap(err, "motion.SensorModel.periodTimerHandler"))
		return
	}
}

func (o *SensorModel) getMotionState() (bool, error) {
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

func (o *SensorModel) setMotionState(state, save bool) error {
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

func (o *SensorModel) Shutdown() error {
	if err := o.SensorModel.Shutdown(); err != nil {
		return errors.Wrap(err, "motion.SensorModel.Shutdown")
	}

	enable, err := o.GetProps().GetBoolValue("enable")
	if err != nil {
		return errors.Wrap(err, "motion.SensorModel.Shutdown")
	}

	if !enable {
		return nil
	}

	context.Logger.Debugf("motion.SensorModel(%d) stopped", o.GetID())

	return nil
}
