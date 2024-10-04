package event

import (
	"github.com/VladimirDronik/touchon-server/models"
	"github.com/VladimirDronik/touchon-server/mqtt/messages"
	"github.com/pkg/errors"
)

type Event struct {
	Code        string `json:"code"` // unique
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Props       *Props `json:"props"`

	TargetID   int                 `json:"target_id"`
	TargetType messages.TargetType `json:"target_type"`
}

func (o *Event) Check() error {
	if _, ok := messages.TargetTypes[o.TargetType]; !ok {
		return errors.Wrap(errors.Errorf("unknown target type %q", o.TargetType), "Event.Check")
	}

	switch {
	case o.Code == "":
		return errors.Wrap(errors.New("code is empty"), "Event.Check")
	case o.Name == "":
		return errors.Wrap(errors.New("name is empty"), "Event.Check")
	case o.Props == nil:
		return errors.Wrap(errors.New("props is empty"), "Event.Check")
	}

	if err := o.Props.Check(); err != nil {
		return errors.Wrap(err, "Event.Check")
	}

	return nil
}

func (o *Event) ToMqttMessage(topic string) (messages.Message, error) {
	payload := make(map[string]interface{}, o.Props.Len())
	for _, p := range o.Props.GetOrderedMap().GetValueList() {
		payload[p.Code] = p.GetValue()
	}

	m, err := messages.NewMessage(messages.MessageTypeEvent, o.Code, o.TargetID, o.TargetType, payload)
	if err != nil {
		return nil, errors.Wrap(err, "ToMqttMessage")
	}

	m.SetTopic(topic)

	return m, nil
}

func FromMqttMessage(msg messages.Message, parseUnknownEvent bool) (*Event, error) {
	maker, err := GetMaker(msg.GetName())
	if !parseUnknownEvent && err != nil {
		// Ругаемся на незарегистрированное событие
		return nil, errors.Wrap(err, "FromMqttMessage")
	}

	if err == nil {
		// Разбираем известное (зарегистрированное) событие
		e, err := maker()
		if err != nil {
			return nil, errors.Wrap(err, "FromMqttMessage")
		}

		e.Code = msg.GetName()
		e.TargetID = msg.GetTargetID()
		e.TargetType = msg.GetTargetType()

		for k, v := range msg.GetPayload() {
			p, err := e.Props.Get(k)
			if err != nil {
				return nil, errors.Wrap(err, "FromMqttMessage")
			}

			if err := p.SetValue(v); err != nil {
				return nil, errors.Wrap(err, "FromMqttMessage")
			}
		}

		return e, nil
	} else {
		// Разбираем неизвестное (не зарегистрированное) событие
		e := &Event{Props: NewProps()}

		e.Code = msg.GetName()
		e.TargetID = msg.GetTargetID()
		e.TargetType = msg.GetTargetType()

		for k, v := range msg.GetPayload() {
			p := Prop{
				Code: k,
				Item: &models.Item{},
			}

			if err := p.SetValue(v); err != nil {
				return nil, errors.Wrap(err, "FromMqttMessage")
			}
		}

		return e, nil
	}
}
