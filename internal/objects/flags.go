package objects

import "github.com/pkg/errors"

type Flags uint64

const (
	CreationForbidden Flags = 1 << iota // Запрещено создание через API
	DeletionForbidden                   // Запрещено удаление через API
	Hidden                              // Скрывать объект в дереве устройств в панели администратора и не только там
)

const (
	Internal = Hidden | CreationForbidden | DeletionForbidden
)

func (o Flags) Err() error {
	switch o {
	case CreationForbidden:
		return errors.New("creation forbidden")
	case DeletionForbidden:
		return errors.New("deletion forbidden")
	case Hidden:
		return errors.New("hidden object")
	}

	return errors.Errorf("unknown flag %d", o)
}

func (o Flags) Has(flag Flags) bool {
	return o&Flags(flag) != 0
}
