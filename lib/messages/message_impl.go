package messages

import (
	"github.com/pkg/errors"
	"touchon-server/lib/interfaces"
)

func NewCommand(name string, targetType interfaces.TargetType, targetID int, args map[string]interface{}) (interfaces.Message, error) {
	msg, err := NewMessage(interfaces.MessageTypeCommand, name, targetType, targetID)
	if err != nil {
		return nil, errors.Wrap(err, "NewCommand")
	}

	o := &Command{
		MessageImpl: msg,
		Args:        args,
	}

	return o, nil
}

type Command struct {
	*MessageImpl
	Args map[string]interface{} `json:"args"` // Аргументы вызова
}

func NewEvent(targetType interfaces.TargetType, targetID int) (*MessageImpl, error) {
	msg, err := NewMessage(interfaces.MessageTypeEvent, "", targetType, targetID)
	if err != nil {
		return nil, errors.Wrap(err, "NewEvent")
	}

	return msg, nil
}

func NewMessage(msgType interfaces.MessageType, name string, targetType interfaces.TargetType, targetID int) (*MessageImpl, error) {
	switch {
	case msgType != interfaces.MessageTypeCommand && msgType != interfaces.MessageTypeEvent:
		return nil, errors.Wrap(errors.Errorf("unknown message type %q", msgType), "NewMessage")
	case targetID < 0:
		return nil, errors.Wrap(errors.New("object ID < 0"), "NewMessage")
	}

	if _, ok := interfaces.TargetTypes[targetType]; !ok {
		return nil, errors.Wrap(errors.Errorf("unknown target type %q", targetType), "NewMessage")
	}

	return &MessageImpl{
		Type:       msgType,
		Name:       name,
		TargetType: targetType,
		TargetID:   targetID,
	}, nil
}

type MessageImpl struct {
	Type       interfaces.MessageType `json:"type"`        // event,command
	Name       string                 `json:"name"`        // onChange,check
	TargetType interfaces.TargetType  `json:"target_type"` //
	TargetID   int                    `json:"target_id"`   // 82
}

func (o *MessageImpl) GetType() interfaces.MessageType {
	return o.Type
}

func (o *MessageImpl) GetName() string {
	return o.Name
}

func (o *MessageImpl) GetTargetType() interfaces.TargetType {
	return o.TargetType
}

func (o *MessageImpl) GetTargetID() int {
	return o.TargetID
}

func (o *MessageImpl) SetType(v interfaces.MessageType) {
	o.Type = v
}

func (o *MessageImpl) SetName(v string) {
	o.Name = v
}

func (o *MessageImpl) SetTargetType(v interfaces.TargetType) {
	o.TargetType = v
}

func (o *MessageImpl) SetTargetID(v int) {
	o.TargetID = v
}
