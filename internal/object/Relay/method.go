package Relay

import (
	"strings"
	"time"

	"github.com/pkg/errors"
	"touchon-server/internal/model"
	"touchon-server/internal/objects"
	"touchon-server/lib/events"
	_ "touchon-server/lib/events"
	"touchon-server/lib/events/object/relay"
	"touchon-server/lib/mqtt/messages"
)

func (o *RelayModel) On(args map[string]interface{}) ([]messages.Message, error) {
	portObjectID, err := o.GetProps().GetIntValue("address")
	if err != nil {
		return nil, err
	}

	portObj, err := objects.LoadPort(portObjectID, model.ChildTypeNobody)
	if err != nil {
		return nil, errors.Wrap(err, "getValues")
	}

	relayMsg, err := relay.NewOnStateMessage("object_manager/object/event", o.GetID())
	if err != nil {
		return nil, err
	}

	portMsg, err := portObj.On(nil)
	if err != nil {
		return nil, err
	}

	//TODO:: сделать тут сохранение статуса объекта в БД

	return append(portMsg, relayMsg), nil
}

func (o *RelayModel) Off(args map[string]interface{}) ([]messages.Message, error) {
	portObjectID, err := o.GetProps().GetIntValue("address")
	if err != nil {
		return nil, err
	}

	portObj, err := objects.LoadPort(portObjectID, model.ChildTypeNobody)
	if err != nil {
		return nil, errors.Wrap(err, "getValues")
	}

	relayMsg, err := relay.NewOffStateMessage("object_manager/object/event", o.GetID())
	if err != nil {
		return nil, err
	}

	portMsg, err := portObj.Off(nil)
	if err != nil {
		return nil, err
	}

	//TODO:: сделать тут сохранение статуса объекта в БД

	return append(portMsg, relayMsg), nil
}

func (o *RelayModel) Toggle(args map[string]interface{}) ([]messages.Message, error) {
	portObjectID, err := o.GetProps().GetIntValue("address")
	if err != nil {
		return nil, err
	}

	portObj, err := objects.LoadPort(portObjectID, model.ChildTypeNobody)
	if err != nil {
		return nil, errors.Wrap(err, "getValues")
	}

	portMsg, err := portObj.Toggle(nil)
	if err != nil {
		return nil, err
	}

	relayMsg, err := events.NewOnChangeStateMessage("object_manager/object/event", messages.TargetTypeObject, o.GetID(), strings.ToLower(portMsg[0].GetPayload()["state"].(string)), "")
	if err != nil {
		return nil, err
	}

	//TODO:: сделать тут сохранение статуса объекта в БД

	return append(portMsg, relayMsg), nil
}

func (o *RelayModel) Check(args map[string]interface{}) ([]messages.Message, error) {
	portObjectID, err := o.GetProps().GetIntValue("address")
	if err != nil {
		return nil, err
	}

	portObj, err := objects.LoadPort(portObjectID, model.ChildTypeNobody)
	if err != nil {
		return nil, errors.Wrap(err, "getValues")
	}

	stateRelay, err := portObj.GetPortState("get", nil, time.Duration(1)*time.Second)
	if err != nil {
		return nil, errors.Wrap(err, "getValues")
	}

	relayMsg, err := relay.NewCheckMessage("object_manager/object/event", messages.TargetTypeObject, o.GetID(), strings.ToLower(stateRelay), "")
	if err != nil {
		return nil, err
	}

	return []messages.Message{relayMsg}, nil
}
