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
}

// New ...
func New(db *gorm.DB) *Store {
	return &Store{
		db: db,
	}
}

// ObjectRepository Инициализация
func (s *Store) ObjectRepository() store.ObjectRepository {
	if s.objectRepo == nil {
		s.objectRepo = &ObjectRepository{store: s}
	}

	return s.objectRepo
}

// PortRepository Инициализация
func (s *Store) PortRepository() store.PortRepository {
	if s.portRepo == nil {
		s.portRepo = &PortRepository{store: s}
	}

	return s.portRepo
}

// DeviceRepository Инициализация
func (s *Store) DeviceRepository() store.DeviceRepository {
	if s.deviceRepo == nil {
		s.deviceRepo = &DeviceRepository{store: s}
	}

	return s.deviceRepo
}

// ScriptRepository Инициализация
func (s *Store) ScriptRepository() store.ScriptRepository {
	if s.scriptRepo == nil {
		s.scriptRepo = &ScriptRepository{store: s}
	}

	return s.scriptRepo
}

// AR

// EventsRepo Инициализация
func (s *Store) EventsRepo() store.EventsRepo {
	if s.eventsRepo == nil {
		s.eventsRepo = &EventsRepo{store: s}
	}

	return s.eventsRepo
}

// EventActionsRepo Инициализация
func (s *Store) EventActionsRepo() store.EventActionsRepo {
	if s.eventActionsRepo == nil {
		s.eventActionsRepo = &EventActionsRepo{store: s}
	}

	return s.eventActionsRepo
}

// CronRepo Инициализация
func (s *Store) CronRepo() store.CronRepo {
	if s.cronRepo == nil {
		s.cronRepo = &CronRepo{store: s}
	}

	return s.cronRepo
}
