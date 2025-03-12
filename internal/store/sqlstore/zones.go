package sqlstore

import (
	"sort"

	"github.com/pkg/errors"
	"touchon-server/internal/model"
)

type Zones struct {
	store *Store
}

// CreateZone Создание нового помещения
func (z *Zones) CreateZone(zone *model.Zone) (int, error) {
	if err := z.store.db.Create(zone).Error; err != nil {
		return 0, err
	}
	return zone.ID, nil
}

// GetZoneTrees получение всех помещений
func (z *Zones) GetZoneTrees(typeZones string, parentIDs ...int) ([]*model.Zone, error) {
	m := make(map[int]*model.Zone, len(parentIDs))

	var rows []*model.Zone
	q := z.store.db

	if parentIDs != nil {
		q = q.Where("id in ?", parentIDs)
	}

	if typeZones == "groups_only" {
		q = q.Where("is_group")
	}

	if err := q.Find(&rows).Error; err != nil {
		return nil, errors.Wrap(err, "GetZoneTrees")
	}

	for _, row := range rows {
		if _, ok := m[row.ID]; !ok {
			m[row.ID] = row
		}
	}

	if err := z.getChildren(m, parentIDs); err != nil {
		return nil, errors.Wrap(err, "GetZoneTrees")
	}

	for _, item := range m {
		if item.ParentID == nil {
			continue
		}

		if parent, ok := m[*item.ParentID]; ok {
			parent.Children = append(parent.Children, item)
		}
	}

	roots := 0
	for _, item := range m {
		if len(item.Children) > 0 {
			sort.Slice(item.Children, func(i, j int) bool {
				switch {
				case item.Children[i].Sort != item.Children[j].Sort:
					return item.Children[i].Sort < item.Children[j].Sort
				default:
					return item.Children[i].Name < item.Children[j].Name
				}
			})
		}

		if item.ParentID == nil {
			roots += 1
		}
	}

	r := make([]*model.Zone, 0, roots)
	for _, item := range m {
		if item.ParentID == nil {
			r = append(r, item)
		}
	}

	sort.Slice(r, func(i, j int) bool {
		switch {
		case r[i].Sort != r[j].Sort:
			return r[i].Sort < r[j].Sort
		default:
			return r[i].Name < r[j].Name
		}
	})

	return r, nil
}

func (z *Zones) getChildren(m map[int]*model.Zone, parentIDs []int) error {
	var rows []*model.Zone

	if err := z.store.db.Where("parent_id in ?", parentIDs).Find(&rows).Error; err != nil {
		return errors.Wrap(err, "getChildren")
	}

	parentIDs = parentIDs[:0]
	for _, row := range rows {
		if _, ok := m[row.ID]; !ok {
			m[row.ID] = row
			parentIDs = append(parentIDs, row.ID)
		}
	}

	if len(parentIDs) > 0 {
		if err := z.getChildren(m, parentIDs); err != nil {
			return err
		}
	}

	return nil
}

// UpdateZones сохранение модели комнат, например при изменении сортировки
func (z *Zones) UpdateZones(zones []*model.Zone) error {
	for _, item := range zones {
		if err := z.store.db.Table("zones").Updates(item).Error; err != nil {
			return errors.Wrap(err, "UpdateZones")
		}

		if len(item.Children) > 0 {
			if err := z.UpdateZones(item.Children); err != nil {
				return err
			}
		}
	}

	return nil
}

// SetOrder задает порядок сортировки
func (z *Zones) SetOrder(zoneIDs []int) error {
	for i, zoneID := range zoneIDs {
		if err := z.SetFieldValue(zoneID, "sort", i+1); err != nil {
			return errors.Wrap(err, "SetOrder")
		}
	}

	return nil
}

func (z *Zones) SetFieldValue(zoneID int, field string, value interface{}) error {
	err := z.store.db.
		Table("zones").
		Where("id = ?", zoneID).
		Update(field, value).
		Error

	if err != nil {
		return errors.Wrap(err, "SetFieldValue")
	}

	return nil
}

// DeleteZone Удаление помещения
func (z *Zones) DeleteZone(zoneID int) error {
	zone := model.Zone{}
	return z.store.db.Where("id = ?", zoneID).
		Or("parent_id = ?", zoneID).Delete(&zone).Error
}
