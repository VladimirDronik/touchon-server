package event

import "github.com/pkg/errors"

type Maker func() (*Event, error)

var register = make(map[string]Maker, 20)

func Register(maker Maker) error {
	if maker == nil {
		return errors.Wrap(errors.New("maker is nil"), "Register")
	}

	e, err := maker()
	if err != nil {
		return errors.Wrap(err, "event.Register")
	}

	if _, ok := register[e.Code]; ok {
		return errors.Wrap(errors.Errorf("event %q is exists", e.Code), "Register")
	}

	if err := e.Check(); err != nil {
		return errors.Wrap(err, "Register")
	}

	register[e.Code] = maker

	return nil
}

func GetMaker(eventName string) (Maker, error) {
	maker, ok := register[eventName]
	if !ok {
		return nil, errors.Wrap(errors.Errorf("event %q not registered", eventName), "GetMaker")
	}

	return maker, nil
}

func GetEvent(eventName string) (*Event, error) {
	maker, ok := register[eventName]
	if !ok {
		return nil, errors.Wrap(errors.Errorf("event %q not registered", eventName), "GetEvent")
	}

	event, err := maker()
	if err != nil {
		return nil, errors.Wrap(err, "GetEvent")
	}

	return event, nil
}
