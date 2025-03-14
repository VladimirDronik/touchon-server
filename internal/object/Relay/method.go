package Relay

import (
	"strconv"
	"strings"
	"time"
	"touchon-server/internal/model"

	"github.com/pkg/errors"
	"touchon-server/internal/objects"
	"touchon-server/lib/events"
	"touchon-server/lib/events/object/relay"
	"touchon-server/lib/interfaces"
)

func (o *RelayModel) On(args map[string]interface{}) ([]interfaces.Message, error) {
	portObjectID, err := o.GetProps().GetIntValue("address")
	if err != nil {
		return nil, errors.Wrap(err, "RelayModel.On")
	}

	portObj, err := objects.LoadPort(portObjectID, false)
	if err != nil {
		return nil, errors.Wrap(err, "RelayModel.On")
	}

	relayMsg, err := relay.NewOnStateOn(o.GetID())
	if err != nil {
		return nil, errors.Wrap(err, "RelayModel.On")
	}

	//Отправляем стандартное сообщение по смене состояния для объекта
	relayStateMsg, err := events.NewOnChangeState(interfaces.TargetTypeObject, o.GetID(), "on", "")
	if err != nil {
		return nil, errors.Wrap(err, "RelayModel.On")
	}

	portMsg, err := portObj.On(nil)
	if err != nil {
		//return nil, errors.Wrap(err, "RelayModel.On")
		relayMsg, err = events.NewOnError(interfaces.TargetTypeObject, o.GetID(), "relay not available")
		if err != nil {
			return nil, errors.Wrap(err, "RelayModel.On")
		}

		relayStateMsg, err = events.NewOnChangeState(interfaces.TargetTypeObject, o.GetID(), model.StatusUnavailable, "")
		if err != nil {
			return nil, errors.Wrap(err, "RelayModel.On")
		}
	}

	//TODO:: сделать тут сохранение статуса объекта в БД

	return append(portMsg, relayMsg, relayStateMsg), nil
}

func (o *RelayModel) Off(args map[string]interface{}) ([]interfaces.Message, error) {
	portObjectID, err := o.GetProps().GetIntValue("address")
	if err != nil {
		return nil, err
	}

	portObj, err := objects.LoadPort(portObjectID, false)
	if err != nil {
		return nil, errors.Wrap(err, "RelayModel.Off")
	}

	relayMsg, err := relay.NewOnStateOff(o.GetID())
	if err != nil {
		return nil, errors.Wrap(err, "RelayModel.Off")
	}

	//Отправляем стандартное сообщение по смене состояния для объекта
	relayStateMsg, err := events.NewOnChangeState(interfaces.TargetTypeObject, o.GetID(), "off", "")
	if err != nil {
		return nil, errors.Wrap(err, "RelayModel.Toggle")
	}

	portMsg, err := portObj.Off(nil)
	if err != nil {
		return nil, errors.Wrap(err, "RelayModel.Off")
	}

	//TODO:: сделать тут сохранение статуса объекта в БД

	return append(portMsg, relayMsg, relayStateMsg), nil
}

func (o *RelayModel) Toggle(args map[string]interface{}) ([]interfaces.Message, error) {
	portObjectID, err := o.GetProps().GetIntValue("address")
	if err != nil {
		return nil, errors.Wrap(err, "RelayModel.Toggle")
	}

	portObj, err := objects.LoadPort(portObjectID, false)
	if err != nil {
		return nil, errors.Wrap(err, "RelayModel.Toggle")
	}

	//Отправляем стандартное сообщение по смене состояния для объекта
	portMsg, err := portObj.Toggle(nil)
	if err != nil {
		return nil, errors.Wrap(err, "RelayModel.Toggle")
	}

	var state string
	if len(portMsg) > 0 && portMsg[0].GetPayload() != nil {
		if v, ok := portMsg[0].GetPayload()["state"]; ok {
			if v, ok := v.(string); ok {
				state = v
			}
		}
	}

	state = strings.ToLower(state)

	relayMsg, err := events.NewOnChangeState(interfaces.TargetTypeObject, o.GetID(), state, "")
	if err != nil {
		return nil, errors.Wrap(err, "RelayModel.Toggle")
	}

	var relayStateMsq interfaces.Message
	switch state {
	case "on":
		relayStateMsq, err = relay.NewOnStateOn(o.GetID())
		if err != nil {
			return nil, errors.Wrap(err, "RelayModel.On")
		}
		break
	case "off":
		relayStateMsq, err = relay.NewOnStateOff(o.GetID())
		if err != nil {
			return nil, errors.Wrap(err, "RelayModel.Off")
		}
		break
	default:
		return nil, errors.Wrap(err, "Relay "+strconv.Itoa(o.GetID())+": the device did not return a state")
	}

	//TODO:: сделать тут сохранение статуса объекта в БД

	return append(portMsg, relayMsg, relayStateMsq), nil
}

func (o *RelayModel) Check(args map[string]interface{}) ([]interfaces.Message, error) {
	portObjectID, err := o.GetProps().GetIntValue("address")
	if err != nil {
		return nil, errors.Wrap(err, "RelayModel.check")
	}

	portObj, err := objects.LoadPort(portObjectID, false)
	if err != nil {
		return nil, errors.Wrap(err, "RelayModel.check")
	}

	stateRelay, err := portObj.GetPortState("get", nil, time.Duration(1)*time.Second)
	if err != nil {
		return nil, errors.Wrap(err, "RelayModel.check")
	}

	relayMsg, err := relay.NewOnCheck(o.GetID(), strings.ToLower(stateRelay), "")
	if err != nil {
		return nil, errors.Wrap(err, "RelayModel.check")
	}

	return []interfaces.Message{relayMsg}, nil
}
