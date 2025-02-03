package model

// Sensor Структура итемов сенсора
type Sensor struct {
	ID           int     `json:"-"`
	ViewItemID   int     `json:"item_id"`
	ZoneID       int     `json:"zone_id"`
	Type         string  `json:"type"`
	Name         string  `json:"name"`
	Current      float32 `json:"current"`
	Optimal      float32 `json:"optimal,omitempty"`
	MinThreshold float32 `json:"min_threshold,omitempty"`
	MaxThreshold float32 `json:"max_threshold,omitempty"`
	Icon         string  `json:"icon,omitempty"`
	PositionLeft int     `json:"position_left"`
	PositionTop  int     `json:"position_top"`
	Sort         int     `json:"sort"`
	Auth         string  `json:"auth,omitempty"`
	Enabled      bool    `json:"enabled"`

	ObjectID    int    `json:"object_id,omitempty" gorm:"-"`
	ObjectParam string `json:"object_param,omitempty" gorm:"-"`
	ObjectEvent string `json:"object_event,omitempty" gorm:"-"`

	History *HistoryPoints `json:"history,omitempty" gorm:"-"`

	// Backward compatibility
	Status string `json:"status,omitempty" gorm:"-"` // on|off
}

// LightParams Структура параметров источника света
type LightParams struct {
	ID         int     `json:"id"`
	ViewItemID int     `json:"view_item_id"`
	Hue        int     `json:"hue,omitempty"`
	Saturation float32 `json:"saturation,omitempty"`
	Brightness float32 `json:"brightness,omitempty"`
	Cct        int     `json:"cct,omitempty"`
}

// Dimmer Структура димера
type Dimmer struct {
	ID         int    `json:"Id"`
	ViewItemID int    `json:"view_item_id"`
	Name       string `json:"Name"`
	Value      int    `json:"Value"`
	Enabled    bool   `json:"Status"`
}
