package objects

import (
	"encoding/json"

	"github.com/pkg/errors"
)

func NewChildren() *Children {
	return &Children{
		items: make([]Object, 0, 10),
	}
}

type Children struct {
	items []Object // Для сохранения порядка детей
}

func (o *Children) Len() int {
	return len(o.items)
}

func (o *Children) GetAll() []Object {
	return o.items
}

func (o *Children) DeleteByID(objectID int) {
	for i := 0; i < len(o.items); i++ {
		if o.items[i].GetID() == objectID {
			// Сдвигаем правую часть списка - влево на одну позицию
			copy(o.items[i:], o.items[i+1:])
			// Уменьшаем массив на один элемент
			o.items = o.items[:len(o.items)-1]
			return
		}
	}
}

func (o *Children) DeleteAll() {
	o.items = o.items[:0]
}

func (o *Children) Add(items ...Object) {
	o.items = append(o.items, items...)
}

func (o *Children) MarshalJSON() ([]byte, error) {
	return json.Marshal(o.items)
}

func (o *Children) UnmarshalJSON(data []byte) error {
	children := make([]json.RawMessage, 0, len(o.items))

	if err := json.Unmarshal(data, &children); err != nil {
		return errors.Wrap(err, "Children.UnmarshalJSON")
	}

	if len(o.items) != len(children) {
		return errors.Wrap(errors.Errorf("len(o.items) != len(Children), %d != %d", len(o.items), len(children)), "Children.UnmarshalJSON")
	}

	for i, childData := range children {
		if err := json.Unmarshal(childData, &o.items[i]); err != nil {
			return errors.Wrap(err, "Children.UnmarshalJSON")
		}
	}

	return nil
}
