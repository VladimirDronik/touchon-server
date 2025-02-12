package model

import (
	"touchon-server/lib/interfaces"
)

type Notification struct {
	ID     int                         `json:"id,omitempty"`
	Type   interfaces.NotificationType `json:"type,omitempty"`
	Date   string                      `json:"date"`
	Text   string                      `json:"text"`
	IsRead bool                        `json:"is_read"`
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
