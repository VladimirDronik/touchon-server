package objects

import (
	"encoding/json"
	"errors"

	"touchon-server/internal/model"
)

var ErrObjectDisabled = errors.New("object disabled")

// Object Абстрактный тип объекта
type Object interface {
	// Init заполняет модель данными из БД
	Init(storeObj *model.StoreObject, childType model.ChildType) error
	Save() error

	GetID() int
	SetID(int)

	GetParentID() *int
	SetParentID(*int)

	GetZoneID() *int
	SetZoneID(*int)

	// GetCategory Возвращает категорию объекта (controller, sensor, module, ext,  etc)
	GetCategory() model.Category
	SetCategory(model.Category)

	// GetType Возвращает конкретный тип объекта (MegaD)
	GetType() string
	SetType(string)

	GetInternal() bool
	SetInternal(bool)

	// GetName Возвращает название объекта (MegaD2561)
	GetName() string
	SetName(string)

	GetStatus() model.ObjectStatus
	SetStatus(model.ObjectStatus)

	// GetProps Возвращает объект для работы со свойствами
	GetProps() *Props

	GetChildren() *Children
	GetEvents() *Events
	GetMethods() *Methods

	GetTags() []string
	ReplaceTags(tags ...string)
	SetTags(tags ...string)
	DeleteTags(tags ...string)
	GetTagsMap() map[string]bool
	SetTagsMap(map[string]bool)

	GetEnabled() bool
	SetEnabled(bool)

	// CheckEnabled метод проверяет, включен ли объект.
	// Вызывается в методах Start и Shutdown.
	// Присутствует в интерфейсе для того, чтобы можно было
	// его вызвать "вручную" в производных типах.
	CheckEnabled() error

	// Start запускает логику объекта
	Start() error

	// Shutdown завершает логику объекта
	Shutdown() error

	GetStoreObject() *model.StoreObject

	DeleteChildren() error

	// Marshaler Объект должен уметь генерировать json-модель.
	json.Marshaler

	// Unmarshaler Объект должен создаваться из json
	json.Unmarshaler
}

type ObjectModel struct {
	ID       int  `json:"id"`
	ParentID *int `json:"parent_id"`
	ZoneID   *int `json:"zone_id"`

	Category model.Category     `json:"category"`
	Type     string             `json:"type"`
	Internal bool               `json:"internal"` // Признак внутреннего объекта (port, sensor_value)
	Name     string             `json:"name"`
	Status   model.ObjectStatus `json:"status"`
	Tags     []string           `json:"tags,omitempty"`
	Enabled  bool               `json:"enabled"`

	Props    *Props    `json:"props,omitempty"`
	Children *Children `json:"children,omitempty"`
	Events   *Events   `json:"events,omitempty"`
	Methods  *Methods  `json:"methods,omitempty"`
}
