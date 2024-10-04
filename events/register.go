//go:build ignore

package events

import "github.com/pkg/errors"

type EventMaker func() (Event, error)

var register = make(map[string]EventMaker, 20)

func RegisterEvent(eventName string, maker EventMaker) error {
	if _, ok := register[eventName]; ok {
		return errors.Wrap(errors.Errorf("event %q is exists", eventName), "RegisterEvent")
	}

	register[eventName] = maker

	return nil
}

func GetEventMaker(eventName string) (EventMaker, error) {
	maker, ok := register[eventName]
	if !ok {
		return nil, errors.Wrap(errors.Errorf("event %q not registered", eventName), "GetEventMaker")
	}

	return maker, nil
}
