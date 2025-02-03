// Пакет для работы с MQTT. Реализует прием и отправку сообщений в шину.

package mqtt

import (
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"touchon-server/internal/context"
	"touchon-server/internal/model"
	"touchon-server/internal/mqtt/subscribers"
	"touchon-server/internal/scripts"
	"touchon-server/internal/store"
	"touchon-server/lib/event"
	"touchon-server/lib/events"
	mqttClient "touchon-server/lib/mqtt/client"
	"touchon-server/lib/mqtt/messages"
	"touchon-server/lib/mqtt/service"
)

type Service interface {
	SetHandler(handler func(messages.Message) error)
	GetLogger() *logrus.Logger
	GetClient() mqttClient.Client
	GetConfig() map[string]string
	Start() error
	Shutdown() error
	Subscribe(publisher, topic string, msgType messages.MessageType, name string, targetType messages.TargetType, targetID *int, handler subscribers.MqttMsgHandler) (int, error)
	Unsubscribe(handlerIDs ...int)
}

// Global instance
var I Service

func New(bufferSize int, threads int) (Service, error) {
	baseService, err := service.New(mqttClient.I, context.Config, bufferSize, threads, context.Logger)
	if err != nil {
		return nil, errors.Wrap(err, "service.New")
	}

	o := &ServiceImpl{
		Service:     baseService,
		subscribers: subscribers.New(2000),
	}

	o.Service.SetHandler(func(msg messages.Message) error {
		// AR
		go o.processObjectMessage(msg)

		for _, handler := range o.subscribers.GetHandlers(msg) {
			msgs, err := handler(msg)
			if err != nil {
				o.GetLogger().Errorf("mqtt.Service.messageHandler: %v", err)

				msg, err := events.NewOnErrorMessage("object_manager/error", msg.GetTargetType(), msg.GetTargetID(), err.Error())
				if err != nil {
					o.GetLogger().Error(errors.Wrap(err, "processObjectMessage"))
					continue
				}

				msgs = append(msgs, msg)
			}

			for _, msg := range msgs {
				if err := o.GetClient().Send(msg); err != nil {
					o.GetLogger().Errorf("mqtt.Service.messageHandler: %v", err)
				}
			}
		}

		return nil
	})

	// Подписываем скрипты на обработку команд на их выполнение
	o.scriptsHandlerID, err = o.Subscribe("", "", messages.MessageTypeCommand, "exec", messages.TargetTypeScript, nil, scripts.I.MqttMsgHandler)
	if err != nil {
		return nil, errors.Wrap(err, "service.New")
	}

	return o, nil
}

type ServiceImpl struct {
	*service.Service
	subscribers      *subscribers.Subscribers
	scriptsHandlerID int
}

func (o *ServiceImpl) Subscribe(publisher, topic string, msgType messages.MessageType, name string, targetType messages.TargetType, targetID *int, handler subscribers.MqttMsgHandler) (int, error) {
	handlerID, err := o.subscribers.AddHandler(publisher, topic, msgType, name, targetType, targetID, handler)
	if err != nil {
		return 0, errors.Wrap(err, "mqtt.Service.Subscribe")
	}

	return handlerID, nil
}

func (o *ServiceImpl) Unsubscribe(handlerIDs ...int) {
	for _, id := range handlerIDs {
		o.subscribers.DeleteHandler(id)
	}
}

func (o *ServiceImpl) Shutdown() error {
	o.Unsubscribe(o.scriptsHandlerID)
	return o.Service.Shutdown()
}

// AR

// processObjectMessage обработка топика объекта
func (o *ServiceImpl) processObjectMessage(msg messages.Message) {
	switch msg.GetType() {
	case messages.MessageTypeEvent:
		if err := o.processEvent(msg); err != nil {
			o.GetLogger().Error(errors.Wrap(err, "processObjectMessage"))

			msg, err := events.NewOnErrorMessage("action_router/error", msg.GetTargetType(), msg.GetTargetID(), err.Error())
			if err != nil {
				o.GetLogger().Error(errors.Wrap(err, "processObjectMessage"))
				return
			}

			if err := o.GetClient().Send(msg); err != nil {
				o.GetLogger().Error(errors.Wrap(err, "processObjectMessage"))
			}
		}

	default:
		topic := msg.GetTopic()
		msg.SetTopic("")
		o.GetLogger().Infof("MQTT: unhandled msg [%s] %s", topic, msg)
	}
}

func (o *ServiceImpl) processEvent(msg messages.Message) error {
	ev, err := event.FromMqttMessage(msg, false)
	if err != nil {
		return errors.Wrap(err, "processEvent")
	}

	storeEvent, err := store.I.EventsRepo().GetEvent(msg.GetTargetType(), msg.GetTargetID(), ev.Code)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return nil
		}

		return errors.Wrap(err, "processEvent")
	}

	actions, err := store.I.EventActionsRepo().GetActions(storeEvent.ID)
	if err != nil {
		return errors.Wrap(err, "processEvent")
	}

	for _, act := range actions[storeEvent.ID] {
		if !act.Enabled {
			continue
		}

		switch act.Type {
		case model.ActionTypeDelay:
			v, ok := act.Args["duration"]
			if !ok {
				return errors.Wrap(errors.New("duration not found"), "processEvent")
			}

			s, ok := v.(string)
			if !ok {
				return errors.Wrap(errors.New("duration is not string"), "processEvent")
			}

			d, err := time.ParseDuration(s)
			if err != nil {
				return errors.Wrap(err, "processEvent")
			}

			time.Sleep(d)

		case model.ActionTypeMethod:
			msg, err := messages.NewCommand(act.Name, act.TargetType, act.TargetID, act.Args)
			if err != nil {
				return errors.Wrap(err, "processEvent")
			}

			msg.SetTopic(o.GetConfig()["service_name"] + "/" + string(msg.GetTargetType()) + "/method")

			if err := o.GetClient().Send(msg); err != nil {
				return errors.Wrap(err, "processEvent")
			}

		case model.ActionTypeNotification:
			v, ok := act.Args["type"]
			if !ok {
				return errors.Wrap(errors.New("type not found"), "processEvent")
			}

			notType, ok := v.(string)
			if !ok {
				return errors.Wrap(errors.New("type is not string"), "processEvent")
			}

			v, ok = act.Args["text"]
			if !ok {
				return errors.Wrap(errors.New("text not found"), "processEvent")
			}

			notText, ok := v.(string)
			if !ok {
				return errors.Wrap(errors.New("type is not string"), "processEvent")
			}

			msg, err := messages.NewNotification(messages.NotificationType(notType), notText)
			if err != nil {
				return errors.Wrap(err, "processEvent")
			}

			if err := o.GetClient().SendRaw(o.GetConfig()["service_name"]+"/notification", act.QoS, false, msg); err != nil {
				return errors.Wrap(err, "processEvent")
			}

		default:
			return errors.Wrap(errors.Errorf("unknown action type %q", act.Type), "processEvent")
		}
	}

	return nil
}
