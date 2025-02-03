package sqlstore

import (
	"time"

	"github.com/pkg/errors"
	"translator/internal/helpers"
	"translator/internal/model"
)

type Users struct {
	store *Store
}

// AddRefreshToken добавление refresh_token в таблицу пользователей, если пользователь существует и вывод ошибки, если
// такого пользователя нет
func (o *Users) AddRefreshToken(deviceID int, refreshToken string, refreshTokenTTL time.Duration) error {
	err := o.store.db.Table("users").
		Where("device_id = ?", deviceID).
		Update("refresh_token", refreshToken).
		Update("token_expired", time.Now().Add(refreshTokenTTL)).Error

	if err != nil {
		return errors.Wrap(err, "AddRefreshToken")
	}

	return nil
}

// GetByToken Получаем юзера по его токену, одновременно проверяя не истек ли он
func (o *Users) GetByToken(token string) (*model.User, error) {
	r := &model.User{}

	err := o.store.db.Where("refresh_token = ?", token).
		Where("token_expired > ?", time.Now()).
		First(r).Error

	if err != nil {
		return nil, errors.Wrap(err, "GetByToken")
	}

	return r, nil
}

// GetByLoginAndPassword получить юзера по связке логин/пароль
func (o *Users) GetByLoginAndPassword(login string, password string) (*model.User, error) {
	r := &model.User{}

	if err := o.store.db.Where("login = ? AND password = ?", login, helpers.MD5(password)).First(r).Error; err != nil {
		return nil, errors.Wrap(err, "GetByLoginAndPassword")
	}

	return r, nil
}

// GetByToken Получаем юзера по его токену, одновременно проверяя не истек ли он
func (o *Users) GetByDeviceID(deviceID int) (*model.User, error) {
	r := &model.User{}

	if err := o.store.db.Where("device_id = ?", deviceID).First(r).Error; err != nil {
		return nil, errors.Wrap(err, "GetByDeviceID")
	}

	return r, nil
}

// GetByToken Получаем юзера по его токену, одновременно проверяя не истек ли он
func (o *Users) GetAllUsers() ([]*model.User, error) {
	var r []*model.User

	err := o.store.db.Select("*").
		Find(&r).Error

	if err != nil {
		return nil, errors.Wrap(err, "GetAllUsers")
	}

	return r, nil
}

func (o *Users) Create(user *model.User) (int, error) {
	user.Password = helpers.MD5(user.Password)

	err := o.store.db.Table("users").Create(user).Error
	if err != nil {
		return 0, err
	}

	return user.ID, nil
}

func (o *Users) Delete(userID int) error {
	return o.store.db.Where("id = ?", userID).Delete(&model.User{}).Error
}

// RemoveToken Удаление данных о сессии в таблице токенов
func (o *Users) RemoveToken(refreshToken string) error {
	if err := o.store.db.Where("refresh_token = ?", refreshToken).Delete(&model.Tokens{}).Error; err != nil {
		return errors.Wrap(err, "RemoveToken")
	}

	return nil
}

// LinkDeviceToken указать токен устройства для отправки push уведомлений
func (o *Users) LinkDeviceToken(userID int, token string, deviceType string) error {
	var user *model.User

	if err := o.store.db.Where("device_id = ?", userID).First(&user).Error; err != nil {
		return errors.Wrap(err, "LinkDeviceToken")
	}

	user.DeviceToken = token
	user.DeviceType = deviceType

	if err := o.store.db.Save(&user).Error; err != nil {
		return errors.Wrap(err, "LinkDeviceToken")
	}

	return nil
}

// GetDeviceToken получить токен устройства и его тип для уведомлений
func (o *Users) GetDeviceToken(deviceID int) (*model.DeviceTokens, error) {
	user := &model.User{}

	if err := o.store.db.Where("device_id = ?", deviceID).Where("device_token != '' AND device_type != ''").First(user).Error; err != nil {
		return nil, errors.Wrap(err, "LinkDeviceToken")
	}

	return &user.DeviceTokens, nil
}
