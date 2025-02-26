package objects

import (
	"encoding/json"
	"sort"
	"sync"

	"github.com/pkg/errors"
	"touchon-server/lib/interfaces"
)

func NewMethods(mu *sync.RWMutex) *Methods {
	return &Methods{
		mu:      mu,
		m:       make(map[string]*Method, 20),
		enabled: true,
	}
}

type Methods struct {
	mu      *sync.RWMutex
	m       map[string]*Method
	enabled bool
}

func (o *Methods) Len() int {
	o.mu.RLock()
	defer o.mu.RUnlock()

	return len(o.m)
}

func (o *Methods) Get(name string) (*Method, error) {
	o.mu.RLock()
	defer o.mu.RUnlock()

	m, ok := o.m[name]
	if !ok {
		return nil, errors.Wrap(errors.Errorf("method %q not found", name), "Get")
	}

	// Если объект отключен, заменяет тело метода заглушкой
	if !o.enabled {
		m = &Method{
			Name:        m.Name,
			Description: m.Description,
			Params:      m.Params,
			Func: func(params map[string]interface{}) ([]interfaces.Message, error) {
				return nil, ErrObjectDisabled
			},
		}
	}

	return m, nil
}

func (o *Methods) GetAll() map[string]*Method {
	o.mu.RLock()
	defer o.mu.RUnlock()

	if o.enabled {
		return o.m
	}

	// Если объект отключен, заменяет тело метода заглушкой
	m := make(map[string]*Method, len(o.m))
	for k := range o.m {
		m[k], _ = o.Get(k)
	}

	return m
}

func (o *Methods) Add(items ...*Method) {
	o.mu.Lock()
	defer o.mu.Unlock()

	for _, item := range items {
		o.m[item.Name] = item

		// Перед вызовом каждого метода блокируем доступ к объекту, т.к. методы
		// могут выполняться одновременно и одновременно обращаться к одной и той
		// же области памяти на запись
		item.Func = func(params map[string]interface{}) ([]interfaces.Message, error) {
			o.mu.Lock()
			defer o.mu.Unlock()

			return item.Func(params)
		}
	}
}

func (o *Methods) MarshalJSON() ([]byte, error) {
	o.mu.RLock()
	defer o.mu.RUnlock()

	items := make([]*Method, 0, len(o.m))
	for _, m := range o.m {
		items = append(items, m)
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].Name < items[j].Name
	})

	return json.Marshal(items)
}

func (o *Methods) UnmarshalJSON([]byte) error {
	// Методы нельзя переопределять с фронта
	return nil
}

func (o *Methods) SetEnabled(v bool) {
	o.mu.Lock()
	defer o.mu.Unlock()

	o.enabled = v
}
