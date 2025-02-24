package wb_mrm2_mini

import (
	"strconv"

	"github.com/pkg/errors"
	"touchon-server/internal/g"
	"touchon-server/internal/object/Modbus"
	"touchon-server/internal/object/Modbus/ModbusDevice"
	"touchon-server/internal/objects"
	"touchon-server/internal/scripts"
	"touchon-server/lib/events/object/wiren_board/wb_mrm2_mini"
	"touchon-server/lib/helpers"
	"touchon-server/lib/interfaces"
	"touchon-server/lib/models"
)

func init() {
	_ = objects.Register(MakeModel)
}

func MakeModel() (objects.Object, error) {
	baseObj, err := ModbusDevice.MakeModel()
	if err != nil {
		return nil, errors.Wrap(err, "wb_mrm2_mini.MakeModel")
	}

	obj := &DeviceModel{
		outputCount: 2,
		coilAddresses: map[int]uint16{
			1: 0x0000,
			2: 0x0001,
		},
	}

	obj.ModbusDevice = baseObj.(*ModbusDevice.ModbusDeviceImpl)

	obj.SetType("wb_mrm2_mini")
	obj.SetName("WB-MRM2-mini Двухканальный модуль реле")
	obj.SetTags("wb_mrm2_mini")

	// Добавляем свои события
	onCheck, err := wb_mrm2_mini.NewOnCheck(0, nil)
	if err != nil {
		return nil, errors.Wrap(err, "wb_mrm2_mini.MakeModel")
	}

	if err := obj.GetEvents().Add(onCheck); err != nil {
		return nil, errors.Wrap(err, "wb_mrm2_mini.MakeModel")
	}

	getOutputsState, err := objects.NewMethod("get_outputs_state", "Получение состояния всех выходов", nil, obj.GetOutputsState)
	if err != nil {
		return nil, errors.Wrap(err, "wb_mrm2_mini.MakeModel")
	}

	params := []*scripts.Param{
		{
			Code:        "k",
			Name:        "Номер выхода",
			Description: "Нумерация начинается с 1",
			Item: &models.Item{
				Type:         models.DataTypeInt,
				DefaultValue: 1,
			},
		},
	}

	getOutputState, err := objects.NewMethod("get_output_state", "Получение состояния выхода", params, obj.GetOutputState)
	if err != nil {
		return nil, errors.Wrap(err, "wb_mrm2_mini.MakeModel")
	}

	params = []*scripts.Param{
		{
			Code:        "k",
			Name:        "Номер выхода",
			Description: "Нумерация начинается с 1",
			Item: &models.Item{
				Type:         models.DataTypeInt,
				DefaultValue: 1,
			},
		},
		{
			Code:        "state",
			Name:        "Состояние выхода",
			Description: "",
			Item: &models.Item{
				Type:         models.DataTypeBool,
				DefaultValue: false,
			},
		},
	}

	setOutputState, err := objects.NewMethod("set_output_state", "Выставление состояния выхода", params, obj.SetOutputState)
	if err != nil {
		return nil, errors.Wrap(err, "wb_mrm2_mini.MakeModel")
	}

	obj.GetMethods().Add(getOutputsState, getOutputState, setOutputState)

	return obj, nil
}

type DeviceModel struct {
	ModbusDevice.ModbusDevice
	outputCount   int
	coilAddresses map[int]uint16
	unitID        int
}

func (o *DeviceModel) Start() error {
	if err := o.ModbusDevice.Start(); err != nil {
		return errors.Wrap(err, "wb_mrm2_mini.DeviceModel.Start")
	}

	address, err := o.GetProps().GetIntValue("address")
	if err != nil {
		return errors.Wrap(err, "wb_mrm2_mini.DeviceModel.Start")
	}
	o.unitID = address

	g.Logger.Debugf("WB-MRM2-mini(%d) started", o.GetID())

	return nil
}

func (o *DeviceModel) Shutdown() error {
	if err := o.ModbusDevice.Shutdown(); err != nil {
		return errors.Wrap(err, "wb_mrm2_mini.DeviceModel.Shutdown")
	}

	g.Logger.Debugf("WB-MRM2-mini(%d) stopped", o.GetID())

	return nil
}

func (o *DeviceModel) GetOutputsState(map[string]interface{}) ([]interfaces.Message, error) {
	action := func(client Modbus.Client) (interface{}, error) {
		return client.ReadCoils(0x0000, uint16(o.outputCount))
	}

	resultHandler := func(r interface{}, err error) {
		if err != nil {
			g.Logger.Error(errors.Wrap(err, "wb_mrm2_mini.DeviceModel.GetOutputsState"))
			return
		}

		states, ok := r.([]bool)
		if !ok {
			g.Logger.Error(errors.Wrap(errors.Errorf("ModbusDeviceImpl.DoAction returned bad value"), "wb_mrm2_mini.DeviceModel.GetOutputsState"))
			return
		}

		args := make(map[string]bool, len(states))
		for i, v := range states {
			args["k"+strconv.Itoa(i+1)] = v
		}

		msg, err := wb_mrm2_mini.NewOnCheck(o.GetID(), args)
		if err != nil {
			g.Logger.Error(errors.Wrap(err, "wb_mrm2_mini.DeviceModel.GetOutputsState"))
			return
		}

		if err := g.Msgs.Send(msg); err != nil {
			g.Logger.Error(errors.Wrap(err, "wb_mrm2_mini.DeviceModel.GetOutputsState"))
		}
	}

	if err := o.ModbusDevice.DoAction(o.unitID, action, o.GetDefaultTries(), resultHandler, Modbus.QueueMinPriority); err != nil {
		return nil, errors.Wrap(err, "wb_mrm2_mini.DeviceModel.GetOutputsState")
	}

	return nil, nil
}

func (o *DeviceModel) GetOutputState(args map[string]interface{}) ([]interfaces.Message, error) {
	outputNumber, err := helpers.GetNumber(args["k"])
	if err != nil {
		return nil, errors.Wrap(err, "wb_mrm2_mini.DeviceModel.GetOutputState(k)")
	}

	address, ok := o.coilAddresses[outputNumber]
	if !ok {
		return nil, errors.Wrap(errors.Errorf("address for output %d not found", outputNumber), "wb_mrm2_mini.DeviceModel.GetOutputState")
	}

	action := func(client Modbus.Client) (interface{}, error) {
		return client.ReadCoil(address)
	}

	resultHandler := func(r interface{}, err error) {
		if err != nil {
			g.Logger.Error(errors.Wrap(err, "wb_mrm2_mini.DeviceModel.GetOutputState"))
			return
		}

		state, ok := r.(bool)
		if !ok {
			g.Logger.Error(errors.Wrap(errors.Errorf("ModbusDeviceImpl.DoAction returned bad value"), "wb_mrm2_mini.DeviceModel.GetOutputState"))
			return
		}

		payload := map[string]bool{
			"k" + strconv.Itoa(outputNumber): state,
		}

		msg, err := wb_mrm2_mini.NewOnCheck(o.GetID(), payload)
		if err != nil {
			g.Logger.Error(errors.Wrap(err, "wb_mrm2_mini.DeviceModel.GetOutputState"))
			return
		}

		if err := g.Msgs.Send(msg); err != nil {
			g.Logger.Error(errors.Wrap(err, "wb_mrm2_mini.DeviceModel.GetOutputState"))
		}
	}

	if err := o.ModbusDevice.DoAction(o.unitID, action, o.GetDefaultTries(), resultHandler, Modbus.QueueMinPriority); err != nil {
		return nil, errors.Wrap(err, "wb_mrm2_mini.DeviceModel.GetOutputState")
	}

	return nil, nil
}

func (o *DeviceModel) SetOutputState(args map[string]interface{}) ([]interfaces.Message, error) {
	outputNumber, err := helpers.GetNumber(args["k"])
	if err != nil {
		return nil, errors.Wrap(err, "wb_mrm2_mini.DeviceModel.SetOutputState(k)")
	}

	state, ok := args["state"].(bool)
	if !ok {
		return nil, errors.Wrap(errors.Errorf("state is bad"), "wb_mrm2_mini.DeviceModel.SetOutputState")
	}

	address, ok := o.coilAddresses[outputNumber]
	if !ok {
		return nil, errors.Wrap(errors.Errorf("address for output %d not found", outputNumber), "wb_mrm2_mini.DeviceModel.SetOutputState")
	}

	action := func(client Modbus.Client) (interface{}, error) {
		return nil, client.WriteCoil(address, state)
	}

	resultHandler := func(r interface{}, err error) {
		if err != nil {
			g.Logger.Error(errors.Wrap(err, "wb_mrm2_mini.DeviceModel.SetOutputState"))
			return
		}
	}

	if err = o.ModbusDevice.DoAction(o.unitID, action, o.GetDefaultTries(), resultHandler, Modbus.QueueMinPriority); err != nil {
		return nil, errors.Wrap(err, "wb_mrm2_mini.DeviceModel.SetOutputState")
	}

	return nil, nil
}
