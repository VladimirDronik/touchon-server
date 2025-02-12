package generic_input

import (
	"github.com/pkg/errors"
	"touchon-server/lib/interfaces"
	"touchon-server/lib/messages"
)

func NewOnClick(targetID int) (interfaces.Message, error) {
	msg, err := messages.NewEvent(interfaces.TargetTypeObject, targetID)
	if err != nil {
		return nil, errors.Wrap(err, "NewOnClick")
	}

	o := &OnClick{
		MessageImpl: msg,
	}

	return o, nil
}

// OnClick Одиночное замыкание
type OnClick struct {
	*messages.MessageImpl
}

func NewOnDoubleClick(targetID int) (interfaces.Message, error) {
	msg, err := messages.NewEvent(interfaces.TargetTypeObject, targetID)
	if err != nil {
		return nil, errors.Wrap(err, "NewOnDoubleClick")
	}

	o := &OnDoubleClick{
		MessageImpl: msg,
	}

	return o, nil
}

// OnDoubleClick Двойное замыкание
type OnDoubleClick struct {
	*messages.MessageImpl
}

func NewOnLongPress(targetID int) (interfaces.Message, error) {
	msg, err := messages.NewEvent(interfaces.TargetTypeObject, targetID)
	if err != nil {
		return nil, errors.Wrap(err, "NewOnLongPress")
	}

	o := &OnLongPress{
		MessageImpl: msg,
	}

	return o, nil
}

// OnLongPress Длительное замыкание
type OnLongPress struct {
	*messages.MessageImpl
}
