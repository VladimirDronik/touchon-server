package sqlstore

import (
	"github.com/pkg/errors"
	"translator/internal/model"
)

type Notifications struct {
	store *Store
}

// GetNotifications получить уведомления
func (o *Notifications) GetNotifications(offset int, limit int) ([]*model.Notification, error) {
	var r []*model.Notification

	err := o.store.db.Select("*").
		Order("date DESC").
		Offset(offset).
		Limit(limit).
		Find(&r).Error

	if err != nil {
		return nil, errors.Wrap(err, "GetNotifications")
	}

	return r, nil
}

// SetIsRead пометить уведомление как прочитанное
func (o *Notifications) SetIsRead(id int) error {
	if err := o.SetFieldValue(id, "is_read", true); err != nil {
		return errors.Wrap(err, "SetIsRead")
	}

	return nil
}

// AddNotification добавление нового уведомления в БД
func (o *Notifications) AddNotification(n *model.Notification) error {
	if err := o.store.db.Create(n).Error; err != nil {
		return errors.Wrap(err, "AddNotification")
	}

	return nil
}

// GetUnreadNotificationsCount Получить список не прочитанных уведомлений
func (o *Notifications) GetUnreadNotificationsCount() (int, error) {
	var count int64

	if err := o.store.db.Model(&model.Notification{}).Where("is_read = ?", false).Count(&count).Error; err != nil {
		return 0, errors.Wrap(err, "GetUnreadNotificationsCount")
	}

	return int(count), nil
}

// GetPushTokens отдает токены для отправки пуш уведомлений
func (o *Notifications) GetPushTokens() (map[string]string, error) {
	var deviceTokens []*model.DeviceTokens

	err := o.store.db.Table("users").
		Select("device_token", "device_type").
		Where("send_push = true").
		Where("device_token NOT NULL").
		Find(&deviceTokens).Error

	if err != nil {
		return nil, errors.Wrap(err, "GetPushTokens")
	}

	r := make(map[string]string, len(deviceTokens))
	for _, v := range deviceTokens {
		r[v.DeviceToken] = v.DeviceType
	}

	return r, nil
}

func (o *Notifications) SetFieldValue(id int, field string, value interface{}) error {
	err := o.store.db.
		Table("notifications").
		Where("id = ?", id).
		Update(field, value).
		Error

	if err != nil {
		return errors.Wrap(err, "SetFieldValue")
	}

	return nil
}
