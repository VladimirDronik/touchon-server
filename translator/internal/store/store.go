package store

import (
	"time"

	"github.com/VladimirDronik/touchon-server/mqtt/messages"
	"translator/internal/model"
)

type Store interface {
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
	SaveItem(item *model.ViewItem) (int, error)
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
	AddEvent(event *model.Event) error
	DeleteEvent(itemID int) error
}
