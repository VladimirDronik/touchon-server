package sqlstore

import (
	"encoding/json"

	"github.com/pkg/errors"
	"touchon-server/lib/interfaces"
)

type EventActionsRepo struct {
	store *Store
}

// GetActions возвращает список действий для события.
func (o *EventActionsRepo) GetActions(eventIDs ...int) (map[int][]*interfaces.EventAction, error) {
	type Row struct {
		interfaces.EventAction
		Args string
	}
	rows := make([]*Row, 0, 10)

	err := o.store.db.Model(&interfaces.EventAction{}).
		Where("event_id in ?", eventIDs).
		Order("event_id, sort, id").
		Scan(&rows).Error

	if err != nil {
		return nil, errors.Wrap(err, "GetActions")
	}

	r := make(map[int][]*interfaces.EventAction, len(rows))
	for _, row := range rows {
		if err := json.Unmarshal([]byte(row.Args), &row.EventAction.Args); err != nil {
			return nil, errors.Wrap(err, "GetActions")
		}

		r[row.EventID] = append(r[row.EventID], &row.EventAction)
	}

	return r, nil
}

// GetActionsCount возвращает количество действий для событий.
func (o *EventActionsRepo) GetActionsCount(eventIDs ...int) (map[int]int, error) {
	type Row struct {
		EventID int `gorm:"event_id"`
		Count   int `gorm:"count"`
	}
	rows := make([]*Row, 0, len(eventIDs))

	err := o.store.db.
		Select("event_id, count(id) as count").
		Where("event_id in ?", eventIDs).
		Table("ar_event_actions").
		Group("event_id").
		Find(&rows).Error

	if err != nil {
		return nil, errors.Wrap(err, "GetActionsCount")
	}

	m := make(map[int]int, len(rows))
	for _, row := range rows {
		m[row.EventID] = row.Count
	}

	return m, nil
}

// SaveAction добавляет или обновляет действие.
func (o *EventActionsRepo) SaveAction(act *interfaces.EventAction) error {
	if act == nil {
		return errors.Wrap(errors.New("act is nil"), "SaveAction")
	}

	count := int64(0)
	if err := o.store.db.Model(act).Where("id = ?", act.ID).Count(&count).Error; err != nil {
		return errors.Wrap(err, "SaveAction")
	}
	actIsExists := count == 1

	type Row struct {
		*interfaces.EventAction
		Args string
	}

	args, err := json.Marshal(act.Args)
	if err != nil {
		return errors.Wrap(err, "SaveAction")
	}

	row := &Row{
		EventAction: act,
		Args:        string(args),
	}

	if actIsExists {
		if err := o.store.db.Updates(row).Error; err != nil {
			return errors.Wrap(err, "SaveAction(update)")
		}
	} else {
		if err := o.store.db.Create(row).Error; err != nil {
			return errors.Wrap(err, "SaveAction(create)")
		}
	}

	return nil
}

// DeleteAction удаляет действие.
func (o *EventActionsRepo) DeleteAction(actID int) error {
	// foreign_keys - для каскадного удаления записей
	err := o.store.db.Exec("DELETE FROM ar_event_actions WHERE id = ?", actID).Error
	if err != nil {
		return errors.Wrap(err, "DeleteAction")
	}

	return nil
}

// DeleteAction удаляет действие.
func (o *EventActionsRepo) DeleteActionByObject(targetType string, objectID int) error {
	// foreign_keys - для каскадного удаления записей
	err := o.store.db.Exec("DELETE FROM ar_event_actions WHERE target_type = ? AND target_id = ?", targetType, objectID).Error
	if err != nil {
		return errors.Wrap(err, "DeleteAction")
	}

	return nil
}

// OrderActions меняет порядок действий.
func (o *EventActionsRepo) OrderActions(actIDs []int) error {
	for i, actID := range actIDs {
		err := o.store.db.Model(&interfaces.EventAction{}).
			Where("id = ?", actID).
			Update("sort", i).Error

		if err != nil {
			return errors.Wrap(err, "OrderActions")
		}
	}

	return nil
}
