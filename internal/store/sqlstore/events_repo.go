package sqlstore

import (
	"github.com/pkg/errors"
	"touchon-server/internal/store"
	"touchon-server/lib/interfaces"
)

type EventsRepo struct {
	store *Store
}

// GetEvents возвращает события сущности.
func (o *EventsRepo) GetEvents(targetType interfaces.TargetType, targetID int) ([]*interfaces.AREvent, error) {
	rows := make([]*interfaces.AREvent, 0, 10)

	err := o.store.db.
		Where("target_type = ?", targetType).
		Where("target_id = ?", targetID).
		Find(&rows).Error

	if err != nil {
		return nil, errors.Wrap(err, "GetEvents")
	}

	return rows, nil
}

// GetEvent возвращает событие.
func (o *EventsRepo) GetEvent(targetType interfaces.TargetType, targetID int, eventName string) (*interfaces.AREvent, error) {
	row := &interfaces.AREvent{}

	err := o.store.db.
		Where("target_type = ?", targetType).
		Where("target_id = ?", targetID).
		Where("event_name = ?", eventName).
		Find(row).Error

	if err != nil {
		return nil, errors.Wrap(err, "GetEvent")
	}

	if row.ID == 0 {
		return nil, errors.Wrap(store.ErrNotFound, "GetEvent")
	}

	return row, nil
}

// SaveEvent добавляет или обновляет событие.
func (o *EventsRepo) SaveEvent(event *interfaces.AREvent) error {
	if event == nil {
		return errors.Wrap(errors.New("event is nil"), "SaveEvent")
	}

	if err := o.store.db.Save(event).Error; err != nil {
		return errors.Wrap(err, "SaveEvent")
	}

	return nil
}

// DeleteEvent удаляет событие.
func (o *EventsRepo) DeleteEvent(targetType interfaces.TargetType, targetID int, eventName string) error {
	var sqlString string
	// foreign_keys - для каскадного удаления записей

	if eventName == "all" {
		sqlString = "DELETE FROM ar_events WHERE target_type = ? AND target_id = ?"
	} else {
		sqlString = "DELETE FROM ar_events WHERE target_type = ? AND target_id = ? AND event_name = ?"
	}

	err := o.store.db.
		Exec(sqlString, targetType, targetID, eventName).
		Error

	if err != nil {
		return errors.Wrap(err, "DeleteEvent")
	}

	return nil
}
