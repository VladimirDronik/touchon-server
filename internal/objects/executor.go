package objects

import (
	"github.com/pkg/errors"
	"touchon-server/internal/model"
	"touchon-server/internal/scripts"
	"touchon-server/internal/store"
)

func NewExecutor() scripts.ObjectMethodExecutor {
	// В замыкании сохраняем store & logger
	// Возвращаем ошибку в виде текста, потому что
	// не IT'шнику проще будет работать с текстовыми ошибками
	return func(id int, objCat, objType, method string, args map[string]interface{}) ([]interface{}, string) {
		// Для удобства внутри возвращаем ошибку, здесь разворачиваем ее в текст
		r, err := func() ([]interface{}, error) {
			// Получаем объекты по ID или по категории и типу
			objModels, err := getObjects(id, model.Category(objCat), objType)
			if err != nil {
				return nil, errors.Wrap(err, "getObjects")
			}

			// Список результатов
			items := make([]interface{}, 0, len(objModels))

			for _, objModel := range objModels {
				method, err := objModel.GetMethods().Get(method)
				if err != nil {
					return nil, errors.Wrap(err, "getObjects")
				}

				r, err := method.Func(args)
				if err != nil {
					return nil, errors.Wrap(err, "getObjects")
				}

				items = append(items, r)
			}

			return items, nil
		}()

		if err != nil {
			return nil, err.Error()
		}

		return r, ""
	}
}

// getObjects возвращает либо один объект по его ID, либо список объектов по категории и типу
func getObjects(objectID int, objCat model.Category, objType string) ([]Object, error) {
	if objectID > 0 {
		objModel, err := LoadObject(objectID, objCat, objType, model.ChildTypeNobody)
		if err != nil {
			return nil, errors.Wrap(err, "getObjects")
		}

		return []Object{objModel}, nil
	}

	filters := map[string]interface{}{
		"category": objCat,
		"type":     objType,
	}

	objects, err := store.I.ObjectRepository().GetObjects(filters, nil, 0, 1000, model.ChildTypeNobody)
	if err != nil {
		return nil, errors.Wrap(err, "getObjects")
	}

	objModels := make([]Object, 0, len(objects))
	for _, obj := range objects {
		objModel, err := LoadObject(obj.ID, "", "", model.ChildTypeNobody)
		if err != nil {
			return nil, errors.Wrap(err, "getObjects")
		}

		objModels = append(objModels, objModel)
	}

	return objModels, nil
}
