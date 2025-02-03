package model

type User struct {
	ID           int    `json:"Id"`
	Login        string `json:"login"`
	Password     string `json:"password"`
	Role         string `json:"role"`
	SendPush     bool   `json:"send_push"`
	RefreshToken string `json:"RefreshToken"`
	TokenExpired string `json:"TokenExpired"`
	DeviceID     int    `json:"DevId"`
	DeviceTokens
}

type DeviceTokens struct {
	DeviceType  string `json:"DeviceType"`
	DeviceToken string `json:"DeviceToken"`
}
