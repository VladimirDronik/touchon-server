package model

// Counter модель счетчика
type Counter struct {
	ID           int      `json:"id"`
	Name         string   `json:"name"`
	Type         string   `json:"type"`
	Unit         string   `json:"unit,omitempty,omitempty"`
	TodayValue   *float32 `json:"today_value,omitempty"`
	WeekValue    *float32 `json:"week_value,omitempty"`
	MonthValue   *float32 `json:"month_value,omitempty"`
	YearValue    *float32 `json:"year_value,omitempty"`
	Value        float32  `json:"value"`
	PriceForUnit *float32 `json:"price_for_unit,omitempty"`
	Sort         int      `json:"sort,omitempty"`
	Enabled      bool     `json:"enabled"`

	History *HistoryPoints `json:"history,omitempty" gorm:"-"`
}
