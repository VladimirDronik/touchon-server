package PortMegaD

import (
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
	"touchon-server/internal/model"
	"touchon-server/internal/objects"
	"touchon-server/internal/store"
	"touchon-server/lib/interfaces"
)

// ResCommand прием команды от megaD
func (o *PortModel) ResCommand(controllerID, portNumber, extPortNumber, clickCount, holdRelease, value, countImpulse string) ([]interfaces.Message, error) {
	objectID, err := store.I.PortRepository().GetPortObjectID(controllerID, portNumber)
	if err != nil {
		return nil, errors.Wrap(err, "ResCommand")
	}

	obj, err := objects.LoadPort(objectID, model.ChildTypeNobody)
	if err != nil {
		return nil, errors.Wrap(err, "ResCommand")
	}

	port, ok := obj.(*PortModel)
	if !ok {
		err := errors.New("MakeModel returns not PortModel")
		return nil, errors.Wrap(err, "ResCommand")
	}

	// Если сработал порт, который находится на модуле расширения
	if extPortNumber != "" {
		// TODO: описать логику обработки сработавшего порта
		return nil, errors.Wrap(errors.New("Поддержка расширителя не реализована"), "ResCommand")
	}

	msgs := make([]interfaces.Message, 0, 10)

	switch {
	case clickCount == "2":
		// вызываем событие двойного нажатия
		msgs = append(msgs, port.OnDoubleClick())

	case holdRelease != "":
		// вызываем событие длительного нажатия
		switch holdRelease {
		case "2":
			msgs = append(msgs, port.OnLongPress())
		case "1":
			msgs = append(msgs, port.OnRelease())
		}

	case value != "":
		mode, err := port.GetProps().GetStringValue("mode")
		if err != nil {
			return nil, errors.Wrap(err, "ResCommand")
		}

		if mode != "PWM" {
			status := "off"

			// Если контроллер шлет сообщение, что сработал выход, то в данном случае просто отвечаем ему кодом 200.
			if value == "1" {
				status = "on"
			}

			if err := port.SetState(status); err != nil {
				// TODO: сделать обработку ошибки при смене статуса объекта
			}

			msgs = append(msgs, port.OnChangeState(status))
		}

		return msgs, nil

	default:
		cnt, err := strconv.Atoi(countImpulse)
		if err != nil {
			return nil, errors.Wrap(err, "ResCommand: count impulse is fall")
		}
		// вызываем событие одиночного нажатия
		msgs = append(msgs, port.OnPress(cnt))
	}

	return msgs, nil
}

// sendCommand отправка команды на контроллер
func (o *PortModel) sendCommand(contrAddr string, portNumber int, portOption bool, command string, params map[string]string, timeout time.Duration) (int, []byte, error) {
	args := url.Values{}

	if portOption {
		args.Add("pt", strconv.Itoa(portNumber))
	}

	if command != "" {
		args.Add("cmd", command)
	}

	for k, v := range params {
		args.Add(k, v)
	}

	u := fmt.Sprintf("http://%s/sec/?%s", contrAddr, args.Encode())

	status, body, err := fasthttp.GetTimeout(nil, u, timeout)
	//g.Logger.Debugf("PortModel.sendCommand: GET %s Status=%d Body=%q Err=%v", u, status, body, err)
	if err != nil {
		return 0, nil, errors.Wrap(err, "sendCommand")
	}

	return status, body, err
}
