package objects

import (
	"encoding/json"
	"sort"

	"github.com/pkg/errors"
	"touchon-server/lib/interfaces"
)

func NewMethods() *Methods {
	return &Methods{
		m:       make(map[string]*Method, 20),
		enabled: true,
	}
}

type Methods struct {
	m       map[string]*Method
	enabled bool
}

func (o *Methods) Len() int {
	return len(o.m)
}

func (o *Methods) Get(name string) (*Method, error) {
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
	for _, item := range items {
		o.m[item.Name] = item
	}
}

func (o *Methods) MarshalJSON() ([]byte, error) {
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
	o.enabled = v
}
