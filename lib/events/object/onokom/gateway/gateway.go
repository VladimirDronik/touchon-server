package gateway

import (
	"github.com/pkg/errors"
	"touchon-server/lib/interfaces"
	"touchon-server/lib/messages"
)

func NewOnCheck(targetID int, props map[string]interface{}) (interfaces.Message, error) {
	msg, err := messages.NewEvent(interfaces.TargetTypeObject, targetID)
	if err != nil {
		return nil, errors.Wrap(err, "NewOnCheck")
	}

	o := &OnCheck{
		MessageImpl: msg,
		Props:       props,
	}

	return o, nil
}

// OnCheck Получение состояния устройства
type OnCheck struct {
	*messages.MessageImpl
	Props map[string]interface{} `json:"props"` // Свойства кондиционера
}

func NewOnChange(targetID int, props map[string]interface{}) (interfaces.Message, error) {
	msg, err := messages.NewEvent(interfaces.TargetTypeObject, targetID)
	if err != nil {
		return nil, errors.Wrap(err, "NewOnChange")
	}

	o := &OnChange{
		MessageImpl: msg,
		Props:       props,
	}

	return o, nil
}

// OnChange Изменение состояния устройства
type OnChange struct {
	*messages.MessageImpl
	Props map[string]interface{} `json:"props"` // Свойства кондиционера
}
