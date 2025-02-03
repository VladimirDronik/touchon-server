package store

import "touchon-server/internal/model"

// Global instance
var I Store

type Store interface {
	ObjectRepository() ObjectRepository
	PortRepository() PortRepository
	DeviceRepository() DeviceRepository
	ScriptRepository() ScriptRepository
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
