package update_object

// Request

type Request struct {
	ID       int                    `json:"id"`
	ParentID *int                   `json:"parent_id"`
	ZoneID   *int                   `json:"zone_id"`
	Name     string                 `json:"name"`
	Props    map[string]interface{} `json:"props"`
	Children []Child                `json:"children"`
}

type Child struct {
	Props    map[string]interface{} `json:"props"`
	Children []Child                `json:"children"`
}

// Response

type Response struct {
	//ID       int `json:"id"`
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
