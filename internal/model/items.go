package model

type ItemForWS struct {
	ID        int     `json:"id"`
	Type      string  `json:"type,omitempty"`
	Status    string  `json:"status,omitempty"`
	Params    string  `json:"-"` // Используется для выборки из БД, но не для выдачи в ответе сервера
	Value     float32 `json:"value,omitempty" gorm:"-"`
	EventArgs []byte  `json:"-"` // Параметры, в которых содержатся дополнительные аргументы, которые можно использовать для обработки логики при возникновении события
}

type SensorItem struct {
	ItemID     int     `json:"item_id"`
	ZoneID     int     `json:"zone_id,omitempty"`
	Type       string  `json:"type,omitempty"`
	Icon       string  `json:"icon,omitempty"`
	Title      string  `json:"title,omitempty"`
	Auth       string  `json:"auth,omitempty"`
	Adjustment bool    `json:"adjustment,omitempty"`
	Current    float32 `json:"current,omitempty"`
	ObjectID   int     `json:"-"`
}
