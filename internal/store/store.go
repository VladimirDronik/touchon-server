package store

import (
	"github.com/pkg/errors"
	"touchon-server/internal/model"
	"touchon-server/lib/mqtt/messages"
)

var ErrNotFound = errors.New("not found")

// Global instance
var I Store

type Store interface {
	ObjectRepository() ObjectRepository
	PortRepository() PortRepository
	DeviceRepository() DeviceRepository
	ScriptRepository() ScriptRepository

	// AR

	EventsRepo() EventsRepo
	EventActionsRepo() EventActionsRepo
	CronRepo() CronRepo
}

type ObjectRepository interface {
	GetProp(objectID int, code string) (string, error)
	SetProp(objectID int, code, value string) error // create, update
	DelProp(objectID int, code string) error

	GetProps(objectID int) (map[string]*model.StoreProp, error)
	GetPropsByObjectIDs(objectIDs []int) (map[int]map[string]*model.StoreProp, error)
	GetPropByObjectIDs(objectIDs []int, code string) (map[int]*model.StoreProp, error)
	SetProps(objectID int, props map[string]string) error // create, update
	DelProps(objectID int) error

	SetObjectStatus(objectID int, status string) error
	GetObject(objectID int) (*model.StoreObject, error)
	SetParent(objectID int, parentID *int) error //Установка родителя для объекта
	SaveObject(object *model.StoreObject) error  // create, update
	GetObjects(filters map[string]interface{}, tags []string, offset, limit int, objectType model.ChildType) ([]*model.StoreObject, error)
	GetObjectsByAddress(address []string) ([]*model.StoreObject, error) // Выводит объекты, у которых адрес совпадает с искомым
	GetTotal(filters map[string]interface{}, tags []string, objectType model.ChildType) (int, error)
	GetObjectsByTags(tags []string, offset, limit int, objectType model.ChildType) ([]*model.StoreObject, error)
	GetTotalByTags(tags []string, objectType model.ChildType) (int, error)
	GetObjectsByIDs(ids []int) ([]*model.StoreObject, error)
	GetObjectChildren(childType model.ChildType, objectID ...int) ([]*model.StoreObject, error)
	DelObject(objectID int) error
	GetAllTags() (map[string]int, error)
}

type PortRepository interface {
	GetPortObjectID(controllerID, portNumber string) (int, error)
}

type DeviceRepository interface {
	GetControllerByName(name string) (*model.StoreObject, error)
}

type ScriptRepository interface {
	GetScript(id int) (*model.StoreScript, error)
	GetScriptByCode(code string) (*model.StoreScript, error)
	SetScript(script *model.StoreScript) error // create, update
	GetScripts(code, name string, offset, limit int) ([]*model.StoreScript, error)
	// GetTotal возвращает общее кол-во найденных скриптов
	GetTotal(code, name string, offset, limit int) (int, error)
	DelScript(id int) error
}

// AR

type EventsRepo interface {
	// GetEvents возвращает события сущности.
	GetEvents(targetType messages.TargetType, targetID int) ([]*model.Event, error)

	// GetEvent возвращает событие.
	GetEvent(targetType messages.TargetType, targetID int, eventName string) (*model.Event, error)

	// SaveEvent добавляет или обновляет событие.
	SaveEvent(ev *model.Event) error

	// DeleteEvent удаляет событие.
	DeleteEvent(targetType messages.TargetType, targetID int, eventName string) error

	// GetAllEventsName возвращает названия всех событий, используемых в таблице.
	// Используется для проверки правильности указанных имен.
	GetAllEventsName() ([]string, error)
}

type EventActionsRepo interface {
	// GetActions возвращает список действий для события.
	GetActions(eventIDs ...int) (map[int][]*model.EventAction, error)

	// GetActionsCount возвращает количество действий для событий.
	GetActionsCount(eventIDs ...int) (map[int]int, error)

	// SaveAction добавляет или обновляет действие.
	SaveAction(act *model.EventAction) error

	// DeleteAction удаляет действие.
	DeleteAction(actID int) error
	DeleteActionByObject(targetType string, objectID int) error

	// OrderActions меняет порядок действий.
	OrderActions(actIDs []int) error
}

type CronRepo interface {
	// GetEnabledTasks возвращает активные задачи.
	GetEnabledTasks() ([]*model.CronTask, error)
	CreateTask(task *model.CronTask) (int, error)
	CreateTaskAction(action *model.CronAction) error
	DeleteTask(objectID int, target string) error
	UpdateTask(task *model.CronTask) error
	GetCronAction(objectID int, targetType string) (*model.CronAction, error)
}
