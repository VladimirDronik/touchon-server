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
	"touchon-server/lib/event"
	"touchon-server/lib/events/object/sensor"
	"touchon-server/lib/mqtt/messages"
)

func init() {
	_ = objects.Register(MakeModel)
}

func MakeModel() (objects.Object, error) {
	baseObj, err := motion.MakeModel()
	if err != nil {
		return nil, errors.Wrap(err, "presence.MakeModel")
	}

	obj := &SensorModel{}
	obj.SensorModel = baseObj.(*motion.SensorModel)

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
	for _, eventName := range []string{"object.sensor.on_presence_on", "object.sensor.on_presence_off"} {
		ev, err := event.MakeEvent(eventName, messages.TargetTypeObject, 0, nil)
		if err != nil {
			return nil, errors.Wrap(err, "presence.MakeModel")
		}

		if err := obj.GetEvents().Add(ev); err != nil {
			return nil, errors.Wrap(err, "presence.MakeModel")
		}
	}

	return obj, nil
}

type SensorModel struct {
	*motion.SensorModel
}

func (o *SensorModel) Start() error {
	enable, err := o.GetProps().GetBoolValue("enable")
	if err != nil {
		return errors.Wrap(err, "presence.SensorModel.Start")
	}

	if !enable {
		return nil
	}

	address, err := o.GetProps().GetStringValue("address")
	if err != nil {
		return errors.Wrap(err, "presence.SensorModel.Start")
	}

	portIDs := strings.Split(address, ";")
	if len(portIDs) != 2 {
		return errors.Wrap(errors.New("address is bad"), "presence.SensorModel.Start")
	}

	presencePortID, err := strconv.Atoi(portIDs[1])
	if err != nil {
		return errors.Wrap(err, "presence.SensorModel.Start")
	}

	// Заменяем "<motion port id>;<presence port id>" на "<motion port id>"
	if err := o.GetProps().Set("address", portIDs[0]); err != nil {
		return errors.Wrap(err, "presence.SensorModel.Start")
	}

	// Выполняем логику датчика движения
	if err := o.SensorModel.Start(); err != nil {
		return errors.Wrap(err, "presence.SensorModel.Start")
	}

	// Возвращаем оригинальное значение адреса
	if err := o.GetProps().Set("address", address); err != nil {
		return errors.Wrap(err, "presence.SensorModel.Start")
	}

	err = o.Subscribe(
		"",
		"",
		messages.MessageTypeEvent,
		"object.port.on_long_press",
		messages.TargetTypeObject,
		&presencePortID,
		o.onPresenceOnHandler,
	)
	if err != nil {
		return errors.Wrap(err, "presence.SensorModel.Start")
	}

	err = o.Subscribe(
		"",
		"",
		messages.MessageTypeEvent,
		"object.port.on_release",
		messages.TargetTypeObject,
		&presencePortID,
		o.onPresenceOffHandler,
	)
	if err != nil {
		return errors.Wrap(err, "presence.SensorModel.Start")
	}

	context.Logger.Debugf("presence.SensorModel(%d) started", o.GetID())

	return nil
}

func (o *SensorModel) onPresenceOnHandler(messages.Message) ([]messages.Message, error) {
	enable, err := o.GetProps().GetBoolValue("enable")
	if err != nil {
		return nil, errors.Wrap(err, "presence.SensorModel.onPresenceOnHandler")
	}

	if !enable {
		return nil, nil
	}

	context.Logger.Debugf("presence.SensorModel(%d): onPresenceOnHandler()", o.GetID())

	// получаем текущее значение движения
	currState, err := o.getPresenceState()
	if err != nil {
		return nil, errors.Wrap(err, "presence.SensorModel.onPresenceOnHandler")
	}

	// Если уже true, уходим
	if currState {
		return nil, nil
	}

	if err := o.setPresenceState(true); err != nil {
		return nil, errors.Wrap(err, "presence.SensorModel.onPresenceOnHandler")
	}

	msg, err := sensor.NewOnPresenceOnMessage("object_manager/object/event", o.GetID())
	if err != nil {
		return nil, errors.Wrap(err, "presence.SensorModel.onPresenceOnHandler")
	}

	return []messages.Message{msg}, nil
}

func (o *SensorModel) onPresenceOffHandler(messages.Message) ([]messages.Message, error) {
	enable, err := o.GetProps().GetBoolValue("enable")
	if err != nil {
		return nil, errors.Wrap(err, "presence.SensorModel.onPresenceOffHandler")
	}

	if !enable {
		return nil, nil
	}

	context.Logger.Debugf("presence.SensorModel(%d): onPresenceOffHandler()", o.GetID())

	// получаем текущее значение движения
	currState, err := o.getPresenceState()
	if err != nil {
		return nil, errors.Wrap(err, "presence.SensorModel.onPresenceOffHandler")
	}

	// Если уже false, уходим
	if !currState {
		return nil, nil
	}

	if err := o.setPresenceState(false); err != nil {
		return nil, errors.Wrap(err, "presence.SensorModel.onPresenceOffHandler")
	}

	msg, err := sensor.NewOnPresenceOffMessage("object_manager/object/event", o.GetID())
	if err != nil {
		return nil, errors.Wrap(err, "presence.SensorModel.onPresenceOffHandler")
	}

	return []messages.Message{msg}, nil
}

func (o *SensorModel) getPresenceState() (bool, error) {
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

func (o *SensorModel) setPresenceState(state bool) error {
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

func (o *SensorModel) Shutdown() error {
	if err := o.SensorModel.Shutdown(); err != nil {
		return errors.Wrap(err, "presence.SensorModel.Shutdown")
	}

	enable, err := o.GetProps().GetBoolValue("enable")
	if err != nil {
		return errors.Wrap(err, "presence.SensorModel.Shutdown")
	}

	if !enable {
		return nil
	}

	context.Logger.Debugf("presence.SensorModel(%d) stopped", o.GetID())

	return nil
}
