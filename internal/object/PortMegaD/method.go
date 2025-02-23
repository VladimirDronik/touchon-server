package PortMegaD

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	"touchon-server/internal/g"
	"touchon-server/internal/model"
	"touchon-server/internal/store"
	"touchon-server/lib/interfaces"
	"touchon-server/lib/messages"
)

type method struct {
	port *PortModel
}

// Структура, которая будет отправлена в mqtt топик
type Content struct {
	Address string
	NumPort int
	Status  string
	Method  string
}

var typePort = map[string]string{
	"nc":   "255",
	"in":   "0",
	"out":  "1",
	"adc":  "2",
	"dsen": "3",
	"i2c":  "4",
}

var modePort = map[string]string{
	"":      "",
	"sw":    "m=0",
	"1w":    "d=3",
	"1wbus": "d=5",
	"sda":   "m=1",
	"scl":   "m=2",
	"pr":    "m=1",
}

// Init инициализация объекта
func (o *PortModel) Init(storeObj *model.StoreObject, childType model.ChildType) error {
	if err := o.ObjectModelImpl.Init(storeObj, childType); err != nil {
		return errors.Wrap(err, "PortModel.Init")
	}

	parentID := o.GetParentID()
	if parentID == nil {
		return errors.Wrap(errors.Errorf("parent_id of %d is nil", o.GetID()), "PortModel.Init")
	}

	// Получаем IP адрес контроллера
	contrIP, err := store.I.ObjectRepository().GetProp(*parentID, "address")
	if err != nil {
		return errors.Wrap(err, "PortModel.Init")
	}
	o.contrAddr = contrIP

	// Получаем номер порта
	o.portNumber, err = o.GetProps().GetIntValue("number")
	if err != nil {
		return errors.Wrap(err, "PortModel.Init")
	}

	return nil
}

func (o *PortModel) GetContrAddr() string {
	return o.contrAddr
}

func (o *PortModel) GetParentPortObjectID() int {
	return o.parentPortObjectID
}

func (o *PortModel) GetPortNumber() int {
	return o.portNumber
}

// On Включаем порт
func (o *PortModel) On(args map[string]interface{}) ([]interfaces.Message, error) {
	command := fmt.Sprintf("%d:1", o.GetPortNumber())
	return o.setPortStatus(false, command, nil)
}

// Off Выключаем порт
func (o *PortModel) Off(args map[string]interface{}) ([]interfaces.Message, error) {
	command := fmt.Sprintf("%d:0", o.GetPortNumber())
	return o.setPortStatus(false, command, nil)
}

func (o *PortModel) Toggle(args map[string]interface{}) ([]interfaces.Message, error) {
	command := fmt.Sprintf("%d:2", o.GetPortNumber())
	return o.setPortStatus(false, command, nil)
}

// SetPWM установка ШИМ для порта
func (o *PortModel) SetPWM(args map[string]interface{}) (msgs []interfaces.Message, e error) {
	value := 0
	smooth := 0

	if args["value"] != nil {
		v, _ := args["value"].(float32)
		value = int(v)
	}
	if args["smooth"] != nil {
		v, _ := args["smooth"].(float32)
		smooth = int(v)
	}

	defer func() {
		msg, err := messages.NewCommand("check", interfaces.TargetTypeObject, o.GetID(), nil)
		if err != nil {
			e = errors.Wrap(err, "SetPWM")
		}
		//msg.SetDelay(smooth*1000 + 20)
		msgs = append(msgs, msg)
	}()

	// Если порт не является портом расширения
	if o.GetParentPortObjectID() == 0 {
		msgs, err := o.setPortStatus(true, "", map[string]string{"pwm": strconv.Itoa(value)})
		if err != nil {
			return nil, errors.Wrap(err, "SetPWM")
		}

		return msgs, nil
	}

	// Если порт является портом расширения
	// Ищем номер родительского порта по id родительского объекта
	parentPortNumber, err := store.I.ObjectRepository().GetProp(o.GetParentPortObjectID(), "number")
	if err != nil {
		return nil, errors.Wrap(err, "SetPWM")
	}

	command := fmt.Sprintf("%se%d:%d", parentPortNumber, o.GetPortNumber(), value)
	params := map[string]string{"cnt": strconv.Itoa(smooth)}

	msgs, err = o.setPortStatus(false, command, params)
	if err != nil {
		return nil, errors.Wrap(err, "SetPWM")
	}

	return msgs, nil
}

func (o *PortModel) Impulse(args map[string]interface{}) ([]interfaces.Message, error) {
	value := 0

	if args["value"] != nil {
		v, _ := args["value"].(float32)
		value = int(v)
	}

	command := fmt.Sprintf("%[1]d:2;o%d;%[1]d:2", o.GetPortNumber(), value/100)

	msgs, err := o.setPortStatus(false, command, nil)
	if err != nil {
		return nil, errors.Wrap(err, "Impulse")
	}

	time.Sleep(time.Duration(value) * time.Millisecond)

	checkMsgs, err := o.Check(nil)
	if err != nil {
		return nil, errors.Wrap(err, "Impulse")
	}

	msgs = append(msgs, checkMsgs...)

	return msgs, nil
}

// Check получение статуса порта
func (o *PortModel) Check(args map[string]interface{}) ([]interfaces.Message, error) {
	newState, err := o.GetPortState("get", nil, time.Duration(5)*time.Second)
	if err != nil {
		// сообщаем, что порт не поменял свой статус и отдаем текущее состояние порта
		return nil, errors.Wrap(err, "Check")
	}

	if string(o.GetStatus()) == newState { // если статус порта пришел такой же, как был уже в БД
		msg := o.OnCheck(string(o.GetStatus()))
		return []interfaces.Message{msg}, nil
	}

	switch newState {
	case "ON":
		o.SetStatus(model.StatusOn)
	case "OFF":
		o.SetStatus(model.StatusOff)
	}

	// заносим статус порта в БД
	go func() {
		if err := store.I.ObjectRepository().SetObjectStatus(o.GetID(), newState); err != nil {
			g.Logger.Error(errors.Wrap(err, "PortModel.Check"))
		}
	}()

	msg := o.OnChangeState(newState)
	return []interfaces.Message{msg}, nil
}

// setPortStatus установка статуса для порта
func (o *PortModel) setPortStatus(portOption bool, command string, params map[string]string) ([]interfaces.Message, error) {
	code, _, err := o.sendCommand(o.GetContrAddr(), o.GetPortNumber(), portOption, command, params, time.Duration(5)*time.Second)
	if err != nil {
		return nil, errors.Wrap(err, "setPortStatus")
	}

	// Если код вернулся ошибочный
	if code > 299 {
		// сообщаем, что порт не поменял свой статус и отдаем текущее состояние порта
		return nil, errors.Wrap(errors.New("device not available"), "setPortStatus")
	} else {
		// проверяем статус порта, что он действительно поменялся и генерим событие onChange
		msgs, err := o.Check(nil)
		if err != nil {
			return nil, errors.Wrap(err, "setPortStatus")
		}

		return msgs, nil
	}
}

// SetState устанавливает статус для объекта без смены состояния самого порта на устройстве
// (используется для портов выхода, например)
func (o *PortModel) SetState(status string) error {
	if err := store.I.ObjectRepository().SetObjectStatus(o.GetID(), status); err != nil {
		return errors.Wrap(err, "SetState")
	}

	return nil
}

// SetTypeMode Установка типа и режима порта
func (o *PortModel) SetTypeMode(typePt string, modePt string, title string, extParams map[string]string) error {
	params := make(map[string]string)

	command := ""
	params = extParams
	params["pty"] = typePort[typePt]
	params["pn"] = strconv.Itoa(o.GetPortNumber())

	code, _, err := o.sendCommand(o.GetContrAddr(), o.GetPortNumber(), false, command, params, time.Duration(5)*time.Second)
	if err != nil || code > 299 {
		return errors.Wrap(err, "SetTypeMode")
	}

	params["emt"] = title
	if modePort[modePt] != "" {
		modeParam := strings.Split(modePort[modePt], "=")
		params[modeParam[0]] = modeParam[1]
	}

	time.Sleep(300 * time.Millisecond)
	//Второй раз отправляем на контроллер команду, потому что сразу он не может поменять и тип порта и режим
	code, _, err = o.sendCommand(o.GetContrAddr(), o.GetPortNumber(), false, command, params, time.Duration(5)*time.Second)
	if err != nil || code > 299 {
		return errors.Wrap(err, "SetTypeMode")
	}

	//Если получилось поменять тип и режим порта, то в БД тоже меняем
	if err := store.I.ObjectRepository().SetProp(o.GetID(), "type", typePt); err != nil {
		return errors.Wrap(err, "SetTypeMode")
	}

	if err := store.I.ObjectRepository().SetProp(o.GetID(), "mode", modePt); err != nil {
		return errors.Wrap(err, "SetTypeMode")
	}

	return nil
}

func (o *PortModel) SetPortParams(params map[string]string) (int, error) {
	code, _, err := o.sendCommand(o.GetContrAddr(), o.GetPortNumber(), false, "", params, time.Duration(5)*time.Second)
	return code, err
}
