package events

import (
	"github.com/VladimirDronik/touchon-server/mqtt/messages"
	"github.com/pkg/errors"
)

func NewOnPortCheck(targetID int, targetType messages.TargetType, state, value string) (messages.Message, error) {
	payload := map[string]interface{}{"state": state, "value": value}

	impl, err := NewEvent("onPortCheck", targetID, targetType, payload)
	if err != nil {
		return nil, errors.Wrap(err, "NewOnPortCheck")
	}

	return &OnPortCheck{Message: impl}, nil
}

type OnPortCheck struct {
	messages.Message
}

func (o *OnPortCheck) GetState() (string, error) {
	v, err := o.GetStringValue("state")
	if err != nil {
		return "", errors.Wrap(err, "GetState")
	}

	return v, nil
}

func (o *OnPortCheck) GetValue() (string, error) {
	v, err := o.GetStringValue("value")
	if err != nil {
		return "", errors.Wrap(err, "GetValue")
	}

	return v, nil
}

func (o *OnPortCheck) SetState(v string) {
	o.GetPayload()["state"] = v
}

func (o *OnPortCheck) SetValue(v string) {
	o.GetPayload()["value"] = v
}
