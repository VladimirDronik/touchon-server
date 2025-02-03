package model

import "touchon-server/lib/mqtt/messages"

type CronTask struct {
	ID          int           `json:"id"`               // 1
	Name        string        `json:"name"`             // check temp
	Description string        `json:"description"`      // Проверяет температуру
	Period      string        `json:"period"`           // 5s
	Enabled     bool          `json:"enabled"`          //
	Actions     []*CronAction `json:"actions" gorm:"-"` //
}

type ActionType = string // Enums(method,delay,notification)

type CronAction struct {
	ID         int                    `json:"id"`                      // 2
	TaskID     int                    `json:"task_id"`                 // 4
	TargetType messages.TargetType    `json:"target_type,omitempty"`   // object, item
	TargetID   int                    `json:"target_id,omitempty"`     // 8
	Type       ActionType             `json:"type"`                    // method, delay
	Name       string                 `json:"name"`                    // script_1, check
	Args       map[string]interface{} `json:"args,omitempty" gorm:"-"` // method or script args
	QoS        messages.QoS           `json:"qos" gorm:"column:qos"`   // 0,1,2
	Enabled    bool                   `json:"enabled"`                 //
	Sort       int                    `json:"sort"`                    //
	Comment    string                 `json:"comment,omitempty"`       // отключено потому что...
}
