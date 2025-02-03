package sqlstore

import (
	"github.com/pkg/errors"
	"touchon-server/internal/model"
	"touchon-server/internal/store"
	"touchon-server/lib/mqtt/messages"
)

type EventsRepo struct {
	store *Store
}

// GetEvents возвращает события сущности.
func (o *EventsRepo) GetEvents(targetType messages.TargetType, targetID int) ([]*model.Event, error) {
	rows := make([]*model.Event, 0, 10)

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
func (o *EventsRepo) GetEvent(targetType messages.TargetType, targetID int, eventName string) (*model.Event, error) {
	row := &model.Event{}

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
func (o *EventsRepo) SaveEvent(event *model.Event) error {
	if event == nil {
		return errors.Wrap(errors.New("event is nil"), "SaveEvent")
	}

	count := int64(0)
	if err := o.store.db.Model(event).Where("id = ?", event.ID).Count(&count).Error; err != nil {
		return errors.Wrap(err, "SaveEvent")
	}
	itemIsExists := count == 1

	if itemIsExists {
		if err := o.store.db.Updates(event).Error; err != nil {
			return errors.Wrap(err, "SaveEvent(update)")
		}
	} else {
		if err := o.store.db.Create(event).Error; err != nil {
			return errors.Wrap(err, "SaveEvent(create)")
		}
	}

	return nil
}

// DeleteEvent удаляет событие.
func (o *EventsRepo) DeleteEvent(targetType messages.TargetType, targetID int, eventName string) error {
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

// GetAllEventsName возвращает названия всех событий, используемых в таблице.
// Используется для проверки правильности указанных имен.
func (o *EventsRepo) GetAllEventsName() ([]string, error) {
	type Row struct {
		EventName string
	}
	rows := make([]*Row, 0, 100)

	err := o.store.db.
		Table("ar_events").
		Distinct("event_name").
		Where("enabled").
		Find(&rows).Error

	if err != nil {
		return nil, errors.Wrap(err, "GetAllEventsName")
	}

	r := make([]string, 0, len(rows))
	for _, row := range rows {
		r = append(r, row.EventName)
	}

	return r, nil
}
