// Шина Modbus

package Modbus

import (
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/simonvetter/modbus"
	"touchon-server/internal/context"
	"touchon-server/internal/model"
	"touchon-server/internal/objects"
	"touchon-server/lib/models"
	"touchon-server/lib/priority_queue"
)

const QueueCapabilities = 1000
const QueuePriorities = 10
const QueueMaxPriority = 1
const QueueMinPriority = 10

type Modbus interface {
	objects.Object
	DoAction(deviceAddr int, action Action, actionTries int, resultHandler ResultHandler, priority int) error
	GetDefaultTries() int
	GetClient() Client
	SetClient(client Client)
}

type Client interface {
	Open() (err error)
	Close() (err error)
	SetUnitId(id uint8) (err error)
	SetEncoding(endianness modbus.Endianness, wordOrder modbus.WordOrder) (err error)
	ReadCoils(addr uint16, quantity uint16) (values []bool, err error)
	ReadCoil(addr uint16) (value bool, err error)
	ReadDiscreteInputs(addr uint16, quantity uint16) (values []bool, err error)
	ReadDiscreteInput(addr uint16) (value bool, err error)
	ReadRegisters(addr uint16, quantity uint16, regType modbus.RegType) (values []uint16, err error)
	ReadRegister(addr uint16, regType modbus.RegType) (value uint16, err error)
	ReadUint32s(addr uint16, quantity uint16, regType modbus.RegType) (values []uint32, err error)
	ReadUint32(addr uint16, regType modbus.RegType) (value uint32, err error)
	ReadFloat32s(addr uint16, quantity uint16, regType modbus.RegType) (values []float32, err error)
	ReadFloat32(addr uint16, regType modbus.RegType) (value float32, err error)
	ReadUint64s(addr uint16, quantity uint16, regType modbus.RegType) (values []uint64, err error)
	ReadUint64(addr uint16, regType modbus.RegType) (value uint64, err error)
	ReadFloat64s(addr uint16, quantity uint16, regType modbus.RegType) (values []float64, err error)
	ReadFloat64(addr uint16, regType modbus.RegType) (value float64, err error)
	ReadBytes(addr uint16, quantity uint16, regType modbus.RegType) (values []byte, err error)
	ReadRawBytes(addr uint16, quantity uint16, regType modbus.RegType) (values []byte, err error)
	WriteCoil(addr uint16, value bool) (err error)
	WriteCoils(addr uint16, values []bool) (err error)
	WriteRegister(addr uint16, value uint16) (err error)
	WriteRegisters(addr uint16, values []uint16) (err error)
	WriteUint32s(addr uint16, values []uint32) (err error)
	WriteUint32(addr uint16, value uint32) (err error)
	WriteFloat32s(addr uint16, values []float32) (err error)
	WriteFloat32(addr uint16, value float32) (err error)
	WriteUint64s(addr uint16, values []uint64) (err error)
	WriteUint64(addr uint16, value uint64) (err error)
	WriteFloat64s(addr uint16, values []float64) (err error)
	WriteFloat64(addr uint16, value float64) (err error)
	WriteBytes(addr uint16, values []byte) (err error)
	WriteRawBytes(addr uint16, values []byte) (err error)
}

type Action func(client Client) (interface{}, error)
type ResultHandler func(result interface{}, err error)

type Task struct {
	DeviceAddr    int
	Action        Action
	Tries         int
	ResultHandler ResultHandler
}

func init() {
	_ = objects.Register(MakeModel)
}

func MakeModel() (objects.Object, error) {
	props := []*objects.Prop{
		{
			Code:        "connection_string",
			Name:        "Строка подключения",
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
			Code:        "speed",
			Name:        "Скорость",
			Description: "",
			Item: &models.Item{
				Type: models.DataTypeEnum,
				Values: map[string]string{
					"1200":   "1200",
					"2400":   "2400",
					"4800":   "4800",
					"9600":   "9600",
					"19200":  "19200",
					"38400":  "38400",
					"57600":  "57600",
					"115200": "115200",
					"128000": "128000",
					"256000": "256000",
				},
				DefaultValue: "19200",
			},
			Required: objects.NewRequired(true),
			Editable: objects.NewCondition(),
			Visible:  objects.NewCondition(),
		},
		{
			Code:        "data_bits",
			Name:        "Кол-во бит данных",
			Description: "",
			Item: &models.Item{
				Type:         models.DataTypeInt,
				DefaultValue: 8,
			},
			Required:   objects.NewRequired(true),
			Editable:   objects.NewCondition(),
			Visible:    objects.NewCondition(),
			CheckValue: objects.Above(0),
		},
		{
			Code:        "parity",
			Name:        "Контроль четности",
			Description: "",
			Item: &models.Item{
				Type: models.DataTypeEnum,
				Values: map[string]string{
					"0": "none",
					"1": "even",
					"2": "odd",
				},
				DefaultValue: "0",
			},
			Required: objects.NewRequired(true),
			Editable: objects.NewCondition(),
			Visible:  objects.NewCondition(),
		},
		{
			Code:        "stop_bits",
			Name:        "Стоповые биты",
			Description: "",
			Item: &models.Item{
				Type: models.DataTypeEnum,
				Values: map[string]string{
					"0": "0",
					"1": "1",
					"2": "2",
				},
				DefaultValue: "2",
			},
			Required: objects.NewRequired(true),
			Editable: objects.NewCondition(),
			Visible:  objects.NewCondition(),
		},
		{
			Code:        "timeout",
			Name:        "Таймаут",
			Description: "Время ожидания в сек",
			Item: &models.Item{
				Type:         models.DataTypeInt,
				DefaultValue: 3,
			},
			Required:   objects.NewRequired(true),
			Editable:   objects.NewCondition(),
			Visible:    objects.NewCondition(),
			CheckValue: objects.Between(1, 60),
		},
		{
			Code:        "tries",
			Name:        "Кол-во попыток выполнения операции",
			Description: "",
			Item: &models.Item{
				Type:         models.DataTypeInt,
				DefaultValue: 3,
			},
			Required:   objects.NewRequired(true),
			Editable:   objects.NewCondition(),
			Visible:    objects.NewCondition(),
			CheckValue: objects.Between(1, 10),
		},
	}

	impl, err := objects.NewObjectModelImpl(
		model.CategoryModbus,
		"modbus",
		false,
		"Шина Modbus",
		props,
		nil,
		nil,
		nil,
		[]string{"modbus"},
	)
	if err != nil {
		return nil, errors.Wrap(err, "Modbus.MakeModel")
	}

	queue, err := priority_queue.New[*Task](QueueCapabilities, QueuePriorities)
	if err != nil {
		return nil, errors.Wrap(err, "Modbus.MakeModel")
	}

	o := &ModbusImpl{
		Object: impl,
		queue:  queue,
		done:   make(chan struct{}),
	}

	return o, nil
}

type ModbusImpl struct {
	objects.Object
	client Client
	tries  int

	queue *priority_queue.PriorityQueue[*Task]
	done  chan struct{}
	wg    sync.WaitGroup
}

func (o *ModbusImpl) GetClient() Client {
	return o.client
}

func (o *ModbusImpl) SetClient(client Client) {
	o.client = client
}

func (o *ModbusImpl) Start() error {
	var err error

	if o.tries, err = o.GetProps().GetIntValue("tries"); err != nil {
		return errors.Wrap(err, "ModbusImpl.initClient")
	}

	if o.client == nil {
		if err := o.initClient(); err != nil {
			return errors.Wrap(err, "ModbusImpl.Start")
		}
	}

	// Запускаем обработку
	o.wg.Add(1)
	go o.actionsHandler()

	context.Logger.Debugf("Modbus(%d) started", o.GetID())

	return nil
}

func (o *ModbusImpl) initClient() error {
	if err := o.Object.Start(); err != nil {
		return errors.Wrap(err, "ModbusImpl.initClient")
	}

	connString, err := o.GetProps().GetStringValue("connection_string")
	if err != nil {
		return errors.Wrap(err, "ModbusImpl.initClient")
	}

	speed, err := o.GetProps().GetIntValue("speed")
	if err != nil {
		return errors.Wrap(err, "ModbusImpl.initClient")
	}

	dataBits, err := o.GetProps().GetIntValue("data_bits")
	if err != nil {
		return errors.Wrap(err, "ModbusImpl.initClient")
	}

	parity, err := o.GetProps().GetIntValue("parity")
	if err != nil {
		return errors.Wrap(err, "ModbusImpl.initClient")
	}

	stopBits, err := o.GetProps().GetIntValue("stop_bits")
	if err != nil {
		return errors.Wrap(err, "ModbusImpl.initClient")
	}

	timeout, err := o.GetProps().GetIntValue("timeout")
	if err != nil {
		return errors.Wrap(err, "ModbusImpl.initClient")
	}

	cfg := &modbus.ClientConfiguration{
		URL:      connString,
		Timeout:  time.Duration(timeout) * time.Second,
		Speed:    uint(speed),
		DataBits: uint(dataBits),
		Parity:   uint(parity),
		StopBits: uint(stopBits),
	}

	if o.client, err = modbus.NewClient(cfg); err != nil {
		return errors.Wrap(err, "ModbusImpl.initClient")
	}

	return nil
}

func (o *ModbusImpl) Shutdown() error {
	if err := o.Object.Shutdown(); err != nil {
		return errors.Wrap(err, "ModbusImpl.Shutdown")
	}

	close(o.done)

	// Ждем завершения дополнительных потоков
	o.wg.Wait()

	context.Logger.Debugf("Modbus(%d) stopped", o.GetID())

	return nil
}

func (o *ModbusImpl) DoAction(deviceAddr int, action Action, actionTries int, resultHandler ResultHandler, priority int) error {
	if actionTries < 1 || 10 < actionTries {
		return errors.Wrap(errors.Errorf("actionTries is bad"), "ModbusImpl.DoAction")
	}

	task := &Task{
		DeviceAddr:    deviceAddr,
		Action:        action,
		Tries:         actionTries,
		ResultHandler: resultHandler,
	}

	if err := o.queue.Push(task, priority); err != nil {
		return errors.Wrap(err, "ModbusImpl.DoAction")
	}

	return nil
}

func (o *ModbusImpl) actionsHandler() {
	defer o.wg.Done()

	for {
		// Если сервис завершает работу, выходим из цикла
		select {
		case <-o.done:
			return
		default:
		}

		// Читаем очередную задачу
		task, ok := o.queue.Pop()
		if !ok {
			// Если задач нет, засыпаем на время и переходим в начало
			time.Sleep(10 * time.Millisecond)
			continue
		}

		// Обрабатываем задачу
		result, err := func() (_ interface{}, e error) {
			if err := o.open(); err != nil {
				return nil, errors.Wrap(err, "ModbusImpl.actionsHandler")
			}

			defer func() {
				if err := o.client.Close(); err != nil && e == nil {
					e = errors.Wrap(err, "ModbusImpl.actionsHandler")
				}
			}()

			if err := o.client.SetUnitId(uint8(task.DeviceAddr)); err != nil {
				return nil, errors.Wrap(err, "ModbusImpl.actionsHandler")
			}

			var r interface{}
			var err error

			for i := 0; i < task.Tries; i++ {
				if r, err = task.Action(o.client); err == nil {
					return r, nil
				}
			}

			return nil, errors.Wrap(err, "ModbusImpl.actionsHandler")
		}()

		// Вызываем обработчик результата новой горутине
		go task.ResultHandler(result, err)
	}
}

func (o *ModbusImpl) open() error {
	var err error

	for i := 0; i < o.tries; i++ {
		if err = o.client.Open(); err == nil {
			return nil
		}
	}

	return errors.Wrap(err, "ModbusImpl.open")
}

func (o *ModbusImpl) GetDefaultTries() int {
	return o.tries
}
