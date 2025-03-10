package scripts

import (
	"encoding/json"

	"github.com/pkg/errors"
	"touchon-server/lib/ordered_map"
)

func NewParams() *Params {
	return &Params{
		m: ordered_map.New[string, *Param](10),
	}
}

type Params struct {
	m *ordered_map.OrderedMap[string, *Param]
}

func (o *Params) Len() int {
	return o.m.Len()
}

func (o *Params) GetAll() *ordered_map.OrderedMap[string, *Param] {
	return o.m
}

func (o *Params) Get(code string) (*Param, error) {
	v, err := o.m.Get(code)
	if err != nil {
		return nil, errors.Wrap(err, "Get")
	}

	return v, nil
}

func (o *Params) Set(code string, value interface{}) error {
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
func (o *Params) GetStringValue(code string) (string, error) {
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
func (o *Params) GetBoolValue(code string) (bool, error) {
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
func (o *Params) GetEnumValue(code string) (string, error) {
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

func (o *Params) GetIntValue(code string) (int, error) {
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

func (o *Params) GetFloatValue(code string) (float32, error) {
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

var chars = map[rune]bool{}

func init() {
	for _, c := range "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_" {
		chars[c] = true
	}
}

func sanitizeCode(v string) string {
	r := make([]rune, 0, len(v))

	for _, c := range v {
		if chars[c] {
			r = append(r, c)
		}
	}

	return string(r)
}

func (o *Params) Add(items ...*Param) error {
	for _, item := range items {
		if item == nil {
			return errors.Wrap(errors.New("param is nil"), "Add")
		}

		item.Code = sanitizeCode(item.Code)

		if err := o.m.Add(item.Code, item); err != nil {
			return errors.Wrap(err, "Add")
		}
	}

	return nil
}

func (o *Params) MarshalJSON() ([]byte, error) {
	return json.Marshal(o.m)
}

func (o *Params) UnmarshalJSON(data []byte) error {
	if err := json.Unmarshal(data, &o.m); err != nil {
		return errors.Wrap(err, "Params.UnmarshalJSON")
	}

	return nil
}

func (o *Params) Check() error {
	for _, p := range o.m.GetValueList() {
		// Проверяем определение свойства
		if err := p.Check(); err != nil {
			return errors.Wrap(err, "check")
		}
	}

	return nil
}
