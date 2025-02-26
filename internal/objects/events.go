package objects

import (
	"encoding/json"
	"sync"

	"github.com/pkg/errors"
	"touchon-server/lib/helpers/orderedmap"
	"touchon-server/lib/interfaces"
)

func NewEvents() *Events {
	return &Events{
		mu: sync.RWMutex{},
		m:  orderedmap.New[string, interfaces.Event](10),
	}
}

type Events struct {
	mu sync.RWMutex
	m  *orderedmap.OrderedMap[string, interfaces.Event]
}

func (o *Events) Len() int {
	o.mu.RLock()
	defer o.mu.RUnlock()

	return o.m.Len()
}

func (o *Events) Delete(eventNames ...string) {
	o.mu.Lock()
	defer o.mu.Unlock()

	for _, eventName := range eventNames {
		o.m.Delete(eventName)
	}
}

func (o *Events) DeleteAll() {
	o.mu.Lock()
	defer o.mu.Unlock()

	o.m.Clear()
}

func (o *Events) Add(items ...interfaces.Event) error {
	o.mu.Lock()
	defer o.mu.Unlock()

	items = ToApiEvents(items...)

	for _, item := range items {
		if err := o.m.Add(item.GetName(), item); err != nil {
			return errors.Wrap(err, "Events.Add")
		}
	}

	return nil
}

func (o *Events) MarshalJSON() ([]byte, error) {
	o.mu.RLock()
	defer o.mu.RUnlock()

	return json.Marshal(o.m)
}

func (o *Events) UnmarshalJSON([]byte) error {
	// События нельзя переопределять с фронта
	return nil
}

type ApiEventImpl struct {
	interfaces.Event
}

func (o *ApiEventImpl) MarshalJSON() ([]byte, error) {
	type R struct {
		Code        string `json:"code"`
		Name        string `json:"name"`
		Description string `json:"description,omitempty"`
	}

	return json.Marshal(&R{o.GetName(), o.GetTitle(), o.GetDescription()})
}

// ToApiEvents заменяет реализацию MarshalJSON() для того, чтобы скрыть лишние поля
func ToApiEvents(items ...interfaces.Event) []interfaces.Event {
	r := make([]interfaces.Event, 0, len(items))

	for _, item := range items {
		r = append(r, &ApiEventImpl{Event: item})
	}

	return r
}
