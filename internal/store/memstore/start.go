package memstore

import (
	"github.com/pkg/errors"
	"touchon-server/internal/g"
)

func (o *MemStore) Start() error {
	o.mu.Lock()
	defer o.mu.Unlock()

	// Строим дерево, чтобы сначала запускались родительские объекты, и только после них дочерние
	list, err := o.makeTree()
	if err != nil {
		return errors.Wrap(err, "Start")
	}

	startTree(list)

	return nil
}

func startTree(items []*treeItem) {
	for _, item := range items {
		// Если объект выключен или уже запущен, пропускаем его
		if !item.Object.GetEnabled() || item.Object.GetStarted() {
			continue
		}

		if err := item.Object.Start(); err != nil {
			// При старте сервиса неверно сконфигурированный объект не должен
			// останавливать сервис или прекращать запуск других не дочерних объектов.
			// Сервис должен запуститься, чтобы была возможность изменить св-ва объекта.
			g.Logger.Error(errors.Wrap(err, "startTree"))
			continue
		}

		if len(item.Children) > 0 {
			startTree(item.Children)
		}
	}
}
