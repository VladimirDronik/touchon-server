package objects

import (
	"encoding/json"
	"sort"

	"github.com/pkg/errors"
)

func NewMethods() *Methods {
	return &Methods{
		m: make(map[string]*Method, 20),
	}
}

type Methods struct {
	m map[string]*Method
}

func (o *Methods) Len() int {
	return len(o.m)
}

func (o *Methods) Get(name string) (*Method, error) {
	m, ok := o.m[name]
	if !ok {
		return nil, errors.Wrap(errors.Errorf("method %q not found", name), "Get")
	}

	return m, nil
}

func (o *Methods) GetAll() map[string]*Method {
	return o.m
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
