package item

import (
	"github.com/pkg/errors"
	"touchon-server/lib/interfaces"
	"touchon-server/lib/messages"
)

func NewOnChangeStateOn(itemID int) (interfaces.Message, error) {
	msg, err := messages.NewEvent(interfaces.TargetTypeItem, itemID)
	if err != nil {
		return nil, errors.Wrap(err, "NewOnChangeStateOn")
	}

	o := &OnChangeStateOn{
		MessageImpl: msg,
	}

	return o, nil
}

type OnChangeStateOn struct {
	*messages.MessageImpl
}

func NewOnChangeStateOff(itemID int) (interfaces.Message, error) {
	msg, err := messages.NewEvent(interfaces.TargetTypeItem, itemID)
	if err != nil {
		return nil, errors.Wrap(err, "NewOnChangeStateOff")
	}

	o := &OnChangeStateOff{
		MessageImpl: msg,
	}

	return o, nil
}

type OnChangeStateOff struct {
	*messages.MessageImpl
}
