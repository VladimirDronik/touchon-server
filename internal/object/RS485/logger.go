package RS485

import (
	"encoding/hex"

	"github.com/simonvetter/modbus"
	"github.com/sirupsen/logrus"
)

func NewLogger(client Client, logger *logrus.Logger) *Logger {
	return &Logger{
		client: client,
		logger: logger,
	}
}

type Logger struct {
	client Client
	logger *logrus.Logger
}

func (o *Logger) Open() error {
	err := o.client.Open()
	o.log("Open", nil, s(err))
	return err
}

func (o *Logger) Close() error {
	err := o.client.Close()
	o.log("Close", nil, s(err))
	return err
}

func (o *Logger) SetUnitId(id uint8) error {
	err := o.client.SetUnitId(id)
	o.log("SetUnitId", s(id), s(err))
	return err
}

func (o *Logger) SetEncoding(endianness modbus.Endianness, wordOrder modbus.WordOrder) error {
	err := o.client.SetEncoding(endianness, wordOrder)
	o.log("SetEncoding", s(endianness, wordOrder), s(err))
	return err
}

func (o *Logger) ReadCoils(addr uint16, quantity uint16) ([]bool, error) {
	v, err := o.client.ReadCoils(addr, quantity)
	o.log("ReadCoils", s(addr, quantity), s(v, err))
	return v, err
}

func (o *Logger) ReadCoil(addr uint16) (bool, error) {
	v, err := o.client.ReadCoil(addr)
	o.log("ReadCoil", s(addr), s(v, err))
	return v, err
}

func (o *Logger) ReadDiscreteInputs(addr uint16, quantity uint16) ([]bool, error) {
	v, err := o.client.ReadDiscreteInputs(addr, quantity)
	o.log("ReadDiscreteInputs", s(addr, quantity), s(v, err))
	return v, err
}

func (o *Logger) ReadDiscreteInput(addr uint16) (bool, error) {
	v, err := o.client.ReadDiscreteInput(addr)
	o.log("ReadDiscreteInput", s(addr), s(v, err))
	return v, err
}

func (o *Logger) ReadRegisters(addr uint16, quantity uint16, regType modbus.RegType) (values []uint16, err error) {
	v, err := o.client.ReadRegisters(addr, quantity, regType)
	o.log("ReadRegisters", s(addr, quantity, regType), s(v, err))
	return v, err
}

func (o *Logger) ReadRegister(addr uint16, regType modbus.RegType) (value uint16, err error) {
	v, err := o.client.ReadRegister(addr, regType)
	o.log("ReadRegister", s(addr, regType), s(v, err))
	return v, err
}

func (o *Logger) ReadUint32s(addr uint16, quantity uint16, regType modbus.RegType) (values []uint32, err error) {
	v, err := o.client.ReadUint32s(addr, quantity, regType)
	o.log("ReadUint32s", s(addr, quantity, regType), s(v, err))
	return v, err
}

func (o *Logger) ReadUint32(addr uint16, regType modbus.RegType) (value uint32, err error) {
	v, err := o.client.ReadUint32(addr, regType)
	o.log("ReadUint32", s(addr, regType), s(v, err))
	return v, err
}

func (o *Logger) ReadFloat32s(addr uint16, quantity uint16, regType modbus.RegType) (values []float32, err error) {
	v, err := o.client.ReadFloat32s(addr, quantity, regType)
	o.log("ReadFloat32s", s(addr, quantity, regType), s(v, err))
	return v, err
}

func (o *Logger) ReadFloat32(addr uint16, regType modbus.RegType) (value float32, err error) {
	v, err := o.client.ReadFloat32(addr, regType)
	o.log("ReadFloat32", s(addr, regType), s(v, err))
	return v, err
}

func (o *Logger) ReadUint64s(addr uint16, quantity uint16, regType modbus.RegType) (values []uint64, err error) {
	v, err := o.client.ReadUint64s(addr, quantity, regType)
	o.log("ReadUint64s", s(addr, quantity, regType), s(v, err))
	return v, err
}

func (o *Logger) ReadUint64(addr uint16, regType modbus.RegType) (value uint64, err error) {
	v, err := o.client.ReadUint64(addr, regType)
	o.log("ReadUint64", s(addr, regType), s(v, err))
	return v, err
}

func (o *Logger) ReadFloat64s(addr uint16, quantity uint16, regType modbus.RegType) (values []float64, err error) {
	v, err := o.client.ReadFloat64s(addr, quantity, regType)
	o.log("ReadFloat64s", s(addr, quantity, regType), s(v, err))
	return v, err
}

func (o *Logger) ReadFloat64(addr uint16, regType modbus.RegType) (value float64, err error) {
	v, err := o.client.ReadFloat64(addr, regType)
	o.log("ReadFloat64", s(addr, regType), s(v, err))
	return v, err
}

func (o *Logger) ReadBytes(addr uint16, quantity uint16, regType modbus.RegType) (values []byte, err error) {
	v, err := o.client.ReadBytes(addr, quantity, regType)
	o.log("ReadBytes", s(addr, quantity, regType), s(v, err))
	return v, err
}

func (o *Logger) ReadRawBytes(addr uint16, quantity uint16, regType modbus.RegType) (values []byte, err error) {
	v, err := o.client.ReadRawBytes(addr, quantity, regType)
	o.log("ReadRawBytes", s(addr, quantity, regType), s(v, err))
	return v, err
}

func (o *Logger) WriteCoil(addr uint16, value bool) (err error) {
	err = o.client.WriteCoil(addr, value)
	o.log("WriteCoil", s(addr, value), s(err))
	return err
}

func (o *Logger) WriteCoils(addr uint16, values []bool) (err error) {
	err = o.client.WriteCoils(addr, values)
	o.log("WriteCoils", s(addr, values), s(err))
	return err
}

func (o *Logger) WriteRegister(addr uint16, value uint16) (err error) {
	err = o.client.WriteRegister(addr, value)
	o.log("WriteRegister", s(addr, value), s(err))
	return err
}

func (o *Logger) WriteRegisters(addr uint16, values []uint16) (err error) {
	err = o.client.WriteRegisters(addr, values)
	o.log("WriteRegisters", s(addr, values), s(err))
	return err
}

func (o *Logger) WriteUint32s(addr uint16, values []uint32) (err error) {
	err = o.client.WriteUint32s(addr, values)
	o.log("WriteUint32s", s(addr, values), s(err))
	return err
}

func (o *Logger) WriteUint32(addr uint16, value uint32) (err error) {
	err = o.client.WriteUint32(addr, value)
	o.log("WriteUint32", s(addr, value), s(err))
	return err
}

func (o *Logger) WriteFloat32s(addr uint16, values []float32) (err error) {
	err = o.client.WriteFloat32s(addr, values)
	o.log("WriteFloat32s", s(addr, values), s(err))
	return err
}

func (o *Logger) WriteFloat32(addr uint16, value float32) (err error) {
	err = o.client.WriteFloat32(addr, value)
	o.log("WriteFloat32", s(addr, value), s(err))
	return err
}

func (o *Logger) WriteUint64s(addr uint16, values []uint64) (err error) {
	err = o.client.WriteUint64s(addr, values)
	o.log("WriteUint64s", s(addr, values), s(err))
	return err
}

func (o *Logger) WriteUint64(addr uint16, value uint64) (err error) {
	err = o.client.WriteUint64(addr, value)
	o.log("WriteUint64", s(addr, value), s(err))
	return err
}

func (o *Logger) WriteFloat64s(addr uint16, values []float64) (err error) {
	err = o.client.WriteFloat64s(addr, values)
	o.log("WriteFloat64s", s(addr, values), s(err))
	return err
}

func (o *Logger) WriteFloat64(addr uint16, value float64) (err error) {
	err = o.client.WriteFloat64(addr, value)
	o.log("WriteFloat64", s(addr, value), s(err))
	return err
}

func (o *Logger) WriteBytes(addr uint16, values []byte) (err error) {
	err = o.client.WriteBytes(addr, values)
	o.log("WriteBytes", s(addr, values), s(err))
	return err
}

func (o *Logger) WriteRawBytes(addr uint16, values []byte) (err error) {
	err = o.client.WriteRawBytes(addr, values)
	o.log("WriteRawBytes", s(addr, values), s(err))
	return err
}

// Методы для "сырой" работы с девайсом

func (o *Logger) WriteToDevice(data []byte) error {
	err := o.client.WriteToDevice(data)
	o.log("WriteToDevice", s(hex.EncodeToString(data)), s(err))
	return err
}

func (o *Logger) ReadFromDevice(bytesCount int) ([]byte, error) {
	v, err := o.client.ReadFromDevice(bytesCount)
	o.log("ReadFromDevice", s(bytesCount), s(v, err))
	return v, err
}

func (o *Logger) log(funcName string, args []interface{}, result []interface{}) {
	if o.logger.IsLevelEnabled(logrus.DebugLevel) {
		o.logger.Logf(o.logger.Level, "RS485.Client.%s(%v) -> %v", funcName, args, result)
	}
}

func s(args ...interface{}) []interface{} {
	return args
}
