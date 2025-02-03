package PortMegaD

import (
	"strings"

	"github.com/pkg/errors"
	"touchon-server/internal/context"
	"touchon-server/lib/events"
	"touchon-server/lib/events/object/port"
	"touchon-server/lib/mqtt/messages"
)

// OnPress событие генерируется при замыкании порта
func (o *PortModel) OnPress() messages.Message {
	msg, err := port.NewOnPressMessage("object_manager/object/event", o.GetID())
	if err != nil {
		context.Logger.Error(errors.Wrap(err, "OnPress"))
		return nil
	}

	return msg
}

// OnRelease событие генерируется при отпускании порта
func (o *PortModel) OnRelease() messages.Message {
	msg, err := port.NewOnReleaseMessage("object_manager/object/event", o.GetID())
	if err != nil {
		context.Logger.Error(errors.Wrap(err, "OnRelease"))
		return nil
	}

	return msg
}

// OnLongPress событие генерируется при удержании
func (o *PortModel) OnLongPress() messages.Message {
	msg, err := port.NewOnLongPressMessage("object_manager/object/event", o.GetID())
	if err != nil {
		context.Logger.Error(errors.Wrap(err, "OnLongPress"))
		return nil
	}

	return msg
}

func (o *PortModel) OnDoubleClick() messages.Message {
	msg, err := port.NewOnDoubleClickMessage("object_manager/object/event", o.GetID())
	if err != nil {
		context.Logger.Error(errors.Wrap(err, "OnDoubleClick"))
		return nil
	}

	return msg
}

// OnChangeState Событие, которое возникает при смене статуса объекта
func (o *PortModel) OnChangeState(state string) messages.Message {
	var value string

	mode, err := o.GetProps().GetStringValue("mode")
	if err != nil {
		context.Logger.Error(errors.Wrap(err, "OnChangeState"))
		return nil
	}

	if mode == "PWM" {
		value = state
		state = ""
	}

	msg, err := events.NewOnChangeStateMessage("object_manager/object/event", messages.TargetTypeObject, o.GetID(), strings.ToLower(state), value)
	if err != nil {
		context.Logger.Error(errors.Wrap(err, "OnChangeState"))
		return nil
	}

	return msg
}

// OnCheck событие, которое возникает, когда проверяем состояние порта, но при этом новое пришедшее состояние порта
// не различается с тем, что хранится в БД
func (o *PortModel) OnCheck(state string) messages.Message {
	var value string

	mode, err := o.GetProps().GetStringValue("mode")
	if err != nil {
		context.Logger.Error(errors.Wrap(err, "OnCheck"))
		return nil
	}

	if mode == "PWM" {
		value = state
		state = ""
	}

	msg, err := port.NewOnCheckMessage("object_manager/object/event", o.GetID(), state, value)
	if err != nil {
		context.Logger.Error(errors.Wrap(err, "OnCheck"))
	}

	return msg
}
