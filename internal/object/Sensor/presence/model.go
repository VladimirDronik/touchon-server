package presence

import (
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"touchon-server/internal/context"
	"touchon-server/internal/model"
	"touchon-server/internal/object/Sensor/motion"
	"touchon-server/internal/object/SensorValue"
	"touchon-server/internal/objects"
	"touchon-server/lib/events/object/sensor"
	"touchon-server/lib/interfaces"
)

func init() {
	_ = objects.Register(MakeModel)
}

func MakeModel() (objects.Object, error) {
	baseObj, err := motion.MakeModel()
	if err != nil {
		return nil, errors.Wrap(err, "presence.MakeModel")
	}

	obj := &PresenceSensorModel{}
	obj.MotionSensorModel = baseObj.(*motion.MotionSensorModel)

	obj.SetType("presence")
	obj.SetName("Датчик присутствия")
	obj.SetTags("presence", "присутствие")

	p, err := SensorValue.Make(SensorValue.TypePresence)
	if err != nil {
		return nil, errors.Wrap(err, "presence.MakeModel")
	}

	if err := p.GetProps().Set("value", 0); err != nil {
		return nil, errors.Wrap(err, "presence.MakeModel")
	}

	// Скрываем св-ва значения датчика, не используемые в данном типе датчика
	for _, code := range []string{"unit", "min_error_value", "min_threshold", "max_threshold", "max_error_value"} {
		prop, err := p.GetProps().Get(code)
		if err != nil {
			return nil, errors.Wrap(err, "presence.MakeModel")
		}

		prop.Required = objects.NewRequired(false)
		prop.Editable = objects.NewCondition().AccessLevel(model.AccessLevelDenied)
		prop.Visible = objects.NewCondition().AccessLevel(model.AccessLevelDenied)
		prop.CheckValue = nil
	}

	// Удаляем регулятор
	p.GetChildren().DeleteAll()

	obj.GetChildren().Add(p)

	// Добавляем свои события
	onPresenceOn, err := sensor.NewOnPresenceOn(0)
	if err != nil {
		return nil, errors.Wrap(err, "presence.MakeModel")
	}

	onPresenceOff, err := sensor.NewOnPresenceOff(0)
	if err != nil {
		return nil, errors.Wrap(err, "presence.MakeModel")
	}

	if err := obj.GetEvents().Add(onPresenceOn, onPresenceOff); err != nil {
		return nil, errors.Wrap(err, "presence.MakeModel")
	}

	return obj, nil
}

type PresenceSensorModel struct {
	*motion.MotionSensorModel
}

func (o *PresenceSensorModel) Start() error {
	if err := o.SensorModel.Start(); err != nil {
		return errors.Wrap(err, "PresenceSensorModel.Start")
	}

	address, err := o.GetProps().GetStringValue("address")
	if err != nil {
		return errors.Wrap(err, "PresenceSensorModel.Start")
	}

	portIDs := strings.Split(address, ";")
	if len(portIDs) != 2 {
		return errors.Wrap(errors.New("address is bad"), "PresenceSensorModel.Start")
	}

	presencePortID, err := strconv.Atoi(portIDs[1])
	if err != nil {
		return errors.Wrap(err, "PresenceSensorModel.Start")
	}

	// Заменяем "<motion port id>;<presence port id>" на "<motion port id>"
	if err := o.GetProps().Set("address", portIDs[0]); err != nil {
		return errors.Wrap(err, "PresenceSensorModel.Start")
	}

	// Выполняем логику датчика движения
	if err := o.MotionSensorModel.Start(); err != nil {
		return errors.Wrap(err, "PresenceSensorModel.Start")
	}

	// Возвращаем оригинальное значение адреса
	if err := o.GetProps().Set("address", address); err != nil {
		return errors.Wrap(err, "PresenceSensorModel.Start")
	}

	err = o.Subscribe(
		interfaces.MessageTypeEvent,
		"object.port.on_long_press",
		interfaces.TargetTypeObject,
		&presencePortID,
		o.onPresenceOnHandler,
	)
	if err != nil {
		return errors.Wrap(err, "PresenceSensorModel.Start")
	}

	err = o.Subscribe(
		interfaces.MessageTypeEvent,
		"object.port.on_release",
		interfaces.TargetTypeObject,
		&presencePortID,
		o.onPresenceOffHandler,
	)
	if err != nil {
		return errors.Wrap(err, "PresenceSensorModel.Start")
	}

	context.Logger.Debugf("PresenceSensorModel(%d) started", o.GetID())

	return nil
}

func (o *PresenceSensorModel) onPresenceOnHandler(svc interfaces.MessageSender, _ interfaces.Message) {
	if err := o.CheckEnabled(); err != nil {
		context.Logger.Error(errors.Wrap(err, "PresenceSensorModel.onPresenceOnHandler"))
		return
	}

	context.Logger.Debugf("PresenceSensorModel(%d): onPresenceOnHandler()", o.GetID())

	// получаем текущее значение движения
	currState, err := o.getPresenceState()
	if err != nil {
		context.Logger.Error(errors.Wrap(err, "PresenceSensorModel.onPresenceOnHandler"))
		return
	}

	// Если уже true, уходим
	if currState {
		return
	}

	if err := o.setPresenceState(true); err != nil {
		context.Logger.Error(errors.Wrap(err, "PresenceSensorModel.onPresenceOnHandler"))
		return
	}

	msg, err := sensor.NewOnPresenceOn(o.GetID())
	if err != nil {
		context.Logger.Error(errors.Wrap(err, "PresenceSensorModel.onPresenceOnHandler"))
		return
	}

	if err := svc.Send(msg); err != nil {
		context.Logger.Error(errors.Wrap(err, "PresenceSensorModel.onPresenceOnHandler"))
	}
}

func (o *PresenceSensorModel) onPresenceOffHandler(svc interfaces.MessageSender, _ interfaces.Message) {
	if err := o.CheckEnabled(); err != nil {
		context.Logger.Error(errors.Wrap(err, "PresenceSensorModel.onPresenceOffHandler"))
		return
	}

	context.Logger.Debugf("PresenceSensorModel(%d): onPresenceOffHandler()", o.GetID())

	// получаем текущее значение движения
	currState, err := o.getPresenceState()
	if err != nil {
		context.Logger.Error(errors.Wrap(err, "PresenceSensorModel.onPresenceOffHandler"))
		return
	}

	// Если уже false, уходим
	if !currState {
		return
	}

	if err := o.setPresenceState(false); err != nil {
		context.Logger.Error(errors.Wrap(err, "PresenceSensorModel.onPresenceOffHandler"))
		return
	}

	msg, err := sensor.NewOnPresenceOff(o.GetID())
	if err != nil {
		context.Logger.Error(errors.Wrap(err, "PresenceSensorModel.onPresenceOffHandler"))
		return
	}

	if err := svc.Send(msg); err != nil {
		context.Logger.Error(errors.Wrap(err, "PresenceSensorModel.onPresenceOffHandler"))
	}
}

func (o *PresenceSensorModel) getPresenceState() (bool, error) {
	for _, child := range o.GetChildren().GetAll() {
		if child.GetCategory() == model.CategorySensorValue && child.GetType() == SensorValue.TypePresence {
			v, err := child.GetProps().GetFloatValue("value")
			if err != nil {
				return false, errors.Wrap(err, "getPresenceState")
			}

			switch v {
			case 0:
				return false, nil
			case 1:
				return true, nil
			default:
				return false, errors.Wrap(errors.Errorf("unexpected presence value: %f", v), "getPresenceState")
			}
		}
	}

	return false, errors.Wrap(errors.New("child SensorValue.Presence not found"), "getPresenceState")
}

func (o *PresenceSensorModel) setPresenceState(state bool) error {
	for _, child := range o.GetChildren().GetAll() {
		if child.GetCategory() == model.CategorySensorValue && child.GetType() == SensorValue.TypePresence {
			value := 0
			if state {
				value = 1
			}

			if err := child.GetProps().Set("value", value); err != nil {
				return errors.Wrap(err, "setPresenceState")
			}

			if err := o.SaveSensorValue(child); err != nil {
				return errors.Wrap(err, "setPresenceState")
			}

			return nil
		}
	}

	return errors.Wrap(errors.New("child SensorValue.Presence not found"), "setPresenceState")
}

func (o *PresenceSensorModel) Shutdown() error {
	if err := o.SensorModel.Shutdown(); err != nil {
		return errors.Wrap(err, "PresenceSensorModel.Shutdown")
	}

	context.Logger.Debugf("PresenceSensorModel(%d) stopped", o.GetID())

	return nil
}
