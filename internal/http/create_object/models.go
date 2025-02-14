package create_object

import (
	"touchon-server/internal/model"
	"touchon-server/lib/interfaces"
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
		Name    string                    `json:"name"`
		Actions []*interfaces.EventAction `json:"actions"`
	} `json:"events,omitempty"`
}

type Child struct {
	Props    map[string]interface{} `json:"props"`
	Children []*Child               `json:"children,omitempty"`
}

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
