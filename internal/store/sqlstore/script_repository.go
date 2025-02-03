package sqlstore

import (
	"github.com/pkg/errors"
	"touchon-server/internal/model"
)

//create table if not exists scripts
//(
//    id          INTEGER primary key autoincrement,
//    code        TEXT not null,
//    name        TEXT not null,
//    description TEXT not null default '',
//    params      TEXT not null default '{}',
//    body        TEXT not null
//);
//
//CREATE UNIQUE INDEX if not exists code ON scripts(code);

type ScriptRepository struct {
	store *Store
}

// GetScript возвращает сценарий
func (o *ScriptRepository) GetScript(id int) (*model.StoreScript, error) {
	obj := &model.StoreScript{}

	err := o.store.db.
		Where("id = ?", id).
		Find(obj).Error

	if err != nil {
		return nil, errors.Wrap(err, "GetScript")
	}

	if obj.ID == 0 {
		return nil, errors.Wrap(errNotFound, "GetScript")
	}

	return obj, nil
}

// GetScriptByCode возвращает сценарий
func (o *ScriptRepository) GetScriptByCode(code string) (*model.StoreScript, error) {
	obj := &model.StoreScript{}

	err := o.store.db.
		Where("code = ?", code).
		Find(obj).Error

	if err != nil {
		return nil, errors.Wrap(err, "GetScriptByCode")
	}

	if obj.ID == 0 {
		return nil, errors.Wrap(errNotFound, "GetScriptByCode")
	}

	return obj, nil
}

// SetScript Создает, либо обновляет сценарий
func (o *ScriptRepository) SetScript(script *model.StoreScript) error {
	if script == nil {
		return errors.Wrap(errors.New("script is nil"), "SetScript")
	}

	count := int64(0)
	if err := o.store.db.Model(script).Where("id = ?", script.ID).Count(&count).Error; err != nil {
		return errors.Wrap(err, "SetScript")
	}
	scriptIsExists := count == 1

	if scriptIsExists {
		if err := o.store.db.Updates(script).Error; err != nil {
			return errors.Wrap(err, "SetScript(update)")
		}
	} else {
		if err := o.store.db.Create(&script).Error; err != nil {
			return errors.Wrap(err, "SetScript(create)")
		}
	}

	return nil
}

// GetScripts возвращает список скриптов
func (o *ScriptRepository) GetScripts(code, name string, offset, limit int) ([]*model.StoreScript, error) {
	rows := make([]*model.StoreScript, 0)

	q := o.store.db.
		Offset(offset).
		Limit(limit)

	if code != "" {
		q = q.Where("code like ?", "%"+code+"%")
	}

	if name != "" {
		q = q.Where("name like ?", "%"+name+"%")
	}

	if err := q.Find(&rows).Error; err != nil {
		return nil, errors.Wrap(err, "GetScripts")
	}

	return rows, nil
}

// GetTotal возвращает общее кол-во найденных скриптов
func (o *ScriptRepository) GetTotal(code, name string, offset, limit int) (int, error) {
	r := int64(0)

	q := o.store.db.
		Model(&model.StoreScript{}).
		Offset(offset).
		Limit(limit)

	if code != "" {
		q = q.Where("code like ?", "%"+code+"%")
	}

	if name != "" {
		q = q.Where("name like ?", "%"+name+"%")
	}

	if err := q.Count(&r).Error; err != nil {
		return 0, errors.Wrap(err, "GetScripts")
	}

	return int(r), nil
}

// DelScript удаляет сценарий
func (o *ScriptRepository) DelScript(id int) error {
	if err := o.store.db.Delete(&model.StoreScript{ID: id}).Error; err != nil {
		return errors.Wrap(err, "DelScript")
	}

	return nil
}
