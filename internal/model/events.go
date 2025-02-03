package model

type Event struct {
	ID         int    `json:"id"`
	TargetType string `json:"target_type"`
	EventName  string `json:"event" gorm:"column:event"`
	TargetID   int    `json:"target_id"`
	Value      string `json:"value"`
	ItemID     int    `json:"item_id"`
}
