package http

import (
	"sort"

	"github.com/pkg/errors"
	"touchon-server/internal/model"
	"touchon-server/internal/store"
)

func loadChildren(m map[int]*model.StoreObject, rows []*model.StoreObject, db store.ObjectRepository, age int, childType model.ChildType, withTags bool) error {
	if len(rows) == 0 || age < 1 {
		return nil
	}

	ids := make([]int, 0, len(rows))
	for _, row := range rows {
		ids = append(ids, row.ID)
	}

	// Получаем всех детей
	children, err := db.GetObjectChildren(childType, ids...)
	if err != nil {
		return errors.Wrap(err, "loadChildren")
	}

	// Добавляем детей в общий список
	for _, row := range children {
		if withTags == false {
			row.Tags = nil
		}
		m[row.ID] = row
	}

	// Пытаемся загрузить следующее поколение детей
	return loadChildren(m, children, db, age-1, childType, withTags)
}

func loadParents(m map[int]*model.StoreObject, rows []*model.StoreObject, db store.ObjectRepository) error {
	// Собираем ID родителей, которых нет в полученном списке
	parentIDsMap := make(map[int]bool, len(rows))
	for _, row := range rows {
		if row.ParentID == nil {
			continue
		}

		if _, ok := m[*row.ParentID]; !ok {
			parentIDsMap[*row.ParentID] = true
		}
	}

	if len(parentIDsMap) == 0 {
		return nil
	}

	parentIDs := make([]int, 0, len(parentIDsMap))
	for id := range parentIDsMap {
		parentIDs = append(parentIDs, id)
	}

	// Получаем отсутствующих родителей
	rows, err := db.GetObjectsByIDs(parentIDs)
	if err != nil {
		return errors.Wrap(err, "loadParents")
	}

	// Добавляем новых родителей в общий список
	for _, row := range rows {
		m[row.ID] = row
	}

	// В списке новых родителей могут быть их родители отсутствующие в списке
	return loadParents(m, rows, db)
}

func sortObjectsTree(items []*model.StoreObject) {
	sort.Slice(items, func(i, j int) bool {
		switch {
		case items[i].Category != items[j].Category:
			return items[i].Category < items[j].Category
		case items[i].Type != items[j].Type:
			return items[i].Type < items[j].Type
		default:
			return items[i].Name < items[j].Name
		}
	})

	for _, item := range items {
		if len(item.Children) > 0 {
			sortObjectsTree(item.Children)
		}

		item.ParentID = nil
	}
}
