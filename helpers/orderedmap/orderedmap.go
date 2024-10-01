package orderedmap

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/pkg/errors"
)

type KeyType interface {
	comparable
}

func New[K KeyType, V any](cap int) *OrderedMap[K, V] {
	return &OrderedMap[K, V]{
		m: make(map[K]V, cap),
	}
}

type Item[K KeyType, V any] struct {
	Key   K
	Value V
}

type OrderedMap[K KeyType, V any] struct {
	m     map[K]V
	order []K
}

func (o *OrderedMap[K, V]) Add(k K, v V) error {
	if _, ok := o.m[k]; ok {
		return errors.Wrap(errors.Errorf("value with key %v is exists", k), "Add")
	}

	o.m[k] = v
	o.order = append(o.order, k)

	return nil
}

func (o *OrderedMap[K, V]) Set(k K, v V) {
	if _, ok := o.m[k]; !ok {
		_ = o.Add(k, v)
		return
	}

	o.m[k] = v
}

func (o *OrderedMap[K, V]) Get(k K) (V, error) {
	var null V

	v, ok := o.m[k]
	if !ok {
		return null, errors.Wrap(errors.Errorf("value with key %v not found", k), "Get")
	}

	return v, nil
}

func (o *OrderedMap[K, V]) GetKeyValueList() []*Item[K, V] {
	r := make([]*Item[K, V], 0, len(o.m))

	for _, k := range o.order {
		r = append(r, &Item[K, V]{k, o.m[k]})
	}

	return r
}

func (o *OrderedMap[K, V]) GetValueList() []V {
	r := make([]V, 0, len(o.m))

	for _, k := range o.order {
		r = append(r, o.m[k])
	}

	return r
}

func (o *OrderedMap[K, V]) GetUnorderedMap() map[K]V {
	r := make(map[K]V, len(o.m))

	// copy map
	for k, v := range o.m {
		r[k] = v
	}

	return r
}

func (o *OrderedMap[K, V]) Len() int {
	return len(o.m)
}

func (o *OrderedMap[K, V]) Clear() {
	for k := range o.m {
		delete(o.m, k)
	}
	o.order = o.order[:0]
}

func (o *OrderedMap[K, V]) MarshalJSON() ([]byte, error) {
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
		return err
	}

	for dec.More() {
		// read key
		t, err = dec.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		k, ok := t.(K)
		if !ok {
			return errors.Errorf("%v is not %T", t, k)
		}

		var v V
		if err := dec.Decode(&v); err != nil {
			return err
		}

		o.Set(k, v)
	}

	// read }
	_, err = dec.Token()
	if err != nil {
		return err
	}

	return nil
}
