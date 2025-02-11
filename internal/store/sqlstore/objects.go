package sqlstore

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	"gorm.io/gorm"
	"touchon-server/internal/model"
)

type ObjectRepository struct {
	store *Store
}

// GetProp читает значение свойства объекта
func (o *ObjectRepository) GetProp(objectID int, code string) (string, error) {
	var propValue string

	err := o.store.db.Model(&model.StoreProp{}).
		Select("value").
		Where("object_id = ?", objectID).
		Where("code = ?", code).
		Take(&propValue).Error

	if err != nil {
		return "", errors.Wrap(err, "GetProp")
	}

	return propValue, nil
}

// SetProp устанавливает значение свойства объекта
func (o *ObjectRepository) SetProp(objectID int, code, value string) error {
	q := `INSERT INTO om_props (object_id, code, value) VALUES (?, ?, ?) ON CONFLICT(object_id, code) DO UPDATE SET value = ?;`
	if err := o.store.db.Exec(q, objectID, code, value, value).Error; err != nil {
		return errors.Wrap(err, "SetProp")
	}

	return nil
}

// DelProp удаляет свойство объекта
func (o *ObjectRepository) DelProp(objectID int, code string) error {
	err := o.store.db.
		Where("object_id = ?", objectID).
		Where("code = ?", code).
		Delete(&model.StoreProp{}).Error

	if err != nil {
		return errors.Wrap(err, "DelProp")
	}

	return nil
}

func (o *ObjectRepository) GetProps(objectID int) (map[string]*model.StoreProp, error) {
	rows := make([]*model.StoreProp, 0)

	err := o.store.db.Model(&model.StoreProp{}).
		Where("object_id = ?", objectID).
		Scan(&rows).Error

	if err != nil {
		return nil, errors.Wrap(err, "GetProps")
	}

	m := make(map[string]*model.StoreProp, len(rows))
	for _, row := range rows {
		m[row.Code] = row
	}

	return m, nil
}

// GetPropsByObjectIDs получение объектов по идентификаторам
func (o *ObjectRepository) GetPropsByObjectIDs(objectIDs []int) (map[int]map[string]*model.StoreProp, error) {
	rows := make([]*model.StoreProp, 0)

	if err := o.store.db.Where("object_id in ?", objectIDs).Find(&rows).Error; err != nil {
		return nil, errors.Wrap(err, "GetObjectsByIDs")
	}

	r := make(map[int]map[string]*model.StoreProp, len(rows))
	for _, row := range rows {
		m, ok := r[row.ObjectID]
		if !ok {
			m = make(map[string]*model.StoreProp, 10)
			r[row.ObjectID] = m
		}
		m[row.Code] = row
	}

	return r, nil
}

// GetPropByObjectIDs получение объектов по идентификаторам
func (o *ObjectRepository) GetPropByObjectIDs(objectIDs []int, code string) (map[int]*model.StoreProp, error) {
	rows := make([]*model.StoreProp, 0)

	err := o.store.db.
		Where("object_id in ?", objectIDs).
		Where("code = ?", code).
		Find(&rows).Error
	if err != nil {
		return nil, errors.Wrap(err, "GetObjectsByIDs")
	}

	r := make(map[int]*model.StoreProp, len(rows))
	for _, row := range rows {
		r[row.ObjectID] = row
	}

	return r, nil
}

// SetProps устанавливает значения свойств объекта
func (o *ObjectRepository) SetProps(objectID int, props map[string]string) error {
	for code, value := range props {
		if err := o.SetProp(objectID, code, value); err != nil {
			return errors.Wrap(err, "SetProps")
		}
	}

	return nil
}

// DelProps удаляет все свойства объекта
func (o *ObjectRepository) DelProps(objectID int) error {
	err := o.store.db.
		Where("object_id = ?", objectID).
		Delete(&model.StoreProp{}).Error

	if err != nil {
		return errors.Wrap(err, "DelProps")
	}

	return nil
}

// GetObject Возвращает объект
func (o *ObjectRepository) GetObject(objectID int) (*model.StoreObject, error) {
	obj := &model.StoreObject{}

	err := o.store.db.
		Where("id = ?", objectID).
		Find(obj).Error

	if err != nil {
		return nil, errors.Wrap(err, "GetObject")
	}

	if obj.ID == 0 {
		return nil, errors.Wrap(errNotFound, "GetObject")
	}

	return obj, nil
}

func (o *ObjectRepository) GetObjectByParent(parentID int, typeObject string) (*model.StoreObject, error) {
	obj := &model.StoreObject{}

	err := o.store.db.
		Where("parent_id = ?", parentID).
		Where("type = ?", typeObject).
		Find(obj).Error

	if err != nil {
		return nil, errors.Wrap(err, "GetObjectByParent")
	}

	if obj.ID == 0 {
		return nil, errors.Wrap(errNotFound, "GetObjectByParent")
	}

	return obj, nil
}

// SetObjectStatus Задает статус объекта
func (o *ObjectRepository) SetObjectStatus(objectID int, status string) error {
	err := o.store.db.
		Model(&model.StoreObject{}).
		Where("id = ?", objectID).
		Update("status", status).Error

	if err != nil {
		return errors.Wrap(err, "SetObject")
	}

	return nil
}

// SaveObject Создает, либо обновляет объект
func (o *ObjectRepository) SaveObject(object *model.StoreObject) error {
	if object == nil {
		return errors.Wrap(errors.New("object is nil"), "SaveObject")
	}

	if object.ParentID != nil && *object.ParentID <= 0 {
		object.ParentID = nil
	}

	if err := o.store.db.Save(object).Error; err != nil {
		return errors.Wrap(err, "SaveObject")
	}

	return nil
}

func (o *ObjectRepository) prepareQuery(filters map[string]interface{}, tags []string, objectType model.ChildType) (*gorm.DB, error) {
	q := o.store.db.Model(&model.StoreObject{})

	for k, v := range filters {
		switch v := v.(type) {
		case string:
			if k == "name" {
				q = q.Where("name like ?", "%"+v+"%")
				continue
			}
			q = q.Where(k+" = ?", v)
		case int:
			q = q.Where(k+" = ?", v)
		case nil:
			q = q.Where(k + " is null")
		default:
			return nil, errors.Wrap(errors.Errorf("unexpected filter %q type: %T", k, v), "prepareQuery")
		}
	}

	if tags != nil {
		for _, tag := range tags {
			q = q.Where(fmt.Sprintf("json_extract(tags, '$.%s')", tag))
		}
	}

	switch objectType {
	case model.ChildTypeInternal:
		q = q.Where("internal = true")
	case model.ChildTypeExternal:
		q = q.Where("internal = false")
	}

	return q, nil
}

// GetObjects получение объектов с учетом фильтров
func (o *ObjectRepository) GetObjects(filters map[string]interface{}, tags []string, offset, limit int, objectType model.ChildType) ([]*model.StoreObject, error) {
	q, err := o.prepareQuery(filters, tags, objectType)
	if err != nil {
		return nil, errors.Wrap(err, "GetObjects")
	}

	q = q.Offset(offset).Limit(limit)

	rows := make([]*model.StoreObject, 0)
	if err := q.Find(&rows).Error; err != nil {
		return nil, errors.Wrap(err, "GetObjects")
	}

	return rows, nil
}

// GetTotal получение общего кол-ва объектов с учетом фильтров
func (o *ObjectRepository) GetTotal(filters map[string]interface{}, tags []string, objectType model.ChildType) (int, error) {
	q, err := o.prepareQuery(filters, tags, objectType)
	if err != nil {
		return 0, errors.Wrap(err, "GetTotal")
	}

	r := int64(0)
	if err := q.Count(&r).Error; err != nil {
		return 0, errors.Wrap(err, "GetTotal")
	}

	return int(r), nil
}

// GetObjectsByTags получение объектов по тегам
func (o *ObjectRepository) GetObjectsByTags(tags []string, offset, limit int, objectType model.ChildType) ([]*model.StoreObject, error) {
	q := o.store.db.Model(&model.StoreObject{}).Offset(offset).Limit(limit)

	for _, tag := range tags {
		q = q.Where(fmt.Sprintf("json_extract(tags, '$.%s')", tag))
	}

	switch objectType {
	case model.ChildTypeInternal:
		q = q.Where("internal = true")
	case model.ChildTypeExternal:
		q = q.Where("internal = false")
	}

	rows := make([]*model.StoreObject, 0)
	if err := q.Find(&rows).Error; err != nil {
		return nil, errors.Wrap(err, "GetObjectsByTags")
	}

	return rows, nil
}

// GetTotalByTags получение общего кол-ва объектов по тегам
func (o *ObjectRepository) GetTotalByTags(tags []string, objectType model.ChildType) (int, error) {
	q := o.store.db.Model(&model.StoreObject{})

	for _, tag := range tags {
		q = q.Where(fmt.Sprintf("json_extract(tags, '$.%s')", tag))
	}

	switch objectType {
	case model.ChildTypeInternal:
		q = q.Where("internal = true")
	case model.ChildTypeExternal:
		q = q.Where("internal = false")
	}

	r := int64(0)
	if err := q.Count(&r).Error; err != nil {
		return 0, errors.Wrap(err, "GetTotalByTags")
	}

	return int(r), nil
}

// GetObjectsByIDs получение объектов по идентификаторам
func (o *ObjectRepository) GetObjectsByIDs(ids []int) ([]*model.StoreObject, error) {
	rows := make([]*model.StoreObject, 0)

	if err := o.store.db.Where("id in ?", ids).Find(&rows).Error; err != nil {
		return nil, errors.Wrap(err, "GetObjectsByIDs")
	}

	return rows, nil
}

// GetObjectChildren получение дочерних объектов
func (o *ObjectRepository) GetObjectChildren(childType model.ChildType, objectIDs ...int) ([]*model.StoreObject, error) {
	rows := make([]*model.StoreObject, 0)

	q := o.store.db.Where("parent_id in ?", objectIDs)

	switch childType {
	case model.ChildTypeInternal:
		q = q.Where("internal = true")
	case model.ChildTypeExternal:
		q = q.Where("internal = false")
	}

	if err := q.Find(&rows).Error; err != nil {
		return nil, errors.Wrap(err, "GetObjectChildren")
	}

	return rows, nil
}

// DelObject удаляет объект
func (o *ObjectRepository) DelObject(objectID int) error {
	// foreign_keys - для каскадного удаления записей
	err := o.store.db.Exec("DELETE FROM om_objects WHERE id = ?", objectID).Error
	if err != nil {
		return errors.Wrap(err, "DelObject")
	}

	return nil
}

func (o *ObjectRepository) GetAllTags() (map[string]int, error) {
	rows := make([]string, 0, 2000)

	if err := o.store.db.Model(&model.StoreObject{}).Select("tags").Find(&rows).Error; err != nil {
		return nil, errors.Wrap(err, "GetAllTags")
	}

	r := make(map[string]int, 500)
	for _, row := range rows {
		if json.Valid([]byte(row)) {
			tags := make(map[string]bool, 10)
			if err := json.Unmarshal([]byte(row), &tags); err != nil {
				return nil, errors.Wrap(err, "GetAllTags")
			}

			for tag := range tags {
				r[tag] += 1
			}
		}
	}

	return r, nil
}

func (o *ObjectRepository) GetObjectsByAddress(address []string) ([]*model.StoreObject, error) {
	q := o.store.db.Table("om_objects")
	q.Joins("INNER JOIN om_props on om_props.object_id = om_objects.id")
	q.Where("code = ?", "address")

	q.Where("value LIKE ?", "%"+address[0]+"%")
	if len(address) > 1 {
		q.Or("value LIKE ?", "%"+address[1]+"%")
	}

	rows := make([]*model.StoreObject, 0)
	if err := q.Find(&rows).Error; err != nil {
		return nil, errors.Wrap(err, "GetObjects")
	}

	return rows, nil
}

func (o *ObjectRepository) SetParent(objectID int, parentID *int) error {
	q := o.store.db.Model(&model.StoreObject{})
	q.Where("id = ?", objectID)
	q.Update("parent_id", parentID)
	err := q.Error

	if err != nil {
		return errors.Wrap(err, "SetParent")
	}

	return nil
}
