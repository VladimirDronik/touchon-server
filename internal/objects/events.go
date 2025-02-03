package objects

import (
	"encoding/json"

	"github.com/pkg/errors"
	"touchon-server/lib/event"
	"touchon-server/lib/helpers/orderedmap"
)

func NewEvents() *Events {
	return &Events{
		m: orderedmap.New[string, *event.Event](10),
	}
}

type Events struct {
	m *orderedmap.OrderedMap[string, *event.Event]
}

func (o *Events) Len() int {
	return o.m.Len()
}

func (o *Events) GetAll() []*event.Event {
	return o.m.GetValueList()
}

func (o *Events) Delete(eventNames ...string) {
	for _, eventName := range eventNames {
		o.m.Delete(eventName)
	}
}

func (o *Events) DeleteAll() {
	o.m.Clear()
}

func (o *Events) Add(items ...*event.Event) error {
	for _, item := range items {
		if err := o.m.Add(item.Code, item); err != nil {
			return errors.Wrap(err, "Events.Add")
		}
	}

	return nil
}

func (o *Events) MarshalJSON() ([]byte, error) {
	return json.Marshal(o.m)
}

func (o *Events) UnmarshalJSON([]byte) error {
	// События нельзя переопределять с фронта
	return nil
}
