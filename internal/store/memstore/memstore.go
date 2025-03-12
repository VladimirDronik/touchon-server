package memstore

import (
	"sync"

	"github.com/pkg/errors"
	"touchon-server/internal/g"
	"touchon-server/internal/objects"
	"touchon-server/internal/store"
)

// Global instance
var I *MemStore

func New() (*MemStore, error) {
	rows, err := store.I.ObjectRepository().GetObjects(map[string]interface{}{"parent_id": nil}, nil, 0, 10000)
	if err != nil {
		return nil, errors.Wrap(err, "memstore.New")
	}

	tree := make(map[int]objects.Object, len(rows))

	for _, row := range rows {
		obj, err := objects.LoadObject(row.ID, "", "", true)
		if err != nil {
			return nil, errors.Wrap(err, "memstore.New")
		}

		tree[obj.GetID()] = obj
	}

	list := make(map[int]objects.Object, 10000)
	for _, obj := range tree {
		addObjectsToList(obj, list)
	}
	g.Logger.Infof("Объекты: корневые %d, всего %d", len(tree), len(list))

	return &MemStore{
		mu:      sync.RWMutex{},
		objects: list,
	}, nil
}

type MemStore struct {
	mu      sync.RWMutex
	objects map[int]objects.Object
}

func (o *MemStore) GetObject(objectID int) (objects.Object, error) {
	o.mu.Lock()
	defer o.mu.Unlock()

	obj, err := o.GetObjectUnsafe(objectID)
	if err != nil {
		return nil, errors.Wrap(err, "GetObject")
	}

	return obj, nil
}

// GetObjectUnsafe возвращает объект не устанавливая блокировку на хранилище
func (o *MemStore) GetObjectUnsafe(objectID int) (objects.Object, error) {
	obj, ok := o.objects[objectID]
	if !ok {
		return nil, errors.Wrapf(errors.New("object not found"), "GetObjectUnsafe(%d)", objectID)
	}

	return obj, nil
}

func (o *MemStore) DeleteObject(objectID int) error {
	o.mu.Lock()
	defer o.mu.Unlock()

	// full=false - оставляем приемных детей (например, сенсоры у контроллера)
	if err := o.deleteObject(objectID, false); err != nil {
		return errors.Wrap(err, "DeleteObject")
	}

	return nil
}

// Удаляет объект и его детей из хранилища, удаляем ссылку на него у родителя
func (o *MemStore) deleteObject(objectID int, full bool) error {
	if obj, ok := o.objects[objectID]; ok {
		// Останавливаем объект с детьми
		errs := o.shutdownObjectTree(obj)

		// Удаляем его у родителя
		if err := o.deleteFromParent(obj); err != nil {
			errs = append(errs, err)
		}

		// Удаляем детей
		o.deleteChildren(obj, full)

		// Удаляем объект из общего списка
		delete(o.objects, objectID)

		// Выводим в лог все ошибки
		for _, err := range errs {
			g.Logger.Error(errors.Wrap(err, "deleteObject"))
		}

		// Возвращаем первую
		if len(errs) > 0 {
			return errors.Wrap(errs[0], "deleteObject")
		}
	}

	return nil
}

// EnableObject включает объект и запускает его (метод Start())
func (o *MemStore) EnableObject(objectID int) error {
	o.mu.Lock()
	defer o.mu.Unlock()

	if obj, ok := o.objects[objectID]; ok {
		obj.SetEnabled(true)

		if err := o.startObjectTree(obj); err != nil {
			return errors.Wrap(err, "EnableObject")
		}
	}

	return nil
}

// DisableObject останавливает объект и отключает его запуск
func (o *MemStore) DisableObject(objectID int) error {
	o.mu.Lock()
	defer o.mu.Unlock()

	if obj, ok := o.objects[objectID]; ok {
		errs := o.shutdownObjectTree(obj)
		obj.SetEnabled(false)

		// Выводим в лог все ошибки
		for _, err := range errs {
			g.Logger.Error(errors.Wrap(err, "DisableObject"))
		}

		// Возвращаем первую
		if len(errs) > 0 {
			return errors.Wrap(errs[0], "DisableObject")
		}
	}

	return nil
}

func (o *MemStore) Search(f func(items map[int]objects.Object) ([]objects.Object, error)) ([]objects.Object, error) {
	o.mu.Lock()
	defer o.mu.Unlock()

	rows, err := f(o.objects)
	if err != nil {
		return nil, errors.Wrap(err, "MemStore.Search")
	}

	return rows, nil
}

func (o *MemStore) addToParent(obj objects.Object) error {
	if obj.GetParentID() == nil {
		return nil
	}

	parent, err := o.GetObjectUnsafe(*obj.GetParentID())
	if err != nil {
		return errors.Wrap(err, "deleteFromParent")
	}

	parent.GetChildren().Add(obj)

	return nil
}

func (o *MemStore) deleteFromParent(obj objects.Object) error {
	if obj.GetParentID() == nil {
		return nil
	}

	parent, err := o.GetObjectUnsafe(*obj.GetParentID())
	if err != nil {
		return errors.Wrap(err, "deleteFromParent")
	}

	parent.GetChildren().DeleteByID(obj.GetID())

	return nil
}

func (o *MemStore) deleteChildren(obj objects.Object, all bool) {
	for _, child := range obj.GetChildren().GetAll() {
		if child.GetChildren().Len() > 0 {
			o.deleteChildren(child, all)
		}

		if !all {
			// TODO удалять только родных детей (например, у контроллера удалять только порты)
			// TODO у оставшихся детей выставлять enabled=false, parent_id=nil
		}

		// Удаляем объект из общего списка
		delete(o.objects, child.GetID())
	}

	// Удаляем ссылки на объекты, для уменьшения вероятности утечки памяти
	obj.GetChildren().DeleteAll()
}

func (o *MemStore) SaveObject(obj objects.Object) error {
	o.mu.Lock()
	defer o.mu.Unlock()

	// Удаляем объект и его детей из хранилища, удаляем ссылку на него у родителя
	if err := o.deleteObject(obj.GetID(), true); err != nil {
		return errors.Wrap(err, "SaveObject")
	}

	// Добавляем объект в общий список
	o.objects[obj.GetID()] = obj

	// Добавляем объект к родителю
	if err := o.addToParent(obj); err != nil {
		return errors.Wrap(err, "SaveObject")
	}

	// Добавляем всех детей в общий список
	addObjectsToList(obj, o.objects)

	// Запускаем объект и его детей
	if err := o.startObjectTree(obj); err != nil {
		return errors.Wrap(err, "SaveObject")
	}

	return nil
}
