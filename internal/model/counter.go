package model

// Counter модель счетчика
type Counter struct {
	ID           int      `json:"id"`
	Name         string   `json:"name"`
	Type         string   `json:"type"`
	Unit         string   `json:"unit,omitempty"`
	TodayValue   float32  `json:"today_value"`
	WeekValue    float32  `json:"week_value"`
	MonthValue   float32  `json:"month_value"`
	YearValue    float32  `json:"year_value"`
	Value        float32  `json:"value"`
	PriceForUnit *float32 `json:"price_for_unit,omitempty"`
	Sort         int      `json:"sort,omitempty"`
	Enabled      bool     `json:"enabled"`

	History *HistoryPoints `json:"history,omitempty" gorm:"-"`
}
