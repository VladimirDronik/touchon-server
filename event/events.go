package event

import (
	"encoding/json"

	"github.com/VladimirDronik/touchon-server/helpers/orderedmap"
	"github.com/pkg/errors"
)

func NewEvents() *Events {
	return &Events{
		m: orderedmap.New[string, *Event](10),
	}
}

type Events struct {
	m *orderedmap.OrderedMap[string, *Event]
}

func (o *Events) Len() int {
	return o.m.Len()
}

func (o *Events) GetAll() []*Event {
	return o.m.GetValueList()
}

func (o *Events) DeleteAll() {
	o.m.Clear()
}

func (o *Events) Add(items ...*Event) error {
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

func (o *Events) UnmarshalJSON(data []byte) error {
	if err := json.Unmarshal(data, o.m); err != nil {
		return errors.Wrap(err, "Events.UnmarshalJSON")
	}

	return nil
}
