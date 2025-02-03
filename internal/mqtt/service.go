// Пакет для работы с MQTT. Реализует прием и отправку сообщений в шину.

package mqtt

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"touchon-server/internal/context"
	"touchon-server/internal/mqtt/subscribers"
	"touchon-server/internal/scripts"
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
