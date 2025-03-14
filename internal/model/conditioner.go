package model

// ConditionerParams модель параметров кондиционера
type ConditionerParams struct {
	CondParams
	OperatingModes       map[string]string `json:"operating_modes"`
	FanSpeeds            []string          `json:"fan_speeds"`
	VerticalDirections   []string          `json:"vertical_directions"`
	HorizontalDirections []string          `json:"horizontal_directions"`
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
	OptimalTemp float32 `json:"optimal_temp"`

	MinThreshold float32 `json:"min_threshold"`
	MaxThreshold float32 `json:"max_threshold"`

	SilentMode bool `json:"silent_mode,omitempty"`
	EcoMode    bool `json:"eco_mode,omitempty"`
	TurboMode  bool `json:"turbo_mode,omitempty"`
	SleepMode  bool `json:"sleep_mode,omitempty"`

	OperatingMode       string `json:"operating_mode,omitempty"`
	FanSpeed            string `json:"fan_speed,omitempty"`
	VerticalDirection   string `json:"vertical_direction,omitempty"`
	HorizontalDirection string `json:"horizontal_direction,omitempty"`

	Ionisation    bool `json:"ionisation,omitempty"`
	SelfCleaning  bool `json:"self_cleaning,omitempty"`
	AntiMold      bool `json:"anti_mold,omitempty"`
	Sound         bool `json:"sound,omitempty"`
	OnDutyHeating bool `json:"on_duty_heating,omitempty"`
	SoftTop       bool `json:"soft_top,omitempty"`
}
