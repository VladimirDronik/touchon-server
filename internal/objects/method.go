package objects

import (
	"github.com/pkg/errors"
	"touchon-server/internal/scripts"
	"touchon-server/lib/interfaces"
)

type MethodFunc func(params map[string]interface{}) ([]interfaces.Message, error)

func NewMethod(name, description string, params []*scripts.Param, f MethodFunc) (*Method, error) {
	s := &Method{
		Name:        name,
		Description: description,
		Params:      scripts.NewParams(),
		Func:        f,
	}

	switch {
	case name == "":
		return nil, errors.New("name is empty")
	case f == nil:
		return nil, errors.New("f is nil")
	}

	for _, p := range params {
		if err := s.Params.Add(p); err != nil {
			return nil, errors.Wrap(err, "NewMethod")
		}
	}

	return s, nil
}

type Method struct {
	Name        string          `json:"name"`
	Description string          `json:"description,omitempty"`
	Params      *scripts.Params `json:"params"`
	Func        MethodFunc      `json:"-"`
}
