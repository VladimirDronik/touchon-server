package model

// Boiler модель котла
type Boiler struct {
	ID                 int     `json:"id"`
	HeatingStatus      string  `json:"heating_status,omitempty"`
	WaterStatus        string  `json:"water_status,omitempty"`
	HeatingCurrentTemp float32 `json:"heating_current_temp,omitempty"`
	HeatingOptimalTemp float32 `json:"heating_optimal_temp,omitempty"`
	WaterCurrentTemp   float32 `json:"water_current_temp,omitempty"`
	HeatingMode        string  `json:"heating_mode,omitempty"`
	IndoorTemp         float32 `json:"indoor_temp"`
	OutdoorTemp        float32 `json:"outdoor_temp"`
	MinThreshold       float32 `json:"min_threshold,omitempty"`
	MaxThreshold       float32 `json:"max_threshold,omitempty"`
	Icon               string  `json:"icon,omitempty"`
	Title              string  `json:"title,omitempty"`
	Color              string  `json:"color,omitempty"`
	Auth               string  `json:"auth,omitempty"`

	Presets    []*BoilerPreset   `json:"presets,omitempty" gorm:"-"`
	Properties []*BoilerProperty `json:"properties,omitempty" gorm:"-"`
	History    *HistoryPoints    `json:"history,omitempty" gorm:"-"`
}

// BoilerPreset модель предустановок для котла
type BoilerPreset struct {
	ID          int     `json:"id,omitempty" swaggerignore:"true"`
	BoilerID    int     `json:"id_boiler,omitempty" swaggerignore:"true"`
	TempOut     float32 `json:"temp_out" default:"-20"`
	TempCoolant float32 `json:"temp_coolant,omitempty" default:"60"`
}

// BoilerProperty модель свойства котла
type BoilerProperty struct {
	ID        int    `json:"id"`
	BoilerID  int    `json:"boiler_id,omitempty"`
	Title     string `json:"title,omitempty"`
	ImageName string `json:"image_name,omitempty"`
	Value     string `json:"value,omitempty"`
	Status    string `json:"status,omitempty"`
}
