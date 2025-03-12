package helpers

import (
	"github.com/pkg/errors"
	"touchon-server/internal/g"
	"touchon-server/internal/model"
	"touchon-server/internal/objects"
	"touchon-server/internal/store"
	"touchon-server/internal/store/memstore"
)

// ResetParentAndAddress Убирает у объекта родителя и адрес в свойствах
func ResetParentAndAddress(objectsToReset map[int]string) error {
	for objectID, _ := range objectsToReset {
		if err := store.I.ObjectRepository().SetProp(objectID, "address", "0"); err != nil {
			return errors.Wrap(err, "ResetParentAndAddress")
		}

		if err := store.I.ObjectRepository().SetParent(objectID, nil); err != nil {
			return errors.Wrap(err, "ResetParentAndAddress")
		}

		if err := store.I.ObjectRepository().SetObjectStatus(objectID, string(model.StatusDisabled)); err != nil {
			return errors.Wrap(err, "ResetParentAndAddress")
		}

		if err := store.I.ObjectRepository().SetEnabled(objectID, false); err != nil {
			return errors.Wrap(err, "ResetParentAndAddress")
		}

		if err := memstore.I.DisableObject(objectID); err != nil {
			return errors.Wrap(err, "ResetParentAndAddress")
		}
	}

	return nil
}

// SaveAndSendStatus Установка статуса объекту и отправка этого статуса в вебсокеты
func SaveAndSendStatus(obj objects.Object, status model.ObjectStatus) error {
	// Выставляем статус у объекта в памяти
	obj.SetStatus(status)

	// Сохраняем статус в базе
	if err := store.I.ObjectRepository().SetStatus(obj.GetID(), status); err != nil {
		return errors.Wrap(err, "SaveAndSendStatus")
	}

	wsMsg := model.ObjectForWS{
		ID:     obj.GetID(),
		Status: obj.GetStatus(),
	}
	g.WSServer.Send("object", wsMsg)

	return nil
}
