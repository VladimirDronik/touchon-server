//go:build ignore

package Module

import (
	"strconv"

	"touchon-server/internal/model"
	"touchon-server/internal/store"
)

type Module struct {
	IdObject     int //IdObject Объекта
	Name         string
	Status       string //состояние модуля вкл/выкл
	ControllerId int    `json:"controller_id"` //ид контроллера, на котором находится модуль
	method       *method
	event        *event
	Port         string `json:"number_port"` // номер порта, к которому подключен
	Store        store.Store
	Config       config.Config
}

func (c *Module) Captions() *Module {
	return &Module{}
}

type event struct {
	contr *Module
}

func (c *Module) GenerateEvent(eventName string, params string) model.Message {

	return model.Message{}
}

// OnChangeStatus Событие, которое возникает при смене статуса объекта
func (c *Module) OnChangeStatus(values ...interface{}) model.Message {
	var message model.Message

	//TODO:: сделать обработку смены статуса

	println(c.IdObject)

	return message

}

type method struct {
	controller *Module
}

// Струкутура, которая будет отправлена в mqtt топик
type content struct {
	Ip string
}

var mqttMess model.Message
var err error

// Init инициализация объекта
func (c *Module) Initialization(objectID int) error {

	var (
		controllerId int
		status       string
	)

	props, err := c.Store.ObjectRepository().GetProps(objectID)
	if err != nil {
		return err
	}

	for _, v := range props {
		switch v.Code {
		case "coontroller_id":
			controllerId, _ = strconv.Atoi(v.Value)
		}
	}

	c.ControllerId = controllerId
	c.Status = status
	c.IdObject = objectID

	return nil
}

// RunMethod выполнение метода по его названию
func (c *Module) RunMethod(methodName string, params map[string]interface{}) ([]model.Message, model.PayloadStruct, error) {

	var message model.Message

	switch methodName {

	}

	return append([]model.Message{}, message), model.PayloadStruct{}, err
}

// Adding добавление нового объекта в БД
func (c *Module) Adding() error {

	return nil
}

// Reboot перезагрузка контроллера
func (c *Module) Reboot() (model.Message, error) {

	//cont := content{Address: method.controller.IP}
	//
	//mqttMess.Topic = "device/MegaD2561/self"
	//mqttMess.Content = "{}"

	return model.Message{}, nil
}
