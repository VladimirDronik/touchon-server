package memstore

import (
	"sync"

	"github.com/pkg/errors"
	"touchon-server/internal/context"
	"touchon-server/internal/model"
	"touchon-server/internal/objects"
	"touchon-server/internal/store"
)

// Global instance
var I *MemStore

func New() (*MemStore, error) {
	rows, err := store.I.ObjectRepository().GetObjects(map[string]interface{}{"parent_id": nil}, nil, 0, 10000, model.ChildTypeAll)
	if err != nil {
		return nil, errors.Wrap(err, "memstore.New")
	}

	tree := make(map[int]objects.Object, len(rows))

	for _, row := range rows {
		obj, err := objects.LoadObject(row.ID, "", "", model.ChildTypeAll)
		if err != nil {
			return nil, errors.Wrap(err, "memstore.New")
		}

		tree[obj.GetID()] = obj
	}

	list := make(map[int]objects.Object, 10000)
	for _, obj := range tree {
		walk(obj, list)
	}
	context.Logger.Infof("Объекты: корневые %d, всего %d", len(tree), len(list))

	return &MemStore{
		mu:      sync.RWMutex{},
		objects: list,
	}, nil
}

func walk(obj objects.Object, list map[int]objects.Object) {
	list[obj.GetID()] = obj
	for _, child := range obj.GetChildren().GetAll() {
		walk(child, list)
	}
}

type MemStore struct {
	mu      sync.RWMutex
	objects map[int]objects.Object
}

type treeItem struct {
	Object   objects.Object
	Children []*treeItem
}

func (o *MemStore) Start() error {
	o.mu.Lock()
	defer o.mu.Unlock()

	// Строим дерево, чтобы сначала запускались родительские объекты, и только после них дочерние
	tree := make(map[int]*treeItem, len(o.objects))
	for _, obj := range o.objects {
		tree[obj.GetID()] = &treeItem{Object: obj}
	}

	for _, item := range tree {
		if parentID := item.Object.GetParentID(); parentID != nil {
			parent, ok := tree[*parentID]
			if !ok {
				return errors.Wrap(errors.Errorf("parentID %d not found for item %d", *parentID, item.Object.GetID()), "Start")
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

	if err := startTree(list); err != nil {
		return errors.Wrap(err, "Start")
	}

	return nil
}

func startTree(items []*treeItem) error {
	for _, item := range items {
		if err := item.Object.Start(); err != nil {
			return errors.Wrap(err, "Start")
		}

		if len(item.Children) > 0 {
			if err := startTree(item.Children); err != nil {
				return err
			}
		}
	}

	return nil
}

func (o *MemStore) Shutdown() error {
	o.mu.Lock()
	defer o.mu.Unlock()

	errs := make([]error, 0, len(o.objects))
	for _, obj := range o.objects {
		if err := obj.Shutdown(); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return errors.Wrap(errs[0], "Shutdown")
	}

	return nil
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
		return nil, errors.Wrap(errors.New("object not found"), "GetObjectUnsafe")
	}

	return obj, nil
}

func (o *MemStore) GetObjects(params map[string]interface{}, offset, limit int, objType model.ChildType) ([]*model.StoreObject, int, error) {
	type Filters struct {
		ID       int
		ParentID int
		ZoneID   int
		Category string
		Type     string
		Name     string
		Status   string
	}
	filters := &Filters{}
	var ok bool

	for k, v := range params {
		switch k {
		case "id":
			filters.ID, ok = v.(int)
			if !ok {
				return nil, 0, errors.Wrap(errors.Errorf("id is not int (%T)", v), "GetObjects")
			}

		case "parent_id":
			filters.ParentID, ok = v.(int)
			if !ok {
				return nil, 0, errors.Wrap(errors.Errorf("parent_id is not int (%T)", v), "GetObjects")
			}

		case "zone_id":
			filters.ZoneID, ok = v.(int)
			if !ok {
				return nil, 0, errors.Wrap(errors.Errorf("zone_id is not int (%T)", v), "GetObjects")
			}

		case "category":
			filters.Category, ok = v.(string)
			if !ok {
				return nil, 0, errors.Wrap(errors.Errorf("category is not string (%T)", v), "GetObjects")
			}

		case "type":
			filters.Type, ok = v.(string)
			if !ok {
				return nil, 0, errors.Wrap(errors.Errorf("type is not string (%T)", v), "GetObjects")
			}

		case "name":
			filters.Name, ok = v.(string)
			if !ok {
				return nil, 0, errors.Wrap(errors.Errorf("name is not string (%T)", v), "GetObjects")
			}

		case "status":
			filters.Status, ok = v.(string)
			if !ok {
				return nil, 0, errors.Wrap(errors.Errorf("status is not string (%T)", v), "GetObjects")
			}
		}
	}

	o.mu.Lock()
	defer o.mu.Unlock()

	r := make([]*model.StoreObject, 0, len(o.objects))
	for _, obj := range o.objects {
		parentID := obj.GetParentID()

		equal := true
		switch {
		case filters.ID > 0 && obj.GetID() != filters.ID:
			equal = false
		case filters.ParentID > 0 && parentID != nil && *parentID != filters.ParentID:
			equal = false
		case filters.ZoneID > 0 && obj.GetZoneID() != nil && *obj.GetZoneID() != filters.ZoneID:
			equal = false
		case filters.Category != "" && string(obj.GetCategory()) != filters.Category:
			equal = false
		case filters.Type != "" && obj.GetType() != filters.Type:
			equal = false
		case filters.Name != "" && obj.GetName() != filters.Name:
			equal = false
		case filters.Status != "" && string(obj.GetStatus()) != filters.Status:
			equal = false
		case objType == model.ChildTypeInternal && !obj.GetInternal():
			equal = false
		case objType == model.ChildTypeExternal && obj.GetInternal():
			equal = false
		}

		if equal {
			r = append(r, obj.GetStoreObject())
		}
	}

	total := len(r)

	if offset > len(r) {
		return nil, total, nil
	}

	r = r[offset:]

	if limit > len(r) {
		return r, total, nil
	}

	return r[:limit], total, nil
}

func (o *MemStore) GetObjectChildren(childType model.ChildType, objectIDs ...int) ([]*model.StoreObject, error) {
	objectIDsMap := make(map[int]bool, len(objectIDs))
	for _, id := range objectIDs {
		objectIDsMap[id] = true
	}

	o.mu.Lock()
	defer o.mu.Unlock()

	r := make([]*model.StoreObject, 0, len(o.objects))
	for _, obj := range o.objects {
		parentID := obj.GetParentID()
		if parentID == nil {
			continue
		}

		if _, ok := objectIDsMap[*parentID]; !ok {
			continue
		}

		if childType == model.ChildTypeInternal && !obj.GetInternal() {
			continue
		}

		if childType == model.ChildTypeExternal && obj.GetInternal() {
			continue
		}

		r = append(r, obj.GetStoreObject())
	}

	return r, nil
}

func (o *MemStore) GetObjectsByIDs(ids []int) ([]*model.StoreObject, error) {
	objectIDsMap := make(map[int]bool, len(ids))
	for _, id := range ids {
		objectIDsMap[id] = true
	}

	o.mu.Lock()
	defer o.mu.Unlock()

	r := make([]*model.StoreObject, 0, len(o.objects))
	for _, obj := range o.objects {
		if _, ok := objectIDsMap[obj.GetID()]; ok {
			r = append(r, obj.GetStoreObject())
		}
	}

	return r, nil
}

func (o *MemStore) SaveObject(obj objects.Object) error {
	o.mu.Lock()
	defer o.mu.Unlock()

	if _, ok := o.objects[obj.GetID()]; ok {
		if err := obj.Shutdown(); err != nil {
			return errors.Wrap(err, "SaveObject")
		}
	}

	o.objects[obj.GetID()] = obj

	if err := obj.Start(); err != nil {
		return errors.Wrap(err, "SaveObject")
	}

	return nil
}

func (o *MemStore) DeleteObject(objectID int) error {
	o.mu.Lock()
	defer o.mu.Unlock()

	if obj, ok := o.objects[objectID]; ok {
		delete(o.objects, objectID)

		if err := obj.Shutdown(); err != nil {
			return errors.Wrap(err, "DeleteObject")
		}
	}

	return nil
}
