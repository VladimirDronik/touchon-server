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
	"touchon-server/lib/interfaces"
)

func (o *RelayModel) On(args map[string]interface{}) ([]interfaces.Message, error) {
	portObjectID, err := o.GetProps().GetIntValue("address")
	if err != nil {
		return nil, err
	}

	portObj, err := objects.LoadPort(portObjectID, model.ChildTypeNobody)
	if err != nil {
		return nil, errors.Wrap(err, "getValues")
	}

	relayMsg, err := relay.NewOnStateOn(o.GetID())
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

func (o *RelayModel) Off(args map[string]interface{}) ([]interfaces.Message, error) {
	portObjectID, err := o.GetProps().GetIntValue("address")
	if err != nil {
		return nil, err
	}

	portObj, err := objects.LoadPort(portObjectID, model.ChildTypeNobody)
	if err != nil {
		return nil, errors.Wrap(err, "getValues")
	}

	relayMsg, err := relay.NewOnStateOff(o.GetID())
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

func (o *RelayModel) Toggle(args map[string]interface{}) ([]interfaces.Message, error) {
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

	var state string
	if len(portMsg) > 0 {
		// TODO
		var payload map[string]interface{} // = portMsg[0].GetPayload()
		if v, ok := payload["state"]; ok {
			if v, ok := v.(string); ok {
				state = v
			}
		}
	}

	relayMsg, err := events.NewOnChangeState(interfaces.TargetTypeObject, o.GetID(), strings.ToLower(state), "")
	if err != nil {
		return nil, err
	}

	//TODO:: сделать тут сохранение статуса объекта в БД

	return append(portMsg, relayMsg), nil
}

func (o *RelayModel) Check(args map[string]interface{}) ([]interfaces.Message, error) {
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

	relayMsg, err := relay.NewOnCheck(o.GetID(), strings.ToLower(stateRelay), "")
	if err != nil {
		return nil, err
	}

	return []interfaces.Message{relayMsg}, nil
}
