package messages

import (
	"github.com/pkg/errors"
	"touchon-server/lib/interfaces"
)

func NewCommand(name string, targetType interfaces.TargetType, targetID int, args map[string]interface{}) (interfaces.Command, error) {
	msg, err := NewMessage(interfaces.MessageTypeCommand, name, targetType, targetID)
	if err != nil {
		return nil, errors.Wrap(err, "NewCommand")
	}

	o := &CommandImpl{
		Message: msg,
	}

	o.SetPayload(args)

	return o, nil
}

type CommandImpl struct {
	interfaces.Message
}

func (o *CommandImpl) GetArgs() map[string]interface{} {
	return o.GetPayload()
}

func NewEvent(name string, targetType interfaces.TargetType, targetID int) (*MessageImpl, error) {
	msg, err := NewMessage(interfaces.MessageTypeEvent, name, targetType, targetID)
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
		msgType:    msgType,
		name:       name,
		targetType: targetType,
		targetID:   targetID,
	}, nil
}

type MessageImpl struct {
	msgType    interfaces.MessageType // event,command
	name       string                 // onChange,check
	targetType interfaces.TargetType  //
	targetID   int                    // 82
	payload    map[string]interface{} //
}

func (o *MessageImpl) GetType() interfaces.MessageType {
	return o.msgType
}

func (o *MessageImpl) GetName() string {
	return o.name
}

func (o *MessageImpl) GetTargetType() interfaces.TargetType {
	return o.targetType
}

func (o *MessageImpl) GetTargetID() int {
	return o.targetID
}

func (o *MessageImpl) SetType(v interfaces.MessageType) {
	o.msgType = v
}

func (o *MessageImpl) SetName(v string) {
	o.name = v
}

func (o *MessageImpl) SetTargetType(v interfaces.TargetType) {
	o.targetType = v
}

func (o *MessageImpl) SetTargetID(v int) {
	o.targetID = v
}

func (o *MessageImpl) GetValue(k string) interface{} {
	return o.payload[k]
}

func (o *MessageImpl) SetValue(k string, v interface{}) {
	if o.payload == nil {
		o.payload = make(map[string]interface{})
	}

	o.payload[k] = v
}

func (o *MessageImpl) GetPayload() map[string]interface{} {
	return o.payload
}

func (o *MessageImpl) SetPayload(v map[string]interface{}) {
	if o.payload == nil {
		o.payload = make(map[string]interface{}, len(v))
	}

	for k, v := range v {
		o.payload[k] = v
	}
}

func (o *MessageImpl) GetFloatValue(name string) (float32, error) {
	v, ok := o.GetPayload()[name]
	if !ok {
		return 0, errors.Wrap(errors.Errorf("%s not found", name), "GetFloatValue")
	}

	switch v := v.(type) {
	case float32:
		return v, nil
	case float64:
		return float32(v), nil
	case int:
		return float32(v), nil
	default:
		return 0, errors.Wrap(errors.Errorf("unexpected data type %T", v), "GetFloatValue")
	}
}

func (o *MessageImpl) GetStringValue(name string) (string, error) {
	v, ok := o.GetPayload()[name]
	if !ok {
		return "", errors.Wrap(errors.Errorf("%s not found", name), "GetStringValue")
	}

	switch v := v.(type) {
	case string:
		return v, nil
	default:
		return "", errors.Wrap(errors.Errorf("unexpected data type %T", v), "GetStringValue")
	}
}

func (o *MessageImpl) GetIntValue(name string) (int, error) {
	v, ok := o.GetPayload()[name]
	if !ok {
		return 0, errors.Wrap(errors.Errorf("%s not found", name), "GetIntValue")
	}

	switch v := v.(type) {
	case float32:
		return int(v), nil
	case float64:
		return int(v), nil
	case int:
		return v, nil
	default:
		return 0, errors.Wrap(errors.Errorf("unexpected data type %T", v), "GetIntValue")
	}
}

func (o *MessageImpl) GetBoolValue(name string) (bool, error) {
	v, ok := o.GetPayload()[name]
	if !ok {
		return false, errors.Wrap(errors.Errorf("%s not found", name), "GetBoolValue")
	}

	switch v := v.(type) {
	case bool:
		return v, nil
	default:
		return false, errors.Wrap(errors.Errorf("unexpected data type %T", v), "GetBoolValue")
	}
}
