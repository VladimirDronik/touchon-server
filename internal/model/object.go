package model

import (
	"encoding/json"
	"sort"
)

// ChildType используется для работы с дочерними объектами
type ChildType int

const (
	ChildTypeInternal ChildType = iota // Внутренние объекты (port, sensor_value etc)
	ChildTypeExternal                  // Все, кроме внутренних
	ChildTypeAll                       // Все дочерние объекты
	ChildTypeNobody                    // Ни какие дочерние объекты
)

type ObjectStatus string

const (
	StatusAvailable   ObjectStatus = "available"
	StatusUnavailable ObjectStatus = "unavailable"
	StatusDisabled    ObjectStatus = "disabled"
	StatusNA          ObjectStatus = "N/A"
	StatusOn          ObjectStatus = "ON"
	StatusOff         ObjectStatus = "OFF"
)

// StoreObject Объект
type StoreObject struct {
	ID       int             `json:"id"`                  // ID объекта
	ParentID *int            `json:"parent_id,omitempty"` // ID родительского объекта
	ZoneID   *int            `json:"zone_id,omitempty"`   // ID зоны, в которой размещен объект
	Category Category        `json:"category"`            // Категория объекта
	Type     string          `json:"type"`                // Тип объекта
	Internal bool            `json:"internal"`            // Признак внутреннего объекта (port, sensor_value)
	Name     string          `json:"name"`                // Название объекта
	Status   ObjectStatus    `json:"status,omitempty"`    // Состояние объекта
	Tags     map[string]bool `gorm:"serializer:json"`     //
	Enabled  bool            `json:"enabled"`             // Включает методы Start/Shutdown
	Methods  []Method        `gorm:"-"`

	Children []*StoreObject `json:"children,omitempty" gorm:"-"` // Дочерние объекты
}

type Method struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

type JsonObject struct {
	ID       int          `json:"id"`                  // ID объекта
	ParentID *int         `json:"parent_id,omitempty"` // ID родительского объекта
	ZoneID   *int         `json:"zone_id,omitempty"`   // ID зоны, в которой размещен объект
	Category Category     `json:"category"`            // Категория объекта
	Type     string       `json:"type"`                // Тип объекта
	Internal bool         `json:"internal"`            // Признак внутреннего объекта (port, sensor_value)
	Name     string       `json:"name"`                // Название объекта
	Status   ObjectStatus `json:"status,omitempty"`    // Состояние объекта
	Tags     []string     `json:"tags"`                //
	Enabled  bool         `json:"enabled"`             // Включает методы Start/Shutdown
	Methods  []Method     `json:"methods,omitempty"`

	Children []*JsonObject `json:"children,omitempty" gorm:"-"` // Дочерние объекты
}

func (o *StoreObject) MarshalJSON() ([]byte, error) {
	return json.Marshal(StoreObjectToJsonObject(o))
}

func StoreObjectToJsonObject(o *StoreObject) *JsonObject {
	r := &JsonObject{
		ID:       o.ID,
		ParentID: o.ParentID,
		ZoneID:   o.ZoneID,
		Category: o.Category,
		Type:     o.Type,
		Internal: o.Internal,
		Name:     o.Name,
		Status:   o.Status,
		Enabled:  o.Enabled,
		Methods:  o.Methods,
	}

	for tag := range o.Tags {
		r.Tags = append(r.Tags, tag)
	}

	sort.Strings(r.Tags)

	for _, child := range o.Children {
		r.Children = append(r.Children, StoreObjectToJsonObject(child))
	}

	return r
}

func (o *StoreObject) TableName() string {
	return "om_objects"
}
