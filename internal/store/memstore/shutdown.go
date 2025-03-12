package memstore

import (
	"github.com/pkg/errors"
	"touchon-server/internal/g"
)

func (o *MemStore) Shutdown() error {
	o.mu.Lock()
	defer o.mu.Unlock()

	// Строим дерево, чтобы сначала останавливались дочерние объекты, и только после них родительские
	list, err := o.makeTree()
	if err != nil {
		return errors.Wrap(err, "Shutdown")
	}

	errs := shutdownTree(list)

	// Выводим в лог все ошибки
	for _, err := range errs {
		g.Logger.Error(errors.Wrap(err, "Shutdown"))
	}

	// Возвращаем первую
	if len(errs) > 0 {
		return errors.Wrap(errs[0], "Shutdown")
	}

	return nil
}

func shutdownTree(items []*treeItem) (errs []error) {
	for _, item := range items {
		// Останавливаем только запущенные объекты
		if !item.Object.GetStarted() {
			continue
		}

		// Сначала останавливаем дочерние объекты
		if len(item.Children) > 0 {
			if childrenErrs := shutdownTree(item.Children); len(childrenErrs) > 0 {
				errs = append(errs, childrenErrs...)
			}
		}

		// Затем останавливаем сам объект
		if err := item.Object.Shutdown(); err != nil {
			errs = append(errs, errors.Wrap(err, "shutdownTree"))
		}
	}

	return errs
}
