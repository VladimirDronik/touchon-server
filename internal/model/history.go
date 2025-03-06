package model

// ServerDate структура текущей даты сервера
type ServerDate struct {
	Hour  string `json:"hour"`
	Day   string `json:"day"`
	Month string `json:"month"`
	Year  string `json:"year"`
}

// HistoryPoints структура данных для графиков счетчика
type HistoryPoints struct {
	ServerDate  *ServerDate     `json:"server_date"`
	DayPoints   []*HistoryPoint `json:"day_points,omitempty"`
	WeekPoints  []*HistoryPoint `json:"week_points,omitempty"`
	MonthPoints []*HistoryPoint `json:"month_points,omitempty"`
	YearPoints  []*HistoryPoint `json:"year_points,omitempty"`
}

// HistoryPoint модель элемента в истории значений
type HistoryPoint struct {
	ID       int      `json:"id"`
	ItemID   int      `json:"item_id"`
	Datetime string   `json:"date"`
	Value    *float32 `json:"value,omitempty"`

	FormattedDate string `json:"formatted_date" gorm:"-"`
}

// HistoryFilter ENUM возможных фильтров для истории значений
type HistoryFilter string

const (
	HistoryFilterDay   HistoryFilter = "day"
	HistoryFilterWeek  HistoryFilter = "week"
	HistoryFilterMonth HistoryFilter = "month"
	HistoryFilterYear  HistoryFilter = "year"
)

type DestTable string

const (
	TableDailyHistory   DestTable = "daily_history"
	TableMonthlyHistory DestTable = "monthly_history"
)

// HistoryItemType ENUM возможных типов сущностей имеющих историю значений
type HistoryItemType string

const (
	HistoryItemTypeDeviceObject  HistoryItemType = "device"
	HistoryItemTypeCounterObject HistoryItemType = "counter"
)
