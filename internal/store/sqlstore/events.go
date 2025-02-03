package sqlstore

import (
	"translator/internal/model"
)

type Events struct {
	store *Store
}

func (e *Events) AddEvent(event *model.Event) error {
	return e.store.db.Create(&event).Error
}

func (e *Events) UpdateEvent(event *model.Event) error {
	return e.store.db.Save(&event).Error
}

func (e *Events) DeleteEvent(itemID int) error {
	event := model.Event{}
	return e.store.db.Where("item_id = ?", itemID).Delete(&event).Error
}
