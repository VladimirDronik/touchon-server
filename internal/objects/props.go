package objects

import (
	"encoding/json"

	"github.com/pkg/errors"
	"touchon-server/lib/helpers/orderedmap"
	"touchon-server/lib/models"
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

func (o *Props) GetAll() *orderedmap.OrderedMap[string, *Prop] {
	return o.m
}

func (o *Props) Get(code string) (*Prop, error) {
	v, err := o.m.Get(code)
	if err != nil {
		return nil, errors.Wrap(err, "Get")
	}

	return v, nil
}

func (o *Props) Delete(code string) {
	o.m.Delete(code)
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
		return "", errors.Wrapf(err, "GetStringValue(%s)", code)
	}

	v, err := p.GetStringValue()
	if err != nil {
		return "", errors.Wrapf(err, "GetStringValue(%s)", code)
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

// Check Проверяет значения параметров в отдельности
func (o *Props) Check() error {
	for _, p := range o.m.GetValueList() {
		switch p.Type {
		case models.DataTypeString:
			_, err := p.GetStringValue()
			if p.Required.Check(o) && err != nil {
				return errors.Wrapf(errors.New("value is not set"), "check(%s)", p.Code)
			}
			if err == nil && p.CheckValue != nil {
				if err := p.CheckValue(p, o.m.GetUnorderedMap()); err != nil {
					return errors.Wrapf(err, "check(%s)", p.Code)
				}
			}

		case models.DataTypeEnum:
			v, err := p.GetEnumValue()
			if p.Required.Check(o) && err != nil {
				return errors.Wrapf(errors.New("value is not set"), "check(%s)", p.Code)
			}
			if err == nil {
				if _, ok := p.Values[v]; !ok {
					return errors.Wrapf(errors.New("value is bad"), "check(%s)", p.Code)
				}

				if p.CheckValue != nil {
					if err := p.CheckValue(p, o.m.GetUnorderedMap()); err != nil {
						return errors.Wrapf(err, "check(%s)", p.Code)
					}
				}
			}

		case models.DataTypeBool:
			_, err := p.GetBoolValue()
			if p.Required.Check(o) && err != nil {
				return errors.Wrapf(errors.New("value is not set"), "check(%s)", p.Code)
			}
			if err == nil && p.CheckValue != nil {
				if err := p.CheckValue(p, o.m.GetUnorderedMap()); err != nil {
					return errors.Wrapf(err, "check(%s)", p.Code)
				}
			}

		case models.DataTypeInt:
			_, err := p.GetIntValue()
			if p.Required.Check(o) && err != nil {
				return errors.Wrapf(errors.New("value is not set"), "check(%s)", p.Code)
			}
			if err == nil && p.CheckValue != nil {
				if err := p.CheckValue(p, o.m.GetUnorderedMap()); err != nil {
					return errors.Wrapf(err, "check(%s)", p.Code)
				}
			}

		case models.DataTypeFloat:
			_, err := p.GetFloatValue()
			if p.Required.Check(o) && err != nil {
				return errors.Wrapf(errors.New("value is not set"), "check(%s)", p.Code)
			}
			if err == nil && p.CheckValue != nil {
				if err := p.CheckValue(p, o.m.GetUnorderedMap()); err != nil {
					return errors.Wrapf(err, "check(%s)", p.Code)
				}
			}

		case models.DataTypeInterface:
			// do nothing

		default:
			return errors.Wrapf(errors.Errorf("unexpected prop data type %q", p.Type), "check(%s)", p.Code)
		}
	}

	return nil
}

func (o *Props) MarshalJSON() ([]byte, error) {
	return json.Marshal(o.m)
}

func (o *Props) UnmarshalJSON(data []byte) error {
	src := NewProps()
	dst := o

	if err := json.Unmarshal(data, &src.m); err != nil {
		return errors.Wrap(err, "Props.UnmarshalJSON")
	}

	if src.m.Len() != dst.m.Len() {
		return errors.Wrap(errors.Errorf("src.m.Len() != dst.m.Len(), %d != %d", src.m.Len(), dst.m.Len()), "Props.UnmarshalJSON")
	}

	for _, srcProp := range src.m.GetValueList() {
		dstProp, err := dst.m.Get(srcProp.Code)
		if err != nil {
			return errors.Wrap(err, "Props.UnmarshalJSON")
		}

		if err := dstProp.SetValue(srcProp.GetValue()); err != nil {
			return errors.Wrap(err, "Props.UnmarshalJSON")
		}
	}

	return nil
}

func (o *Props) CheckDefinition() error {
	for _, p := range o.m.GetValueList() {
		// Проверяем определение свойства
		if err := p.CheckDefinition(); err != nil {
			return errors.Wrap(err, "Props.CheckDefinition")
		}
	}

	return nil
}
