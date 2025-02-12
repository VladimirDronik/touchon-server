package Relay

import (
	"strings"

	"github.com/pkg/errors"
	"touchon-server/internal/context"
	"touchon-server/internal/model"
	"touchon-server/internal/msgs"
	"touchon-server/internal/objects"
	"touchon-server/lib/events"
	"touchon-server/lib/interfaces"
	"touchon-server/lib/models"
)

func init() {
	_ = objects.Register(MakeModel)
}

func MakeModel() (objects.Object, error) {
	props := []*objects.Prop{
		{
			Code:        "address",
			Name:        "ИД объекта порта, на котором находится реле",
			Description: "",
			Item: &models.Item{
				Type:         models.DataTypeInt,
				DefaultValue: 0,
			},
			Required: objects.NewRequired(true),
			Editable: objects.NewCondition(),
			Visible:  objects.NewCondition(),
		},
		{
			Code:        "interface",
			Name:        "Интерфейс, по которому подключено реле",
			Description: "",
			Item: &models.Item{
				Type:         models.DataTypeString,
				DefaultValue: "MEGA-OUT",
			},
			Required: objects.NewRequired(true),
			Editable: objects.NewCondition(),
			Visible:  objects.NewCondition(),
		},
	}

	impl, err := objects.NewObjectModelImpl(
		model.CategoryRelay,
		"relay",
		false,
		"Реле",
		props,
		nil,
		[]string{
			"on_change_state",
			"object.relay.on_state_on",
			"object.relay.on_state_off",
			"object.relay.on_check",
		},
		nil,
		[]string{"Лампа", "relay", "output", "Насос", "Вентилятор", "Розетка", "Реле"},
	)
	if err != nil {
		return nil, errors.Wrap(err, "Relay.MakeModel")
	}

	o := &RelayModel{ObjectModelImpl: impl}

	on, err := objects.NewMethod("on", "Включение реле", nil, o.On)
	if err != nil {
		return nil, errors.Wrap(err, "Relay.MakeModel")
	}

	check, err := objects.NewMethod("check", "Получение статуса реле", nil, o.Check)
	if err != nil {
		return nil, errors.Wrap(err, "Relay.MakeModel")
	}

	off, err := objects.NewMethod("off", "Выключение реле", nil, o.Off)
	if err != nil {
		return nil, errors.Wrap(err, "Relay.MakeModel")
	}

	toggle, err := objects.NewMethod("toggle", "Переключение состояния реле вкл/выкл", nil, o.Toggle)
	if err != nil {
		return nil, errors.Wrap(err, "Relay.MakeModel")
	}

	o.GetMethods().Add(on, off, toggle, check)

	return o, nil
}

type RelayModel struct {
	*objects.ObjectModelImpl
}

func (o *RelayModel) Start() error {
	if err := o.ObjectModelImpl.Start(); err != nil {
		return errors.Wrap(err, "Relay.Start")
	}

	portID, err := o.GetProps().GetIntValue("address")
	if err != nil {
		return errors.Wrap(err, "Relay.Start")
	}

	err = o.Subscribe(
		interfaces.MessageTypeEvent,
		"",
		interfaces.TargetTypeObject,
		&portID,
		o.handler,
	)
	if err != nil {
		return errors.Wrap(err, "Relay.Start")
	}

	context.Logger.Debugf("Relay(%d) started", o.GetID())

	return nil
}

func (o *RelayModel) handler(msg interfaces.Message) {
	context.Logger.Debugf("Relay(%d): handler()", o.GetID())

	var err error

	switch msg.GetName() {
	case "object.port.on_change_state":
		// TODO
		var state string // = msg.GetPayload()["status"].(string)
		msg, err = events.NewOnChangeState(interfaces.TargetTypeObject, o.GetID(), strings.ToLower(state), "")
	default:
		return
	}

	if err != nil {
		context.Logger.Error(errors.Wrap(err, "RelayModel.handler"))
		return
	}

	if err := msgs.I.Send(msg); err != nil {
		context.Logger.Error(errors.Wrap(err, "RelayModel.handler"))
	}
}

func (o *RelayModel) Shutdown() error {
	context.Logger.Debugf("RelayModel(%d) stopped", o.GetID())

	return nil
}
