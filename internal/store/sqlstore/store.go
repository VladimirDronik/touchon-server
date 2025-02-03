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
