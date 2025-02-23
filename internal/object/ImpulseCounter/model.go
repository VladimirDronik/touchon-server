package ImpulseCounter

import (
	"github.com/pkg/errors"
	"touchon-server/internal/g"
	helpersObj "touchon-server/internal/helpers"
	"touchon-server/internal/model"
	"touchon-server/internal/objects"
	"touchon-server/lib/events/object/impulse_counter"
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
			Name:        "Адрес устройства на шине или ID порта",
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
			Name:        "Интерфейс, по которому подключено устройство",
			Description: "",
			Item: &models.Item{
				Type:         models.DataTypeString,
				DefaultValue: "",
			},
			Required: objects.True(),
			Editable: objects.True(),
			Visible:  objects.True(),
		},
		{
			Code:        "threshold",
			Name:        "Порог генерации события счетчика",
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
			Code:        "value",
			Name:        "Значение счетчика, хранящееся в БД",
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
			Code:        "server_request",
			Name:        "Опрос счетчика сервером",
			Description: "",
			Item: &models.Item{
				Type:         models.DataTypeBool,
				DefaultValue: true,
			},
			Required: objects.True(),
			Editable: objects.True(),
			Visible:  objects.True(),
		},
		{
			Code:        "polling_interval",
			Name:        "Интервал опроса счетчика сервером, с",
			Description: "",
			Item: &models.Item{
				Type:         models.DataTypeInt,
				DefaultValue: 5,
			},
			Required: objects.True(),
			Editable: objects.True(),
			Visible:  objects.True(),
		},
		{
			Code:        "reset_value",
			Name:        "Значение количества импульсов для обнуления счетчика",
			Description: "",
			Item: &models.Item{
				Type:         models.DataTypeInt,
				DefaultValue: 100,
			},
			Required: objects.True(),
			Editable: objects.True(),
			Visible:  objects.True(),
		},
		{
			Code:        "unit",
			Name:        "Единицы измерения счетчика",
			Description: "",
			Item: &models.Item{
				Type:         models.DataTypeString,
				DefaultValue: "",
			},
			Required: objects.True(),
			Editable: objects.True(),
			Visible:  objects.True(),
		},
		{
			Code:        "math_transformation",
			Name:        "Формула математического преобразования импульсов",
			Description: "",
			Item: &models.Item{
				Type:         models.DataTypeString,
				DefaultValue: "",
			},
			Required: objects.True(),
			Editable: objects.True(),
			Visible:  objects.True(),
		},
		{
			Code:        "write_graph",
			Name:        "Ведение графика для счетчика",
			Description: "",
			Item: &models.Item{
				Type:         models.DataTypeBool,
				DefaultValue: true,
			},
			Required: objects.True(),
			Editable: objects.True(),
			Visible:  objects.True(),
		},
		{
			Code:        "remove_jitter",
			Name:        "Убрать дребезг контакта",
			Description: "",
			Item: &models.Item{
				Type:         models.DataTypeBool,
				DefaultValue: false,
			},
			Required: objects.True(),
			Editable: objects.True(),
			Visible:  objects.True(),
		},
	}

	//events
	onThreshold, err := impulse_counter.NewOnThreshold(0, 0)
	if err != nil {
		return nil, errors.Wrap(err, "ImpulseCounter.MakeModel")
	}

	onCheck, err := impulse_counter.NewOnCheck(0, 0)
	if err != nil {
		return nil, errors.Wrap(err, "ImpulseCounter.MakeModel")
	}

	//model
	impl, err := objects.NewObjectModelImpl(
		model.CategoryGenericInput,
		"impulse_counter",
		false,
		"Счетчик импульсов",
		props,
		nil,
		[]interfaces.Event{onThreshold, onCheck},
		nil,
		[]string{"счетчик", "counter", "импульс"},
	)
	if err != nil {
		return nil, errors.Wrap(err, "ImpulseCounter.MakeModel")
	}

	o := &ImpulseCounter{ObjectModelImpl: impl}

	//methods
	check, err := objects.NewMethod("check", "Получение количества импульсов", nil, o.Check)
	if err != nil {
		return nil, errors.Wrap(err, "PortMegaD.MakeModel")
	}

	o.GetMethods().Add(check)

	return o, nil
}

type ImpulseCounter struct {
	*objects.ObjectModelImpl
}

func (o *ImpulseCounter) Start() error {
	if err := o.ObjectModelImpl.Start(); err != nil {
		return errors.Wrap(err, "ImpulseCounterModel.Start")
	}

	portID, err := o.GetProps().GetIntValue("address")
	if err != nil {
		return errors.Wrap(err, "ImpulseCounterModel.Start")
	}

	err = o.Subscribe(
		interfaces.MessageTypeEvent,
		"",
		interfaces.TargetTypeObject,
		&portID,
		o.handler,
	)
	if err != nil {
		return errors.Wrap(err, "ImpulseCounterModel.Start")
	}

	g.Logger.Debugf("ImpulseCounterModel(%d) started", o.GetID())

	return nil
}

func (o *ImpulseCounter) handler(svc interfaces.MessageSender, msg interfaces.Message) {
	var err error
	status := model.StatusAvailable

	switch msg.GetName() {
	case "object.port.on_press":
		countImpulse, err := msg.GetIntValue("count_impulse")
		if err != nil {
			g.Logger.Error(errors.Wrap(err, "ImpulseCounterModel.handler"))
			return
		}

		//логика сбрасывания значений счетчика, с учетом порогового значения
		resetValue, err := o.GetProps().GetIntValue("reset_value")
		if err != nil {
			g.Logger.Error(errors.Wrap(err, "ImpulseCounterModel.handler"))
			return
		}
		err = o.resetTo(65535 - resetValue)
		if err != nil {
			g.Logger.Error(errors.Wrap(err, "ImpulseCounterModel.handler: reset counter"))
			return
		}
		err = o.saveImpulses(countImpulse)
		if err != nil {
			g.Logger.Error(errors.Wrap(err, "ImpulseCounterModel.handler: save impulse data to DB failed"))
			return
		}

		msg, err = impulse_counter.NewOnThreshold(o.GetID(), countImpulse)
	default:
		return
	}

	if err != nil {
		status = model.StatusUnavailable
		g.Logger.Error(errors.Wrap(err, "ImpulseCounterModel.handler"))
		return
	}

	helpersObj.SaveAndSendStatus(o, status)
	if err := svc.Send(msg); err != nil {
		g.Logger.Error(errors.Wrap(err, "ImpulseCounterModel.handler"))
		return
	}
}

func (o *ImpulseCounter) Shutdown() error {
	g.Logger.Debugf("ImpulseCounterModel(%d) stopped", o.GetID())

	return nil
}
