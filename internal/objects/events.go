package objects

import (
	"encoding/json"

	"github.com/pkg/errors"
	"touchon-server/lib/helpers/orderedmap"
	"touchon-server/lib/interfaces"
)

func NewEvents() *Events {
	return &Events{
		m: orderedmap.New[string, interfaces.Event](10),
	}
}

type Events struct {
	m *orderedmap.OrderedMap[string, interfaces.Event]
}

func (o *Events) Len() int {
	return o.m.Len()
}

func (o *Events) GetAll() []interfaces.Event {
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

func (o *Events) Add(items ...interfaces.Event) error {
	for _, item := range items {
		if err := o.m.Add(item.GetEventCode(), item); err != nil {
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
