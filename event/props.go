package event

import (
	"encoding/json"

	"github.com/VladimirDronik/touchon-server/helpers/orderedmap"
	"github.com/pkg/errors"
)

func NewProps() *Props {
	return &Props{
		m: orderedmap.New[string, *Prop](10),
	}
}

type Props struct {
	m *orderedmap.OrderedMap[string, *Prop]
}

func (o *Props) Len() int {
	return o.m.Len()
}

func (o *Props) GetOrderedMap() *orderedmap.OrderedMap[string, *Prop] {
	return o.m
}

func (o *Props) Get(code string) (*Prop, error) {
	v, err := o.m.Get(code)
	if err != nil {
		return nil, errors.Wrap(err, "Get")
	}

	return v, nil
}

func (o *Props) Set(code string, value interface{}) error {
	p, err := o.Get(code)
	if err != nil {
		return errors.Wrap(err, "Set")
	}

	if err := p.SetValue(value); err != nil {
		return errors.Wrapf(err, "Set(%s)", code)
	}

	return nil
}

// GetStringValue Метод-хэлпер для получения строкового значения
func (o *Props) GetStringValue(code string) (string, error) {
	p, err := o.Get(code)
	if err != nil {
		return "", errors.Wrap(err, "GetStringValue")
	}

	v, err := p.GetStringValue()
	if err != nil {
		return "", errors.Wrap(err, "GetStringValue")
	}

	return v, nil
}

// GetBoolValue Метод-хэлпер для получения логического значения
func (o *Props) GetBoolValue(code string) (bool, error) {
	p, err := o.Get(code)
	if err != nil {
		return false, errors.Wrap(err, "GetBoolValue")
	}

	v, err := p.GetBoolValue()
	if err != nil {
		return false, errors.Wrap(err, "GetBoolValue")
	}

	return v, nil
}

// GetEnumValue Метод-хэлпер для получения значения-перечисления
func (o *Props) GetEnumValue(code string) (string, error) {
	p, err := o.Get(code)
	if err != nil {
		return "", errors.Wrap(err, "GetEnumValue")
	}

	v, err := p.GetEnumValue()
	if err != nil {
		return "", errors.Wrap(err, "GetEnumValue")
	}

	return v, nil
}

func (o *Props) GetIntValue(code string) (int, error) {
	p, err := o.Get(code)
	if err != nil {
		return 0, errors.Wrap(err, "GetIntValue")
	}

	v, err := p.GetIntValue()
	if err != nil {
		return 0, errors.Wrap(err, "GetIntValue")
	}

	return v, nil
}

func (o *Props) GetFloatValue(code string) (float32, error) {
	p, err := o.Get(code)
	if err != nil {
		return 0, errors.Wrap(err, "GetFloatValue")
	}

	v, err := p.GetFloatValue()
	if err != nil {
		return 0, errors.Wrap(err, "GetFloatValue")
	}

	return v, nil
}

func (o *Props) Add(items ...*Prop) error {
	for _, item := range items {
		if item == nil {
			return errors.Wrap(errors.New("prop is nil"), "Add")
		}

		if err := o.m.Add(item.Code, item); err != nil {
			return errors.Wrap(err, "Add")
		}
	}

	return nil
}

func (o *Props) MarshalJSON() ([]byte, error) {
	return json.Marshal(o.m)
}

func (o *Props) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, o.m)
}

func (o *Props) Check() error {
	for _, p := range o.m.GetValueList() {
		// Проверяем определение свойства
		if err := p.Check(); err != nil {
			return errors.Wrap(err, "Check")
		}
	}

	return nil
}
