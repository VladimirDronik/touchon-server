package Relay

import (
	"strings"

	"github.com/pkg/errors"
	"touchon-server/internal/g"
	"touchon-server/internal/model"
	"touchon-server/internal/objects"
	"touchon-server/lib/events"
	"touchon-server/lib/events/object/relay"
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
			Required: objects.True(),
			Editable: objects.True(),
			Visible:  objects.True(),
		},
		{
			Code:        "interface",
			Name:        "Интерфейс, по которому подключено реле",
			Description: "",
			Item: &models.Item{
				Type:         models.DataTypeString,
				DefaultValue: "MEGA-OUT",
			},
			Required: objects.True(),
			Editable: objects.True(),
			Visible:  objects.True(),
		},
	}

	onChangeState, err := events.NewOnChangeState(interfaces.TargetTypeObject, 0, "", "")
	if err != nil {
		return nil, errors.Wrap(err, "Relay.MakeModel")
	}

	onStateOn, err := relay.NewOnStateOn(0)
	if err != nil {
		return nil, errors.Wrap(err, "Relay.MakeModel")
	}

	onStateOff, err := relay.NewOnStateOff(0)
	if err != nil {
		return nil, errors.Wrap(err, "Relay.MakeModel")
	}

	onCheck, err := relay.NewOnCheck(0, "", "")
	if err != nil {
		return nil, errors.Wrap(err, "Relay.MakeModel")
	}

	impl, err := objects.NewObjectModelImpl(
		model.CategoryRelay,
		"relay",
		false,
		"Реле",
		props,
		nil,
		[]interfaces.Event{onChangeState, onStateOn, onStateOff, onCheck},
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

	g.Logger.Debugf("Relay(%d) started", o.GetID())

	return nil
}

func (o *RelayModel) handler(svc interfaces.MessageSender, msg interfaces.Message) {
	g.Logger.Debugf("Relay(%d): handler()", o.GetID())

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
		g.Logger.Error(errors.Wrap(err, "RelayModel.handler"))
		return
	}

	if err := svc.Send(msg); err != nil {
		g.Logger.Error(errors.Wrap(err, "RelayModel.handler"))
	}
}

func (o *RelayModel) Shutdown() error {
	g.Logger.Debugf("RelayModel(%d) stopped", o.GetID())

	return nil
}
