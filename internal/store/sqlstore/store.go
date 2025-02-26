package sqlstore

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"touchon-server/internal/store"
)

var errNotFound = errors.New("not found")

type Store struct {
	db         *gorm.DB
	objectRepo *ObjectRepository
	portRepo   *PortRepository
	deviceRepo *DeviceRepository
	scriptRepo *ScriptRepository

	// AR
	eventsRepo       *EventsRepo
	eventActionsRepo *EventActionsRepo
	cronRepo         *CronRepo

	// TR
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

// New ...
func New(db *gorm.DB) *Store {
	o := &Store{db: db}
	// OM
	o.objectRepo = &ObjectRepository{store: o}
	o.portRepo = &PortRepository{store: o}
	o.deviceRepo = &DeviceRepository{store: o}
	o.scriptRepo = &ScriptRepository{store: o}

	// AR
	o.eventsRepo = &EventsRepo{store: o}
	o.eventActionsRepo = &EventActionsRepo{store: o}
	o.cronRepo = &CronRepo{store: o}

	// TR
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

func (o *Store) GetDB() *gorm.DB {
	return o.db
}

// OM

func (o *Store) ObjectRepository() store.ObjectRepository {
	return o.objectRepo
}

func (o *Store) PortRepository() store.PortRepository {
	return o.portRepo
}

func (o *Store) DeviceRepository() store.DeviceRepository {
	return o.deviceRepo
}

func (o *Store) ScriptRepository() store.ScriptRepository {
	return o.scriptRepo
}

// AR

func (o *Store) EventsRepo() store.EventsRepo {
	return o.eventsRepo
}

func (o *Store) EventActionsRepo() store.EventActionsRepo {
	return o.eventActionsRepo
}

func (o *Store) CronRepo() store.CronRepo {
	return o.cronRepo
}

// TR

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
