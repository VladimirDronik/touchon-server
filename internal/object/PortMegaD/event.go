package PortMegaD

import (
	"strings"

	"github.com/pkg/errors"
	"touchon-server/internal/g"
	"touchon-server/lib/events"
	"touchon-server/lib/events/object/port"
	"touchon-server/lib/interfaces"
)

// OnPress событие генерируется при замыкании порта
func (o *PortModel) OnPress(countImpulse int) interfaces.Message {
	msg, err := port.NewOnPress(o.GetID(), countImpulse)
	if err != nil {
		g.Logger.Error(errors.Wrap(err, "OnPress"))
		return nil
	}

	return msg
}

// OnRelease событие генерируется при отпускании порта
func (o *PortModel) OnRelease() interfaces.Message {
	msg, err := port.NewOnRelease(o.GetID())
	if err != nil {
		g.Logger.Error(errors.Wrap(err, "OnRelease"))
		return nil
	}

	return msg
}

// OnLongPress событие генерируется при удержании
func (o *PortModel) OnLongPress() interfaces.Message {
	msg, err := port.NewOnLongPress(o.GetID())
	if err != nil {
		g.Logger.Error(errors.Wrap(err, "OnLongPress"))
		return nil
	}

	return msg
}

func (o *PortModel) OnDoubleClick() interfaces.Message {
	msg, err := port.NewOnDoubleClick(o.GetID())
	if err != nil {
		g.Logger.Error(errors.Wrap(err, "OnDoubleClick"))
		return nil
	}

	return msg
}

// OnChangeState Событие, которое возникает при смене статуса объекта
func (o *PortModel) OnChangeState(state string) interfaces.Message {
	var value string

	mode, err := o.GetProps().GetStringValue("mode")
	if err != nil {
		g.Logger.Error(errors.Wrap(err, "OnChangeState"))
		return nil
	}

	if mode == "PWM" {
		value = state
		state = ""
	}

	msg, err := events.NewOnChangeState(interfaces.TargetTypeObject, o.GetID(), strings.ToLower(state), value)
	if err != nil {
		g.Logger.Error(errors.Wrap(err, "OnChangeState"))
		return nil
	}

	msg.SetPayload(map[string]interface{}{"state": state})

	return msg
}

// OnCheck событие, которое возникает, когда проверяем состояние порта, но при этом новое пришедшее состояние порта
// не различается с тем, что хранится в БД
func (o *PortModel) OnCheck(state string) interfaces.Message {
	var value string

	mode, err := o.GetProps().GetStringValue("mode")
	if err != nil {
		g.Logger.Error(errors.Wrap(err, "OnCheck"))
		return nil
	}

	if mode == "PWM" {
		value = state
		state = ""
	}

	msg, err := port.NewOnCheck(o.GetID(), state, value)
	if err != nil {
		g.Logger.Error(errors.Wrap(err, "OnCheck"))
	}

	return msg
}
