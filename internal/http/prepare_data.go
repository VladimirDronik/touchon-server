package http

import (
	"encoding/json"

	"github.com/pkg/errors"
	"touchon-server/internal/model"
	"touchon-server/internal/store"
)

// outputZoneItems Формирование JSON с кнопками всех остальных комнат для отправки на сервер
func (o *Server) outputZoneItems() ([]*model.GroupRoom, error) {
	groupZoneRep, err := store.I.Items().GetZoneItems()
	if err != nil {
		return nil, errors.Wrap(err, "outputZoneItems")
	}

	for _, group := range groupZoneRep {
		for _, item := range group.Items {
			if err := o.prepareItem(item); err != nil {
				return nil, errors.Wrap(err, "outputZoneItems")
			}
		}
	}

	return groupZoneRep, nil
}

// outputZoneItem Формирование JSON с кнопками для помещения
func (o *Server) outputZoneItem(id int) (*model.GroupRoom, error) {
	var zone, err = store.I.Items().GetZone(id)
	if err != nil {
		return nil, errors.Wrap(err, "outputZoneItem")
	}

	for _, item := range zone.Items {
		if err := o.prepareItem(item); err != nil {
			return nil, errors.Wrap(err, "outputZoneItem")
		}
	}

	return zone, nil
}

// Подготовка итемов для вывода в правильном формате, который требуется для itemView
func (o *Server) prepareItem(item *model.ViewItem) error {
	var err error

	switch item.Type {
	case "group":
		item.GroupElements, err = store.I.Items().GetGroupElements(item.ID)
		if err != nil {
			return errors.Wrap(err, "prepareItem")
		}
	}

	item.ParamsOutput = json.RawMessage(item.Params)

	return nil
}
