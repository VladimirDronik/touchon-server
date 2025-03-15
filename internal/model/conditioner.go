package model

// Conditioner модель параметров кондиционера
type Conditioner struct {
	CondParams
	OperatingModes       map[string]string `json:"operating_modes"`
	FanSpeeds            map[string]string `json:"fan_speeds"`
	VerticalDirections   map[string]string `json:"vertical_directions"`
	HorizontalDirections map[string]string `json:"horizontal_directions"`
}

type StoreConditionerParams struct {
	CondParams
	OperatingModes       string `json:"-"`
	FanSpeeds            string `json:"-"`
	VerticalDirections   string `json:"-"`
	HorizontalDirections string `json:"-"`
}

type CondParams struct {
	ID          int     `json:"id"`
	ViewItemID  int     `json:"view_item_id"`
	InsideTemp  float32 `json:"inside_temp"`
	OutsideTemp float32 `json:"outside_temp"`
	CurrentTemp float32 `json:"current_temp,omitempty"`
	TargetTemp  float32 `json:"target_temp"`

	MinThreshold float32 `json:"min_threshold"`
	MaxThreshold float32 `json:"max_threshold"`

	SilentMode bool `json:"silent_mode,omitempty"`
	EcoMode    bool `json:"eco_mode,omitempty"`
	TurboMode  bool `json:"turbo_mode,omitempty"`
	SleepMode  bool `json:"sleep_mode,omitempty"`

	OperatingMode       int `json:"operating_mode,omitempty"`
	FanSpeed            int `json:"fan_speed,omitempty"`
	VerticalDirection   int `json:"vertical_direction,omitempty"`
	HorizontalDirection int `json:"horizontal_direction,omitempty"`

	Ionisation       bool `json:"ionisation,omitempty"`
	SelfCleaning     bool `json:"self_cleaning,omitempty"`
	AntiMold         bool `json:"anti_mold,omitempty"`
	Sound            bool `json:"sound,omitempty"`
	OnDutyHeating    bool `json:"on_duty_heating,omitempty"`
	SoftTop          bool `json:"soft_top,omitempty"`
	DisplayBacklight bool `json:"display_backlight,omitempty"`
	PowerStatus      bool `json:"power_status,omitempty"`
}

type ConditionerItem struct {
	ViewItemID   int     `json:"view_item_id"`
	ObjectID     int     `json:"object_id" gorm:"object_id"`
	MinThreshold float32 `json:"min_threshold" gorm:"min_threshold"`
	MaxThreshold float32 `json:"max_threshold" gorm:"max_threshold"`
	Title        string  `json:"title" gorm:"-"`
	Icon         string  `json:"icon" gorm:"-"`
	Auth         string  `json:"auth" gorm:"-"`
}
