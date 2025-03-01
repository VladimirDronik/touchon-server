package objects

import (
	"fmt"

	"github.com/pkg/errors"
	"touchon-server/internal/model"
)

// ObjectMaker Функция для создания экземпляра модели объекта
type ObjectMaker func(withChildren bool) (Object, error)

// Реестр объектов
var register = make(map[string]ObjectMaker, 50)

type objectAttr struct {
	Name string
	Tags []string
}

var categoriesAndTypes = make(map[string]map[string]objectAttr, 50)

func getKey(objCat model.Category, objType string) string {
	return fmt.Sprintf("%s:%s", objCat, objType)
}

// Register проверяет модель объекта.
// Необходимо для проверки определения моделей
func Register(maker ObjectMaker) (e error) {
	defer func() {
		if e != nil {
			panic(errors.Wrap(e, "object.Register"))
		}
	}()

	// Создаем экземпляр модели объекта для проверки модели
	obj, err := maker(true)
	if err != nil {
		return err
	}

	objCat := obj.GetCategory()
	objType := obj.GetType()
	objTags := obj.GetTags()

	switch {
	case objCat == "":
		return errors.New("object category is empty")
	case objType == "":
		return errors.New("object type is empty")
	}

	key := getKey(objCat, objType)
	if _, ok := register[key]; ok {
		return errors.Errorf("object %q registered already", key)
	}

	if props := obj.GetProps(); props != nil {
		if err := props.CheckDefinition(); err != nil {
			return errors.Wrapf(err, "Register(%s)", key)
		}
	}

	register[key] = maker

	types, ok := categoriesAndTypes[string(objCat)]
	if !ok {
		types = make(map[string]objectAttr, 10)
		categoriesAndTypes[string(objCat)] = types
	}

	if _, ok := types[objType]; !ok {
		types[objType] = objectAttr{obj.GetName(), objTags}
	}

	return nil
}

func GetObjectModel(objCat model.Category, objType string, withChildren bool) (Object, error) {
	key := getKey(objCat, objType)
	maker, ok := register[key]
	if !ok {
		return nil, errors.Wrapf(errors.New("object not found"), "GetObjectModel(%s, %s)", string(objCat), objType)
	}

	obj, err := maker(withChildren)
	if err != nil {
		return nil, errors.Wrap(err, "GetObjectModel")
	}

	return obj, nil
}

func GetCategoriesAndTypes() map[string]map[string]objectAttr {
	return categoriesAndTypes
}
