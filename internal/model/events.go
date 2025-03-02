package model

type TrEvent struct {
	ID         int    `json:"id"`
	TargetType string `json:"target_type"`
	EventName  string `json:"event" gorm:"column:event"`
	TargetID   int    `json:"target_id"`
	Enabled    int    `json:"enabled"`
	//Value      string `json:"value"`
	//ItemID     int    `json:"item_id"`
}

func (o *TrEvent) TableName() string {
	return "ar_events"
}

type EventActions struct {
	ID         int    `json:"id"`
	EventID    int    `json:"event_id"`
	TargetType string `json:"target_type"`
	TargetID   int    `json:"target_id"`
	Type       string `json:"type"`
	Name       string `json:"name"`
	Args       string `json:"args"`
	Qos        string `json:"qos"`
	Enabled    bool   `json:"enabled"`
	Sort       int    `json:"sort"`
	Comment    string `json:"comment"`
}

func (o *EventActions) TableName() string {
	return "ar_event_actions"
}
