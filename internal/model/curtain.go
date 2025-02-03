package model

// CurtainParams модель значений графика для штор
type CurtainParams struct {
	ID          int      `json:"id"`
	ViewItemID  int      `json:"view_item_id"`
	Type        string   `json:"type"`
	ControlType string   `json:"control_type"`
	OpenPercent *float32 `json:"open_percent"`
}
