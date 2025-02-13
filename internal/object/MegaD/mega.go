package MegaD

import (
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

type method struct {
	controller *Controller
}

func (o *Controller) Initialization(objectID int) error {
	o.ObjectID = objectID
	return nil
}

// Adding добавление нового объекта в БД
func (o *Controller) Adding() error {
	return nil
}
