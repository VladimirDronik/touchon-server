//go:build ignore

package events

import (
	"encoding/json"

	"github.com/VladimirDronik/touchon-server/mqtt/messages"
	"github.com/pkg/errors"
	"github.com/valyala/fastjson"
)

type Event interface {
	GetName() string                    // onChange,check
	GetTargetID() int                   // 82
	GetTargetType() messages.TargetType // object,item
	GetPayload() map[string]interface{} //

	SetName(string)
	SetTargetID(int)
	SetTargetType(messages.TargetType)
	SetPayload(map[string]interface{})

	GetFloatValue(name string) (float32, error)
	GetStringValue(name string) (string, error)
	GetIntValue(name string) (int, error)
	GetBoolValue(name string) (bool, error)

	json.Marshaler
	json.Unmarshaler
}

func MakeEventFromJSON(data []byte) (Event, error) {
	eventName := fastjson.GetString(data, "name")

	eventMaker, err := GetEventMaker(eventName)
	if err != nil {
		return nil, errors.Wrap(err, "MakeEventFromJSON")
	}

	r, err := eventMaker()
	if err != nil {
		return nil, errors.Wrap(err, "MakeEventFromJSON")
	}

	if err := json.Unmarshal(data, &r); err != nil {
		return nil, errors.Wrap(err, "MakeEventFromJSON")
	}

	if r.GetName() != eventName {
		return nil, errors.Wrap(errors.Errorf("logic error: event name %q not equal event name %q from json", r.GetName(), eventName), "MakeEventFromJSON")
	}

	return r, nil
}
