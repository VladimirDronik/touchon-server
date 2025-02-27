package orderedmap

import (
	"bytes"
	"encoding/json"
	"io"
	"slices"
	"sync"

	"github.com/pkg/errors"
)

type KeyType interface {
	comparable
}

func New[K KeyType, V any](cap int) *OrderedMap[K, V] {
	return &OrderedMap[K, V]{
		mu: sync.RWMutex{},
		m:  make(map[K]V, cap),
	}
}

type Item[K KeyType, V any] struct {
	Key   K
	Value V
}

type OrderedMap[K KeyType, V any] struct {
	mu    sync.RWMutex
	m     map[K]V
	order []K
}

func (o *OrderedMap[K, V]) Add(k K, v V) error {
	o.mu.Lock()
	defer o.mu.Unlock()

	if _, ok := o.m[k]; ok {
		return errors.Wrap(errors.Errorf("value with key %v is exists", k), "Add")
	}

	o.m[k] = v
	o.order = append(o.order, k)

	return nil
}

func (o *OrderedMap[K, V]) Set(k K, v V) {
	o.mu.RLock()
	_, ok := o.m[k]
	o.mu.RUnlock()

	if !ok {
		_ = o.Add(k, v)
		return
	}

	o.mu.Lock()
	o.m[k] = v
	o.mu.Unlock()
}

func (o *OrderedMap[K, V]) Get(k K) (V, error) {
	o.mu.RLock()
	defer o.mu.RUnlock()

	var null V

	v, ok := o.m[k]
	if !ok {
		return null, errors.Wrap(errors.Errorf("value with key %v not found", k), "Get")
	}

	return v, nil
}

func (o *OrderedMap[K, V]) GetKeyValueList() []*Item[K, V] {
	o.mu.RLock()
	defer o.mu.RUnlock()

	r := make([]*Item[K, V], 0, len(o.m))

	for _, k := range o.order {
		r = append(r, &Item[K, V]{k, o.m[k]})
	}

	return r
}

func (o *OrderedMap[K, V]) GetValueList() []V {
	o.mu.RLock()
	defer o.mu.RUnlock()

	r := make([]V, 0, len(o.m))

	for _, k := range o.order {
		r = append(r, o.m[k])
	}

	return r
}

func (o *OrderedMap[K, V]) GetUnorderedMap() map[K]V {
	o.mu.RLock()
	defer o.mu.RUnlock()

	r := make(map[K]V, len(o.m))

	// copy map
	for k, v := range o.m {
		r[k] = v
	}

	return r
}

func (o *OrderedMap[K, V]) Len() int {
	o.mu.RLock()
	defer o.mu.RUnlock()

	return len(o.m)
}

func (o *OrderedMap[K, V]) Delete(k K) {
	o.mu.Lock()
	defer o.mu.Unlock()

	delete(o.m, k)

	for i, item := range o.order {
		if item == k {
			o.order = slices.Delete(o.order, i, i+1)
		}
	}
}

func (o *OrderedMap[K, V]) Clear() {
	o.mu.Lock()
	defer o.mu.Unlock()

	for k := range o.m {
		delete(o.m, k)
	}
	o.order = o.order[:0]
}

func (o *OrderedMap[K, V]) MarshalJSON() ([]byte, error) {
	o.mu.RLock()
	defer o.mu.RUnlock()

	keys := make([]json.RawMessage, 0, len(o.m))
	values := make([]json.RawMessage, 0, len(o.m))

	for _, k := range o.order {
		key, err := json.Marshal(k)
		if err != nil {
			return nil, err
		}

		keys = append(keys, key)

		value, err := json.Marshal(o.m[k])
		if err != nil {
			return nil, err
		}

		values = append(values, value)
	}

	buf := bytes.NewBuffer(nil)
	buf.WriteString("{")

	for i, key := range keys {
		if len(key) == 0 {
			return nil, errors.Errorf("key is empty")
		}

		if key[0] != '"' {
			buf.WriteByte('"')
		}
		buf.Write(keys[i])
		if key[0] != '"' {
			buf.WriteByte('"')
		}
		buf.WriteString(": ")
		buf.Write(values[i])
		if i+1 < len(keys) {
			buf.WriteByte(',')
		}
	}

	buf.WriteByte('}')

	return buf.Bytes(), nil
}

func (o *OrderedMap[K, V]) UnmarshalJSON(data []byte) error {
	dec := json.NewDecoder(bytes.NewReader(data))

	// read {
	t, err := dec.Token()
	if err != nil {
		return errors.Wrap(err, "OrderedMap.UnmarshalJSON")
	}

	for dec.More() {
		// read key
		t, err = dec.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return errors.Wrap(err, "OrderedMap.UnmarshalJSON")
		}

		k, ok := t.(K)
		if !ok {
			return errors.Wrap(errors.Errorf("%v is not %T", t, k), "OrderedMap.UnmarshalJSON")
		}

		var v V
		if err := dec.Decode(&v); err != nil {
			return errors.Wrap(err, "OrderedMap.UnmarshalJSON")
		}

		o.Set(k, v)
	}

	// read }
	_, err = dec.Token()
	if err != nil {
		return errors.Wrap(err, "OrderedMap.UnmarshalJSON")
	}

	return nil
}
