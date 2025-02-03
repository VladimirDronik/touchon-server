// Реализация поиска подписчиков по параметрам сообщения

package subscribers

import (
	"strconv"

	"github.com/pkg/errors"
	"touchon-server/lib/event"
	"touchon-server/lib/intset"
	"touchon-server/lib/mqtt/messages"
)

type MqttMsgHandler func(msg messages.Message) ([]messages.Message, error)

func New(handlersCount int) *Subscribers {
	return &Subscribers{
		handlersCount: handlersCount,
		currHandlerID: 0,
		handlers:      make(map[int]MqttMsgHandler, handlersCount),
		columns: map[string]map[string]*intset.IntSet{
			"publishers":   make(map[string]*intset.IntSet, 20),
			"topics":       make(map[string]*intset.IntSet, 100),
			"types":        make(map[string]*intset.IntSet, 5),
			"names":        make(map[string]*intset.IntSet, 100),
			"target_types": make(map[string]*intset.IntSet, 10),
			"target_ids":   make(map[string]*intset.IntSet, 1000),
		},
	}
}

type Subscribers struct {
	handlersCount int
	currHandlerID int
	handlers      map[int]MqttMsgHandler // Все зарегистрированные обработчики
	columns       map[string]map[string]*intset.IntSet
}

func (o *Subscribers) AddHandler(publisher, topic string, msgType messages.MessageType, name string, targetType messages.TargetType, targetID *int, handler MqttMsgHandler) (int, error) {
	switch {
	case msgType != "" && msgType != messages.MessageTypeCommand && msgType != messages.MessageTypeEvent:
		return 0, errors.Wrap(errors.Errorf("message type is wrong %q", msgType), "Subscribers.AddHandler")
	case targetType != "" && !messages.TargetTypes[targetType]:
		return 0, errors.Wrap(errors.Errorf("target type is wrong %q", targetType), "Subscribers.AddHandler")
	case targetID != nil && *targetID < 1:
		return 0, errors.Wrap(errors.Errorf("target id is wrong %d", targetID), "Subscribers.AddHandler")
	case handler == nil:
		return 0, errors.Wrap(errors.New("handler is nil"), "Subscribers.AddHandler")
	}

	// Проверяем корректность имени события
	if msgType == messages.MessageTypeEvent && name != "" {
		if _, err := event.GetMaker(name); err != nil {
			return 0, errors.Wrap(err, "Subscribers.AddHandler")
		}
	}

	o.currHandlerID += 1

	o.handlers[o.currHandlerID] = handler

	tID := ""
	if targetID != nil && *targetID > 0 {
		tID = strconv.Itoa(*targetID)
	}

	values := map[string]string{
		"publishers":   publisher,
		"topics":       topic,
		"types":        string(msgType),
		"names":        name,
		"target_types": string(targetType),
		"target_ids":   tID,
	}

	for columnName, column := range o.columns {
		v := values[columnName]
		set := column[v]
		if set == nil {
			set = intset.New(o.handlersCount)
			column[v] = set
		}
		set.Add(o.currHandlerID)
	}

	return o.currHandlerID, nil
}

func (o *Subscribers) DeleteHandler(id int) {
	delete(o.handlers, id)

	for _, column := range o.columns {
		for _, set := range column {
			set.Remove(id)
		}
	}
}

func (o *Subscribers) GetHandlers(msg messages.Message) []MqttMsgHandler {
	values := map[string]string{
		"topics":       msg.GetTopic(),
		"types":        string(msg.GetType()),
		"names":        msg.GetName(),
		"target_types": string(msg.GetTargetType()),
		"target_ids":   strconv.Itoa(msg.GetTargetID()),
	}

	s := intset.New(o.handlersCount)

	for _, v := range []string{msg.GetPublisher(), ""} {
		if s2, ok := o.columns["publishers"][v]; ok {
			s.UnionWith(s2)
		}
	}

	if s.Len() == 0 {
		return nil
	}

	for columnName, v := range values {
		column := o.columns[columnName]

		s2 := intset.New(o.handlersCount)

		for _, v := range []string{v, ""} {
			if s3, ok := column[v]; ok {
				s2.UnionWith(s3)
			}
		}

		s.IntersectWith(s2)
	}

	r := make([]MqttMsgHandler, 0, 10)

	for _, handlerID := range s.Elems() {
		r = append(r, o.handlers[handlerID])
	}

	return r
}
