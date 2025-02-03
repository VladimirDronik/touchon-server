package sqlstore

import (
	"touchon-server/internal/model"
)

type Events struct {
	store *Store
}

func (e *Events) AddEvent(event *model.TrEvent) error {
	return e.store.db.Create(&event).Error
}

func (e *Events) UpdateEvent(event *model.TrEvent) error {
	return e.store.db.Save(&event).Error
}

func (e *Events) DeleteEvent(itemID int) error {
	event := model.TrEvent{}
	return e.store.db.Where("item_id = ?", itemID).Delete(&event).Error
}
