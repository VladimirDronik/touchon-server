package sqlstore

import (
	"touchon-server/internal/model"
)

type Events struct {
	store *Store
}

func (e *Events) AddEvent(event *model.TrEvent) (int, error) {
	e.store.db.Where("target_type = ?", event.TargetType).
		Where("target_id = ?", event.TargetID).
		Where("event_name = ?", event.EventName).
		First(&event)

	if event.ID == 0 {
		err := e.store.db.Create(&event).Error
		if err != nil {
			return 0, err
		}
	}

	return event.ID, nil
}

func (e *Events) AddEventAction(eventAction *model.EventActions) (int, error) {
	err := e.store.db.Create(&eventAction).Error
	if err != nil {
		return 0, err
	}
	return eventAction.ID, nil
}

func (e *Events) UpdateEvent(event *model.TrEvent) error {
	return e.store.db.Save(&event).Error
}

func (e *Events) DeleteEvent(target string, itemID int) error {
	event := model.TrEvent{}
	return e.store.db.Where("target_id = ?", itemID).
		Where("target_type = ?", target).Delete(&event).Error
}

func (e *Events) DeleteEventAction(target string, itemID int) error {
	eventAction := model.EventActions{}
	return e.store.db.Where("target_type = ?", target).
		Where("target_id = ?", itemID).Delete(&eventAction).Error
}
