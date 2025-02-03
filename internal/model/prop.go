package model

import "encoding/json"

// StoreProp Свойство объекта
type StoreProp struct {
	ID       int    `json:"id"`        // ID свойства
	ObjectID int    `json:"object_id"` // ID объекта, содержащего данное свойство
	Code     string `json:"code"`      // Уникальный в рамках объекта код свойства
	Value    string `json:"value"`     // Значение свойства
}

func (o *StoreProp) TableName() string {
	return "om_props"
}

// StoreScript Скрипт
type StoreScript struct {
	ID          int             `json:"id"`
	Code        string          `json:"code"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Params      json.RawMessage `json:"params"`
	Body        string          `json:"body"`
}

func (o *StoreScript) TableName() string {
	return "om_scripts"
}
