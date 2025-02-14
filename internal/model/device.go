package model

// Sensor Структура итемов сенсора
type Sensor struct {
	ID           int     `json:"-"`
	ViewItemID   int     `json:"item_id"`
	ObjectID     int     `json:"object_id,omitempty"`
	Type         string  `json:"type"`
	MinThreshold float32 `json:"min_threshold,omitempty"`
	MaxThreshold float32 `json:"max_threshold,omitempty"`
	Adjustment   bool    `json:"adjustment,omitempty"`

	Current float32 `json:"current" gorm:"-"`
	Target  float32 `json:"target,omitempty" gorm:"-"`

	ObjectParam string `json:"object_param,omitempty" gorm:"-"`
	ObjectEvent string `json:"object_event,omitempty" gorm:"-"`

	History *HistoryPoints `json:"history,omitempty" gorm:"-"`

	//Параметры, которые нужны для создания view
	ZoneID  int    `json:"zone_id,omitempty"`
	Icon    string `json:"icon,omitempty"`
	Enabled bool   `json:"enabled,omitempty"`
	Title   string `json:"title,omitempty"`
	Auth    string `json:"auth,omitempty"`
	Sort    int    `json:"sort,omitempty"`

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
