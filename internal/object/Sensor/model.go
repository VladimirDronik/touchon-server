package Sensor

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	"touchon-server/internal/ws"

	"github.com/pkg/errors"
	"touchon-server/internal/g"
	"touchon-server/internal/model"
	"touchon-server/internal/object/SensorValue"
	"touchon-server/internal/objects"
	"touchon-server/internal/store"
	"touchon-server/lib/events"
	"touchon-server/lib/events/object/sensor"
	"touchon-server/lib/interfaces"
	"touchon-server/lib/messages"
	"touchon-server/lib/models"
)

var ErrGetValuesFuncNotSet = errors.New("getValues() not set")

type GetValuesFunc func(timeout time.Duration) (map[SensorValue.Type]float32, error)

const ValueUpdateAtFormat = "02.01.2006 15:04:05"

const timeout = 5

func init() {
	// Базовый тип датчика не регистрируем
	// _ = objects.Register(MakeModel)
}

func MakeModel(withChildren bool) (objects.Object, error) {
	props := []*objects.Prop{
		{
			Code:        "interface",
			Name:        "Интерфейс подключения",
			Description: "",
			Item: &models.Item{
				Type: models.DataTypeEnum,
				Values: map[string]string{
					"1W":       "1-Wire",
					"1WBUS":    "Шина 1-Wire",
					"I2C":      "I2C",
					"MODBUS":   "Modbus",
					"ADC":      "АЦП",
					"DISCRETE": "Дискретный",
					"MEGA-IN":  "IN",
				},
				DefaultValue: "I2C",
			},
			Required: objects.True(),
			// По умолчанию, менять заданный в модели интерфейс датчика нельзя
			Editable: objects.False(),
			Visible:  objects.True(),
		},
		{
			Code:        "address",
			Name:        "Адрес датчика на шине или id порта",
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
			Code:        "update_interval",
			Name:        "Интервал (10s, 1m etc)",
			Description: "Интервал получения значения датчика",
			Item: &models.Item{
				Type:         models.DataTypeString,
				DefaultValue: "5s",
			},
			Required:   objects.True(),
			Editable:   objects.True(),
			Visible:    objects.True(),
			CheckValue: objects.AboveOrEqual1(),
		},
	}

	onCheck, err := sensor.NewOnCheck(0, nil)
	if err != nil {
		return nil, errors.Wrap(err, "Sensor.MakeModel")
	}

	onAlarm, err := sensor.NewOnAlarm(0, "")
	if err != nil {
		return nil, errors.Wrap(err, "Sensor.MakeModel")
	}

	impl, err := objects.NewObjectModelImpl(
		model.CategorySensor,
		"",
		0,
		"",
		props,
		nil,
		[]interfaces.Event{onCheck, onAlarm},
		nil,
		[]string{"sensor"},
	)
	if err != nil {
		return nil, errors.Wrap(err, "Sensor.MakeModel")
	}

	obj := &SensorModel{
		ObjectModelImpl: impl,
	}

	check, err := objects.NewMethod("check", "Опрашивает датчик, обновляет показания датчика в БД", nil, obj.Check)
	if err != nil {
		return nil, errors.Wrap(err, "Sensor.MakeModel")
	}

	obj.GetMethods().Add(check)

	return obj, nil
}

type SensorModel struct {
	*objects.ObjectModelImpl
	getValues GetValuesFunc
}

func (o *SensorModel) SetGetValuesFunc(f GetValuesFunc) {
	o.getValues = f
}

func (o *SensorModel) ParseI2CAddress() (sdaPortObjectID, sclPortObjectID int, _ error) {
	sdaPortObjectID, sclPortObjectID = -1, -1

	addr, err := o.GetProps().GetStringValue("address")
	if err != nil {
		return 0, 0, errors.Wrap(err, "ParseI2CAddress")
	}

	errBadAddress := errors.Errorf("field 'address' has bad value %q", addr)

	portsSdaScl := strings.Split(addr, ";")
	switch len(portsSdaScl) {
	case 1:
		if sdaPortObjectID, err = strconv.Atoi(portsSdaScl[0]); err != nil {
			return 0, 0, errors.Wrap(err, "ParseI2CAddress")
		}

	case 2:
		if sdaPortObjectID, err = strconv.Atoi(portsSdaScl[0]); err != nil {
			return 0, 0, errors.Wrap(err, "ParseI2CAddress")
		}

		if sclPortObjectID, err = strconv.Atoi(portsSdaScl[1]); err != nil {
			return 0, 0, errors.Wrap(err, "ParseI2CAddress")
		}

	default:
		return 0, 0, errors.Wrap(errBadAddress, "ParseI2CAddress")
	}

	if sdaPortObjectID == 0 || sclPortObjectID == 0 {
		return 0, 0, errors.Wrap(errBadAddress, "ParseI2CAddress")
	}

	return
}

func (o *SensorModel) Check(map[string]interface{}) ([]interfaces.Message, error) {
	if o.getValues == nil {
		return nil, errors.Wrap(ErrGetValuesFuncNotSet, "SensorModel.Check")
	}

	g.Logger.Debugf("Sensor %s (%d) check", o.GetType(), o.GetID())

	values, err := o.getValues(time.Duration(timeout) * time.Second)
	if err != nil {
		return nil, errors.Wrap(err, "SensorModel.Check")
	}

	msgs := make([]interfaces.Message, 0, o.GetChildren().Len())
	for _, child := range o.GetChildren().GetAll() {
		msg, err := o.processChild(child, values)
		if err != nil {
			return nil, errors.Wrap(err, "SensorModel.Check")
		}

		if msg != nil {
			msgs = append(msgs, msg)
		}
	}

	vals := make(map[string]float32, len(values))
	for k, v := range values {
		vals[string(k)] = v
	}

	ws.I.Send("object", model.ObjectForWS{ID: o.GetID(), Value: vals})

	msg, err := sensor.NewOnCheck(o.GetID(), vals)
	if err != nil {
		return nil, errors.Wrap(err, "SensorModel.Check")
	}

	msgs = append(msgs, msg)

	if err := g.Msgs.Send(msgs...); err != nil {
		return nil, errors.Wrap(err, "SensorModel.Check")
	}

	return nil, nil
}

func (o *SensorModel) check() ([]interfaces.Message, error) {
	defer o.GetTimer().Reset()
	return o.Check(nil)
}

func (o *SensorModel) Start() error {
	if err := o.ObjectModelImpl.Start(); err != nil {
		return errors.Wrap(err, "SensorModel.Start")
	}

	// Если update_interval отсутствует (motion, presence etc), не запускаем таймер
	updateIntervalS, err := o.GetProps().GetStringValue("update_interval")
	if err != nil {
		g.Logger.Warn(errors.Wrap(err, "SensorModel.Start"))
		return nil
	}

	updateInterval, err := time.ParseDuration(updateIntervalS)
	if err != nil {
		return errors.Wrap(err, "SensorModel.Start")
	}

	_, err = o.Check(nil)
	if err != nil {
		g.Logger.Error(errors.Wrap(err, "SensorModel.Start: Check func"))
	}

	o.SetTimer(updateInterval, func() {
		// Игнорируем ошибку ErrGetValuesFuncNotSet при автоматическом обновлении датчика
		if _, err := o.check(); err != nil && !errors.Is(err, ErrGetValuesFuncNotSet) {
			g.Logger.Error(errors.Wrap(err, "SensorModel.Start"))
		}
	})
	o.GetTimer().Start()

	return nil
}

func (o *SensorModel) Shutdown() error {
	if err := o.ObjectModelImpl.Shutdown(); err != nil {
		return errors.Wrap(err, "SensorModel.Shutdown")
	}

	return nil
}

func (o *SensorModel) processChild(child objects.Object, values map[SensorValue.Type]float32) (interfaces.Message, error) {
	v, ok := values[child.GetType()]
	if !ok {
		return nil, errors.Wrap(errors.Errorf("value for child object %q not found", child.GetName()), "processChild")
	}

	if err := child.GetProps().Set("value", v); err != nil {
		return nil, errors.Wrap(err, "processChild")
	}

	if err := child.GetProps().Set("value_updated_at", time.Now().Format(ValueUpdateAtFormat)); err != nil {
		return nil, errors.Wrap(err, "processChild")
	}

	valueObj, ok := child.(*SensorValue.SensorValueModel)
	if !ok {
		return nil, errors.Wrap(errors.Errorf("child object is not SensorValueModel type (%T)", child), "processChild")
	}

	if err := valueObj.Validate(); err != nil {
		if !errors.Is(err, SensorValue.ErrSensorAlarmValue) && !errors.Is(err, SensorValue.ErrSensorErrorValue) {
			return nil, errors.Wrap(err, "processChild")
		}

		text, err := MakeAlarmMessage(valueObj)
		if err != nil {
			return nil, errors.Wrap(err, "processChild")
		}

		text = fmt.Sprintf("Датчик %q (ID:%d): %s", o.GetName(), o.GetID(), text)
		var msg interfaces.Message

		if errors.Is(err, SensorValue.ErrSensorAlarmValue) {
			msg, err = sensor.NewOnAlarm(o.GetID(), text)
		} else {
			msg, err = events.NewOnError(interfaces.TargetTypeObject, o.GetID(), text)
		}
		if err != nil {
			return nil, errors.Wrap(err, "processChild")
		}

		return msg, nil
	}

	if err := o.SaveSensorValue(child); err != nil {
		return nil, errors.Wrap(err, "processChild")
	}

	return nil, nil
}

// SaveSensorValue добавление данных значения датчика в БД
func (o *SensorModel) SaveSensorValue(valueObj objects.Object) error {
	value, err := valueObj.GetProps().GetFloatValue("value")
	if err != nil {
		return errors.Wrap(err, "SaveSensorValue")
	}

	valueUpdatedAt, err := valueObj.GetProps().GetStringValue("value_updated_at")
	if err != nil {
		return errors.Wrap(err, "SaveSensorValue")
	}

	props := map[string]string{
		"value":            fmt.Sprintf("%.2f", value),
		"value_updated_at": valueUpdatedAt,
	}

	// Добавление значения датчика
	if err := store.I.ObjectRepository().SetProps(valueObj.GetID(), props); err != nil {
		return errors.Wrap(err, "SaveSensorValue")
	}

	return nil
}

func (o *SensorModel) DeleteChildren() error {
	for _, child := range o.GetChildren().GetAll() {
		if err := child.DeleteChildren(); err != nil {
			return errors.Wrap(err, "DeleteChildren")
		}

		// Удаляем только значения датчика
		if child.GetCategory() != model.CategorySensorValue {
			continue
		}

		if err := store.I.ObjectRepository().DelObject(child.GetID()); err != nil {
			return errors.Wrap(err, "DeleteChildren")
		}
	}

	return nil
}

func (o *SensorModel) GetState() (interfaces.Message, error) {
	msg, err := messages.NewEvent("on_get_state", interfaces.TargetTypeObject, o.GetID())
	if err != nil {
		return nil, errors.Wrap(err, "SensorModel.GetState")
	}

	for _, child := range o.GetChildren().GetAll() {
		v, err := child.GetProps().GetFloatValue("value")
		if err != nil {
			return nil, errors.Wrap(err, "SensorModel.GetState")
		}

		msg.SetValue(child.GetType(), v)
	}

	return msg, nil
}
