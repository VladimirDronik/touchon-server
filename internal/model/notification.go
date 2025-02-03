package model

import "touchon-server/lib/events"

type Notification struct {
	ID     int               `json:"id,omitempty"`
	Type   events.NotifyType `json:"type,omitempty"`
	Date   string            `json:"date"`
	Text   string            `json:"text"`
	IsRead bool              `json:"is_read"`
}

type GroupedNotifications struct {
	Day           string          `json:"day"`
	Notifications []*Notification `json:"notifications"`
}

type PushNotification struct {
	Tokens map[string]string `json:"tokens,omitempty"`
	Title  string            `json:"title"`
	Body   string            `json:"body"`
	Sound  string            `json:"sound,omitempty"`
}
