package model

type Menu struct {
	ID       int    `json:"id"`
	ParentID int    `json:"parent_id"`
	Page     string `json:"page"`
	Title    string `json:"title"`
	Image    string `json:"image"`
	Sort     int    `json:"sort"`
	Params   string `json:"-"`
	Enabled  bool   `json:"enabled"`

	ParamsOutput interface{} `json:"params,omitempty" gorm:"-"`
	Children     []*Menu     `json:"child_menu_items,omitempty" gorm:"-"`
}
