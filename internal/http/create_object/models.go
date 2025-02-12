package create_object

import (
	"touchon-server/internal/model"
	"touchon-server/lib/mqtt/messages"
)

// Request

type Request struct {
	Object struct {
		ParentID *int                   `json:"parent_id,omitempty"`
		ZoneID   *int                   `json:"zone_id,omitempty"`
		Category model.Category         `json:"category"`
		Type     string                 `json:"type"`
		Name     string                 `json:"name"`
		Props    map[string]interface{} `json:"props,omitempty"`
		Enabled  bool                   `json:"enabled"`
		Children []*Child               `json:"children,omitempty"`
	} `json:"object"`

	Events []struct {
		Name    string    `json:"name"`
		Actions []*Action `json:"actions"`
	} `json:"events,omitempty"`
}

type Child struct {
	Props    map[string]interface{} `json:"props"`
	Children []*Child               `json:"children,omitempty"`
}

type Action struct {
	TargetType messages.TargetType    `json:"target_type"` // Enums(not_matters,object,item,script,service)
	TargetID   int                    `json:"target_id"`   // 8
	Type       ActionType             `json:"type"`        // Enums(method,delay,notification)
	Name       string                 `json:"name"`        // script_1, check
	Args       map[string]interface{} `json:"args"`        // method or script args
	QoS        messages.QoS           `json:"qos"`         // Enums(0,1,2)
	Enabled    bool                   `json:"enabled"`     // Enums(true,false)
	Sort       int                    `json:"sort"`        //
	Comment    string                 `json:"comment"`     // отключено потому что...
}

type ActionType = string // Enums(method,delay,notification)

// Response

type Response struct {
	ID int `json:"id"`
	//ParentID int `json:"parent_id"`
	//ZoneID   int `json:"zone_id"`
	//
	//Category model.Category `json:"category"`
	//Type     string         `json:"type"`
	//Internal bool           `json:"internal"` // Признак внутреннего объекта (port, sensor_value)
	//Name     string         `json:"name"`
	//Status   string         `json:"status"`
	//Tags     []string       `json:"tags,omitempty"`
	//
	//Props    map[string]interface{} `json:"props,omitempty"`
	//Children []*Response            `json:"children,omitempty"`

	//Events   *Events                `json:"events,omitempty"`
	//Methods  *Methods               `json:"methods,omitempty"`
}
