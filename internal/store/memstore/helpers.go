package memstore

import (
	"github.com/pkg/errors"
	"touchon-server/internal/objects"
)

type treeItem struct {
	Object   objects.Object
	Children []*treeItem
}

func (o *MemStore) makeTree() ([]*treeItem, error) {
	tree := make(map[int]*treeItem, len(o.objects))
	for _, obj := range o.objects {
		tree[obj.GetID()] = &treeItem{Object: obj}
	}

	for _, item := range tree {
		if parentID := item.Object.GetParentID(); parentID != nil {
			parent, ok := tree[*parentID]
			if !ok {
				return nil, errors.Wrap(errors.Errorf("parentID %d not found for item %d", *parentID, item.Object.GetID()), "Start")
			}

			parent.Children = append(parent.Children, item)
		}
	}

	// Удаляем все элементы, кроме корневых
	list := make([]*treeItem, 0, len(tree))
	for _, item := range tree {
		if parentID := item.Object.GetParentID(); parentID == nil {
			list = append(list, item)
		}
	}

	return list, nil
}

func addObjectsToList(obj objects.Object, list map[int]objects.Object) {
	list[obj.GetID()] = obj
	for _, child := range obj.GetChildren().GetAll() {
		addObjectsToList(child, list)
	}
}

func (o *MemStore) startObjectTree(obj objects.Object) error {
	// Если объект выключен или уже запущен, уходим
	if !obj.GetEnabled() || obj.GetStarted() {
		return nil
	}

	if err := obj.Start(); err != nil {
		return errors.Wrap(err, "startObjectTree")
	}

	for _, child := range obj.GetChildren().GetAll() {
		if err := o.startObjectTree(child); err != nil {
			return err
		}
	}

	return nil
}

func (o *MemStore) shutdownObjectTree(obj objects.Object) (errs []error) {
	for _, child := range obj.GetChildren().GetAll() {
		if childrenErrs := o.shutdownObjectTree(child); len(childrenErrs) > 0 {
			errs = append(errs, childrenErrs...)
		}
	}

	// Останавливаем только запущенные объекты
	if !obj.GetStarted() {
		return errs
	}

	if err := obj.Shutdown(); err != nil {
		errs = append(errs, errors.Wrap(err, "shutdownObjectTree"))
	}

	return errs
}
