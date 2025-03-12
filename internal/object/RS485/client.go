package RS485

import (
	"net/url"

	"github.com/pkg/errors"
	"github.com/simonvetter/modbus"
	"touchon-server/internal/g"
)

var ErrBadUrlScheme = errors.New("bad url scheme")

func newClient(cfg *modbus.ClientConfiguration) (Client, error) {
	u, err := url.Parse(cfg.URL)
	if err != nil {
		return nil, errors.Wrap(err, "RS485.newClient")
	}

	mbClient, err := modbus.NewClient(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "RS485.newClient")
	}

	o := &ClientImpl{ModbusClient: mbClient}

	if u.Scheme == "rtu" { // rtuovertcp ??
		o.devicePath = u.Path
	}

	return NewLogger(o, g.Logger), nil
}

// ClientImpl оборачивает modbus.ModbusClient для реализации собственных функций сырой записи/чтения в устройство.
type ClientImpl struct {
	*modbus.ModbusClient
	devicePath string
}

func (o *ClientImpl) WriteToDevice(data []byte) error {
	if o.devicePath == "" {
		return errors.Wrap(ErrBadUrlScheme, "WriteToDevice")
	}

	// TODO...

	return nil
}

func (o *ClientImpl) ReadFromDevice(bytesCount int) ([]byte, error) {
	if o.devicePath == "" {
		return nil, errors.Wrap(ErrBadUrlScheme, "ReadFromDevice")
	}

	// TODO...

	return nil, nil
}
