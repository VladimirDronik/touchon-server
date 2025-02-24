package interfaces

type CronTask struct {
	ID          int           `json:"id"`               // 1
	Name        string        `json:"name"`             // check temp
	Description string        `json:"description"`      // Проверяет температуру
	Period      string        `json:"period"`           // 5s
	Enabled     bool          `json:"enabled"`          //
	Actions     []*CronAction `json:"actions" gorm:"-"` //
}

func (o *CronTask) TableName() string {
	return "ar_cron_tasks"
}

type ActionType = string

const (
	ActionTypeMethod       ActionType = "method"
	ActionTypeDelay        ActionType = "delay"
	ActionTypeNotification ActionType = "notification"
)

type CronAction struct {
	ID         int                    `json:"id"`                      // 2
	TaskID     int                    `json:"task_id"`                 // 4
	TargetType TargetType             `json:"target_type,omitempty"`   // object, item
	TargetID   int                    `json:"target_id,omitempty"`     // 8
	Type       ActionType             `json:"type"`                    // method, delay
	Name       string                 `json:"name"`                    // script_1, check
	Args       map[string]interface{} `json:"args,omitempty" gorm:"-"` // method or script args
	Enabled    bool                   `json:"enabled"`                 //
	Sort       int                    `json:"sort"`                    //
	Comment    string                 `json:"comment,omitempty"`       // отключено потому что...
}

func (o *CronAction) TableName() string {
	return "ar_cron_actions"
}

// AREvent - Action router event
type AREvent struct {
	ID         int            `json:"id"`                    // 1
	TargetType TargetType     `json:"target_type,omitempty"` // object, item
	TargetID   int            `json:"target_id,omitempty"`   // 8
	EventName  string         `json:"event_name"`            //
	Enabled    bool           `json:"enabled"`               //
	Actions    []*EventAction `json:"actions" gorm:"-"`      //
}

func (o *AREvent) TableName() string {
	return "ar_events"
}

type EventAction struct {
	ID         int                    `json:"id,omitempty"`            // 2
	EventID    int                    `json:"event_id,omitempty"`      // 4
	TargetType TargetType             `json:"target_type,omitempty"`   // object, item
	TargetID   int                    `json:"target_id,omitempty"`     // 8
	Type       ActionType             `json:"type"`                    // script, method
	Name       string                 `json:"name"`                    // script_1, check
	Args       map[string]interface{} `json:"args,omitempty" gorm:"-"` // method or script args
	Enabled    bool                   `json:"enabled"`                 //
	Sort       int                    `json:"sort"`                    //
	Comment    string                 `json:"comment,omitempty"`       // отключено потому что...
}

func (o *EventAction) TableName() string {
	return "ar_event_actions"
}
