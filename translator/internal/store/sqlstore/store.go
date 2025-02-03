package sqlstore

import (
	"gorm.io/gorm"
	"translator/internal/store"
)

type Store struct {
	db            *gorm.DB
	users         store.Users
	items         store.Items
	devices       store.Devices
	history       store.History
	notifications store.Notifications
	boilers       store.Boilers
	conditioners  store.Conditioners
	curtains      store.Curtains
	lights        store.Lights
	zones         store.Zones
	events        store.Events
}

func New(db *gorm.DB) *Store {
	o := &Store{db: db}
	o.users = &Users{store: o}
	o.items = &Items{store: o}
	o.devices = &Devices{store: o}
	o.history = &History{store: o}
	o.notifications = &Notifications{store: o}
	o.boilers = &Boilers{store: o}
	o.conditioners = &Conditioners{store: o}
	o.curtains = &Curtains{store: o}
	o.lights = &Lights{store: o}
	o.zones = &Zones{store: o}
	o.events = &Events{store: o}
	return o
}

func (o *Store) Users() store.Users {
	return o.users
}

func (o *Store) Items() store.Items {
	return o.items
}

func (o *Store) Devices() store.Devices {
	return o.devices
}

func (o *Store) History() store.History {
	return o.history
}

func (o *Store) Notifications() store.Notifications {
	return o.notifications
}

func (o *Store) Boilers() store.Boilers {
	return o.boilers
}

func (o *Store) Conditioners() store.Conditioners {
	return o.conditioners
}

func (o *Store) Curtains() store.Curtains {
	return o.curtains
}

func (o *Store) Lights() store.Lights {
	return o.lights
}

func (o *Store) Zones() store.Zones {
	return o.zones
}

func (o *Store) Events() store.Events {
	return o.events
}
