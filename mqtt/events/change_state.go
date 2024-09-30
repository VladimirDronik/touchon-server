package events

import (
	"github.com/VladimirDronik/touchon-server/mqtt/messages"
	"github.com/pkg/errors"
)

func NewOnChangeState(targetID int, targetType messages.TargetType, state, value string) (messages.Message, error) {
	payload := map[string]interface{}{"state": state, "value": value}

	impl, err := NewEvent("onChangeState", targetID, targetType, payload)
	if err != nil {
		return nil, errors.Wrap(err, "NewOnChangeState")
	}

	return &OnChangeState{Message: impl}, nil
}

type OnChangeState struct {
	messages.Message
}

func (o *OnChangeState) GetState() (string, error) {
	v, err := o.GetStringValue("state")
	if err != nil {
		return "", errors.Wrap(err, "GetState")
	}

	return v, nil
}

func (o *OnChangeState) GetValue() (string, error) {
	v, err := o.GetStringValue("value")
	if err != nil {
		return "", errors.Wrap(err, "GetValue")
	}

	return v, nil
}

func (o *OnChangeState) SetState(v string) {
	o.GetPayload()["state"] = v
}

func (o *OnChangeState) SetValue(v string) {
	o.GetPayload()["value"] = v
}
