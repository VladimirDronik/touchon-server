package MegaD

import (
	//"C"
	"github.com/sirupsen/logrus"
	"touchon-server/internal/store"
)

type Controller struct {
	ObjectID     int // ObjectID Объекта
	Name         string
	ControllerID int //ид контроллера, на котором находится порт
	method       *method
	event        *event
	IP           string `json:"ip"`          // IP контроллера
	Number       string `json:"number_port"` // Номер порта
	EventName    string `json:"event_name"`  // Название события, которое сработало
	Params       string `json:"params"`      // Параметры, которые пришли событиями
	Store        store.Store
	Logger       *logrus.Logger
}

type event struct {
	contr *Controller
}

//func (o *Controller) GenerateEvent(eventName string, params string) model.Message {
//
//	return model.Message{}
//}

// OnChangeStatus Событие, которое возникает при смене статуса объекта
//func (o *Controller) OnChangeStatus(values ...interface{}) model.Message {
//	var message model.Message
//
//	//TODO:: сделать обработку смены статуса
//
//	println(o.ObjectID)
//	//act := action.Action{Store: e.controller.Store}
//	//act.RunActions(e.controller.ObjectID, "onChangeStatus")
//	//for i, v := range values {
//	//
//	//	switch i {
//	//	case 1:
//	//		println("cmd=", v)
//	//	case 2:
//	//		println("click=", v)
//	//	case 3:
//	//		println("m=", v)
//	//
//	//	}
//	//}
//
//	return message
//}

type method struct {
	controller *Controller
}

// Структура, которая будет отправлена в mqtt топик
type content struct {
	Ip string
}

func (o *Controller) Initialization(objectID int) error {
	o.ObjectID = objectID
	return nil
}

// RunMethod выполнение метода по его названию
//func (o *Controller) RunMethod(methodName string, params map[string]interface{}) ([]model.Message, model.PayloadStruct, error) {
//	var message model.Message
//
//	switch methodName {
//
//	}
//
//	return append([]model.Message{}, message), model.PayloadStruct{}, nil
//}

//// On Включаем порт
//func (c *Controller) On() (model.Message, error) {
//
//	println("OK")
//	//Меняем статус объекта
//
//	//Отправляем даные контроллеру (формируем сообщение)
//	mqttMess.Topic = "device/MegaD2561/self"
//	mqttMess.Content = "{}"
//
//	return mqttMess, nil
//
//	//Вызываем событие onChangeStatus
//	//method.controller.Events().OnChangeStatus()
//}
//
//// Off Выключаем порт
//func (c *Controller) Off() (model.Message, error) {
//	return model.Message{}, nil
//}

// Adding добавление нового объекта в БД
func (o *Controller) Adding() error {
	return nil
}

// Reboot перезагрузка контроллера
//func (o *Controller) Reboot() (model.Message, error) {
//	//cont := content{Address: method.controller.IP}
//	//
//	//mqttMess.Topic = "device/MegaD2561/self"
//	//mqttMess.Content = "{}"
//
//	return model.Message{}, nil
//}
