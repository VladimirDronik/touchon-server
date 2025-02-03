package GenericInput

import (
	"github.com/pkg/errors"
	"touchon-server/internal/context"
	"touchon-server/internal/model"
	"touchon-server/internal/objects"
	"touchon-server/lib/events/object/generic_input"
	"touchon-server/lib/models"
	"touchon-server/lib/mqtt/messages"
)

func init() {
	_ = objects.Register(MakeModel)
}

func MakeModel() (objects.Object, error) {
	props := []*objects.Prop{
		{
			Code:        "address",
			Name:        "Адрес устройства на шине или ID порта",
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
			Name:        "Интерфейс, по которому подключено устройство",
			Description: "",
			Item: &models.Item{
				Type:         models.DataTypeString,
				DefaultValue: "",
			},
			Required: objects.NewRequired(true),
			Editable: objects.NewCondition(),
			Visible:  objects.NewCondition(),
		},
		{
			Code:        "mode",
			Name:        "Режим работы устройства",
			Description: "",
			Item: &models.Item{
				Type:         models.DataTypeString,
				DefaultValue: "",
			},
			Required: objects.NewRequired(true),
			Editable: objects.NewCondition(),
			Visible:  objects.NewCondition(),
		},
	}

	impl, err := objects.NewObjectModelImpl(
		model.CategoryGenericInput,
		"generic_input",
		false,
		"Универсальный вход",
		props,
		nil,
		[]string{
			"object.generic_input.on_click",
			"object.generic_input.on_double_click",
			"object.generic_input.on_long_press",
		},
		nil,
		[]string{"generic", "input"},
	)
	if err != nil {
		return nil, errors.Wrap(err, "GenericInput.MakeModel")
	}

	return &GenericInputModel{
		ObjectModelImpl: impl,
	}, nil
}

type GenericInputModel struct {
	*objects.ObjectModelImpl
}

func (o *GenericInputModel) Start() error {
	if err := o.ObjectModelImpl.Start(); err != nil {
		return errors.Wrap(err, "GenericInputModel.Start")
	}

	portID, err := o.GetProps().GetIntValue("address")
	if err != nil {
		return errors.Wrap(err, "GenericInputModel.Start")
	}

	err = o.Subscribe(
		"",
		"",
		messages.MessageTypeEvent,
		"",
		messages.TargetTypeObject,
		&portID,
		o.handler,
	)
	if err != nil {
		return errors.Wrap(err, "GenericInputModel.Start")
	}

	context.Logger.Debugf("GenericInputModel(%d) started", o.GetID())

	return nil
}

func (o *GenericInputModel) handler(msg messages.Message) ([]messages.Message, error) {
	context.Logger.Debugf("GenericInputModel(%d): handler()", o.GetID())

	var err error

	switch msg.GetName() {
	case "object.port.on_press":
		msg, err = generic_input.NewOnClickMessage("object_manager/object/event", o.GetID())
	case "object.port.on_double_click":
		msg, err = generic_input.NewOnDoubleClickMessage("object_manager/object/event", o.GetID())
	case "object.port.on_long_press":
		msg, err = generic_input.NewOnLongPressMessage("object_manager/object/event", o.GetID())
	default:
		return nil, nil
	}

	if err != nil {
		return nil, errors.Wrap(err, "GenericInputModel.handler")
	}

	return []messages.Message{msg}, nil
}

func (o *GenericInputModel) Shutdown() error {
	context.Logger.Debugf("GenericInputModel(%d) stopped", o.GetID())

	return nil
}
