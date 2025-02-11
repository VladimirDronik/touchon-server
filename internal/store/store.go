package store

import (
	"time"

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

	// TR

	Users() Users
	Items() Items
	Devices() Devices
	History() History
	Notifications() Notifications
	Boilers() Boilers
	Conditioners() Conditioners
	Curtains() Curtains
	Lights() Lights
	Zones() Zones
	Events() Events
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
	GetObjectByParent(parentID int, typeObject string) (*model.StoreObject, error)
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

// TR

type Users interface {
	AddRefreshToken(deviceID int, refreshToken string, refreshTokenTTL time.Duration) error
	GetByToken(token string) (*model.User, error)
	GetByLoginAndPassword(login string, password string) (*model.User, error)
	GetByDeviceID(deviceID int) (*model.User, error)
	RemoveToken(refreshToken string) error // Удаление данных о сессии в таблице токенов
	LinkDeviceToken(userID int, token string, deviceType string) error
	GetDeviceToken(deviceID int) (*model.DeviceTokens, error)
	GetAllUsers() ([]*model.User, error)
	Create(user *model.User) (int, error)
	Delete(userID int) error
}

type Items interface {
	SaveItem(item *model.ViewItem) error
	GetItem(itemID int) (*model.ViewItem, error)
	UpdateItem(viewItem *model.ViewItem) error
	DeleteItem(itemID int) error
	GetScenarios() ([]*model.Scenario, error) // Отдает элементы для кнопок без комнат
	GetZoneItems() ([]*model.GroupRoom, error)
	GetZoneSensors(zoneID int) ([]*model.Sensor, error)
	GetGroupElements(groupID int) ([]*model.ViewItem, error)
	ChangeItem(itemID int, status string) error                                                                   // Изменение элементов
	GetItemsForChange(targetType messages.TargetType, targetID int, eventName string) ([]*model.ItemForWS, error) // Найти действия для элементов, которые соответствуют событию
	GetZones() ([]*model.Zone, error)                                                                             // Отдает список помещений, где есть элементы
	GetCountersList() ([]*model.Counter, error)
	GetCounter(counterID int) (*model.Counter, error)
	GetZone(zoneID int) (*model.GroupRoom, error) // отдает помещение
	GetMenus(parentIDs ...int) ([]*model.Menu, error)
	SetOrder(itemIDs []int, zoneID int) error // Задает порядок сортировки
}

type Devices interface {
	GetSensor(itemID int) (*model.Sensor, error)
	SetSensorValue(itemID int, value float32) error
	GetDimmer(itemID int) (*model.Dimmer, error)
	SetFieldValue(table string, itemID int, field string, value interface{}) error
	SaveSensor(sensor *model.Sensor) error
	DeleteSensor(itemID int) error
}

type History interface {
	GetHistory(itemID int, itemType model.HistoryItemType, filter model.HistoryFilter) (*model.HistoryPoints, error)
	SetHourlyValue(itemID int, dateTime string, value float32) error
	GenerateHistory(itemID int, itemType model.HistoryItemType, filter model.HistoryFilter, startDate, endDate string, min, max float32) error
}

type Notifications interface {
	GetNotifications(offset int, limit int) ([]*model.Notification, error)
	SetIsRead(id int) error
	AddNotification(n *model.Notification) error
	GetUnreadNotificationsCount() (int, error)
	GetPushTokens() (map[string]string, error)
	SetFieldValue(id int, field string, value interface{}) error
}

type Boilers interface {
	GetBoiler(boilerID int) (*model.Boiler, error)
	SetOutlineStatus(boilerID int, outline string, status string) error
	SetHeatingMode(boilerID int, mode string) error
	SetHeatingTemperature(boilerID int, value float32) error
	UpdateBoilerPresets(boilerID int, presets []*model.BoilerPreset) error
	SetFieldValue(itemID int, field string, value interface{}) error
}

type Conditioners interface {
	GetConditioner(objectID int) (*model.ViewItem, error)
	SetConditionerTemperature(itemID int, value float32) error
	SetConditionerMode(itemID int, mode string, value bool) error
	SetConditionerOperatingMode(itemID int, mode string) error
	SetConditionerDirection(itemID int, plane string, direction string) error
	SetConditionerFanSpeed(itemID int, speed string) error
	SetConditionerExtraMode(itemID int, mode string, value bool) error
	SetFieldValue(itemID int, field string, value interface{}) error
	GetParams(itemID int) (*model.ConditionerParams, error)
}

type Curtains interface {
	GetCurtain(itemID int) (*model.ViewItem, error)
	SetCurtainOpenPercent(itemID int, value float32) error
	SetFieldValue(itemID int, field string, value interface{}) error
	GetParams(itemID int) (*model.CurtainParams, error)
}

type Lights interface {
	GetLight(objectID int) (*model.ViewItem, error)
	SetLightHSVColor(itemID int, hue int, saturation float32, brightness float32) error
	SetLightCCTColor(itemID int, cct int) error
	SetBrightness(itemID int, brightness float32) error
	SetFieldValue(itemID int, field string, value interface{}) error
	GetParams(itemID int) (*model.LightParams, error)
}

type Zones interface {
	CreateZone(zone *model.Zone) (int, error)                        // Создает новое помещение
	GetZoneTrees(parentIDs ...int) ([]*model.Zone, error)            // Отдает список всех помещения
	UpdateZones(zones []*model.Zone) error                           // Изменяет помещения, которые переданы как структура
	SetOrder(zoneIDs []int) error                                    // Задает порядок сортировки
	SetFieldValue(zoneID int, field string, value interface{}) error // Установить значение для поля помещения
	DeleteZone(zoneID int) error                                     //Удаление помещения
}

type Events interface {
	AddEvent(event *model.TrEvent) error
	DeleteEvent(itemID int) error
}
