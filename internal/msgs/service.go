// Пакет для работы с MQTT. Реализует прием и отправку сообщений в шину.

package msgs

import "touchon-server/lib/interfaces"

//import (
//	"bytes"
//	"encoding/json"
//	"net/http"
//	"time"
//
//	"github.com/pkg/errors"
//	"github.com/sirupsen/logrus"
//	"touchon-server/internal/context"
//	"touchon-server/internal/model"
//	"touchon-server/internal/scripts"
//	"touchon-server/internal/store"
//	"touchon-server/internal/ws"
//	"touchon-server/lib/event"
//	"touchon-server/lib/events"
//	mqttClient "touchon-server/lib/mqtt/client"
//	"touchon-server/lib/mqtt/messages"
//	"touchon-server/lib/mqtt/service"
//	"touchon-server/lib/subscribers"
//)
//
//type Service interface {
//	SetHandler(handler func(messages.Message) error)
//	GetLogger() *logrus.Logger
//	GetClient() mqttClient.Client
//	GetConfig() map[string]string
//	Start() error
//	Shutdown() error
//	Subscribe(publisher, topic string, msgType messages.MessageType, name string, targetType messages.TargetType, targetID *int, handler subscribers.MsgHandler) (int, error)
//	Unsubscribe(handlerIDs ...int)
//}

// Global instance
var I interfaces.MessagesService

//func New(bufferSize int, threads int, pushSenderAddress string) (Service, error) {
//	baseService, err := service.New(mqttClient.I, context.Config, bufferSize, threads, context.Logger)
//	if err != nil {
//		return nil, errors.Wrap(err, "service.New")
//	}
//
//	o := &ServiceImpl{
//		Service:           baseService,
//		subscribers:       subscribers.New(2000),
//		pushSenderAddress: pushSenderAddress,
//	}
//
//	o.Service.SetHandler(func(msg messages.Message) error {
//		// AR
//		go o.processObjectMessage(msg)
//
//		// TR
//		go func() {
//			if err := o.errorWrapper(msg); err != nil {
//				context.Logger.Error(errors.Wrap(err, "mqtt.Service"))
//			}
//		}()
//
//		for _, handler := range o.subscribers.GetHandlers(msg) {
//			msgs, err := handler(msg)
//			if err != nil {
//				o.GetLogger().Errorf("mqtt.Service.messageHandler: %v", err)
//
//				msg, err := events.NewOnErrorMessage("object_manager/error", msg.GetTargetType(), msg.GetTargetID(), err.Error())
//				if err != nil {
//					o.GetLogger().Error(errors.Wrap(err, "processObjectMessage"))
//					continue
//				}
//
//				msgs = append(msgs, msg)
//			}
//
//			for _, msg := range msgs {
//				if err := o.GetClient().Send(msg); err != nil {
//					o.GetLogger().Errorf("mqtt.Service.messageHandler: %v", err)
//				}
//			}
//		}
//
//		return nil
//	})
//
//	// Подписываем скрипты на обработку команд на их выполнение
//	o.scriptsHandlerID, err = o.Subscribe("", "", messages.MessageTypeCommand, "exec", messages.TargetTypeScript, nil, scripts.I.MqttMsgHandler)
//	if err != nil {
//		return nil, errors.Wrap(err, "service.New")
//	}
//
//	return o, nil
//}
//
//type ServiceImpl struct {
//	*service.Service
//	subscribers       *subscribers.Subscribers
//	scriptsHandlerID  int
//	pushSenderAddress string
//}
//
//func (o *ServiceImpl) Subscribe(publisher, topic string, msgType messages.MessageType, name string, targetType messages.TargetType, targetID *int, handler subscribers.MsgHandler) (int, error) {
//	handlerID, err := o.subscribers.AddHandler(publisher, topic, msgType, name, targetType, targetID, handler)
//	if err != nil {
//		return 0, errors.Wrap(err, "mqtt.Service.Subscribe")
//	}
//
//	return handlerID, nil
//}
//
//func (o *ServiceImpl) Unsubscribe(handlerIDs ...int) {
//	for _, id := range handlerIDs {
//		o.subscribers.DeleteHandler(id)
//	}
//}
//
//func (o *ServiceImpl) Shutdown() error {
//	o.Unsubscribe(o.scriptsHandlerID)
//	return o.Service.Shutdown()
//}
//
//// AR
//
//// processObjectMessage обработка топика объекта
//func (o *ServiceImpl) processObjectMessage(msg messages.Message) {
//	switch msg.GetType() {
//	case messages.MessageTypeEvent:
//		if err := o.processEvent(msg); err != nil {
//			o.GetLogger().Error(errors.Wrap(err, "processObjectMessage"))
//
//			msg, err := events.NewOnErrorMessage("action_router/error", msg.GetTargetType(), msg.GetTargetID(), err.Error())
//			if err != nil {
//				o.GetLogger().Error(errors.Wrap(err, "processObjectMessage"))
//				return
//			}
//
//			if err := o.GetClient().Send(msg); err != nil {
//				o.GetLogger().Error(errors.Wrap(err, "processObjectMessage"))
//			}
//		}
//
//	default:
//		topic := msg.GetTopic()
//		msg.SetTopic("")
//		o.GetLogger().Infof("MQTT: unhandled msg [%s] %s", topic, msg)
//	}
//}
//
//func (o *ServiceImpl) processEvent(msg messages.Message) error {
//	ev, err := event.FromMqttMessage(msg, false)
//	if err != nil {
//		return errors.Wrap(err, "processEvent")
//	}
//
//	storeEvent, err := store.I.EventsRepo().GetEvent(msg.GetTargetType(), msg.GetTargetID(), ev.Code)
//	if err != nil {
//		if errors.Is(err, store.ErrNotFound) {
//			return nil
//		}
//
//		return errors.Wrap(err, "processEvent")
//	}
//
//	actions, err := store.I.EventActionsRepo().GetActions(storeEvent.ID)
//	if err != nil {
//		return errors.Wrap(err, "processEvent")
//	}
//
//	for _, act := range actions[storeEvent.ID] {
//		if !act.Enabled {
//			continue
//		}
//
//		switch act.Type {
//		case model.ActionTypeDelay:
//			v, ok := act.Args["duration"]
//			if !ok {
//				return errors.Wrap(errors.New("duration not found"), "processEvent")
//			}
//
//			s, ok := v.(string)
//			if !ok {
//				return errors.Wrap(errors.New("duration is not string"), "processEvent")
//			}
//
//			d, err := time.ParseDuration(s)
//			if err != nil {
//				return errors.Wrap(err, "processEvent")
//			}
//
//			time.Sleep(d)
//
//		case model.ActionTypeMethod:
//			msg, err := messages.NewCommand(act.Name, act.TargetType, act.TargetID, act.Args)
//			if err != nil {
//				return errors.Wrap(err, "processEvent")
//			}
//
//			msg.SetTopic(string(msg.GetTargetType()) + "/method")
//
//			if err := o.GetClient().Send(msg); err != nil {
//				return errors.Wrap(err, "processEvent")
//			}
//
//		case model.ActionTypeNotification:
//			v, ok := act.Args["type"]
//			if !ok {
//				return errors.Wrap(errors.New("type not found"), "processEvent")
//			}
//
//			notType, ok := v.(string)
//			if !ok {
//				return errors.Wrap(errors.New("type is not string"), "processEvent")
//			}
//
//			v, ok = act.Args["text"]
//			if !ok {
//				return errors.Wrap(errors.New("text not found"), "processEvent")
//			}
//
//			notText, ok := v.(string)
//			if !ok {
//				return errors.Wrap(errors.New("type is not string"), "processEvent")
//			}
//
//			msg, err := messages.NewNotification(messages.NotificationType(notType), notText)
//			if err != nil {
//				return errors.Wrap(err, "processEvent")
//			}
//
//			if err := o.GetClient().SendRaw("notification", act.QoS, false, msg); err != nil {
//				return errors.Wrap(err, "processEvent")
//			}
//
//		default:
//			return errors.Wrap(errors.Errorf("unknown action type %q", act.Type), "processEvent")
//		}
//	}
//
//	return nil
//}
//
//// TR
//
//func (o *ServiceImpl) errorWrapper(msg messages.Message) error {
//	if err := o.messageHandler(msg); err != nil {
//		msg, err := events.NewOnErrorMessage("touchon-server/error", msg.GetTargetType(), msg.GetTargetID(), err.Error())
//		if err != nil {
//			return errors.Wrap(err, "errorWrapper")
//		}
//
//		if err := o.GetClient().Send(msg); err != nil {
//			return errors.Wrap(err, "errorWrapper")
//		}
//
//		return errors.Wrap(err, "errorWrapper")
//	}
//
//	return nil
//}
//
//// messageHandler Обработка сообщений, которые пришли в MQTT
//func (o *ServiceImpl) messageHandler(msg messages.Message) error {
//	switch {
//	case msg.GetType() == messages.MessageTypeEvent && msg.GetName() == "on_notify":
//		if err := o.processNotification(msg); err != nil {
//			return errors.Wrap(err, "messageHandler")
//		}
//
//	case msg.GetTargetType() == messages.TargetTypeObject && msg.GetType() == messages.MessageTypeEvent:
//		if err := o.processObjectEvent(msg); err != nil {
//			return errors.Wrap(err, "messageHandler")
//		}
//
//	case msg.GetTargetType() == messages.TargetTypeItem && msg.GetType() == messages.MessageTypeCommand:
//		switch msg.GetName() {
//		case "set_state":
//			state, ok := msg.GetPayload()["state"].(string)
//			if !ok {
//				return errors.Wrap(errors.New("state is not string"), "messageHandler")
//			}
//
//			if err := store.I.Items().ChangeItem(msg.GetTargetID(), state); err != nil {
//				return err
//			}
//
//			ws.I.Send(&model.ViewItem{
//				ID:     msg.GetTargetID(),
//				Status: state,
//			})
//		}
//
//	default:
//		topic := msg.GetTopic()
//		msg.SetTopic("")
//		o.GetLogger().Infof("MQTT: unhandled msg [%s] %s", topic, msg)
//	}
//
//	return nil
//}
//
//// processObjectEvent запуск метода объекта
//func (o *ServiceImpl) processObjectEvent(msg messages.Message) error {
//	// Ищем в таблице событие, которое пришло в топике
//	items, err := store.I.Items().GetItemsForChange(msg.GetTargetType(), msg.GetTargetID(), msg.GetName())
//	if err != nil {
//		return errors.Wrap(err, "processObjectEvent")
//	}
//
//	// Перебираем найденные итемы, чтобы произвести с ними действие
//	for _, item := range items {
//		switch item.Type {
//
//		case "button", "switch", "conditioner":
//			item.Status, _ = msg.GetStringValue(item.EventValue)
//			if err := store.I.Items().ChangeItem(item.ID, item.Status); err != nil {
//				return errors.Wrap(err, "processObjectEvent")
//			}
//
//			ws.I.Send(item)
//
//		case "sensor":
//			item.Value, _ = msg.GetFloatValue(item.EventValue)
//
//			sensor, err := store.I.Devices().GetSensor(item.ID)
//			if err != nil {
//				return errors.Wrap(err, "processObjectEvent")
//			}
//
//			if sensor.Current != item.Value {
//				ws.I.Send(item)
//			}
//
//			// Обновление значения в таблице сенсоров
//			if err := store.I.Devices().SetSensorValue(item.ID, item.Value); err != nil {
//				return errors.Wrap(err, "processObjectEvent")
//			}
//
//			// Обновление значения в таблице графиков для сенсора
//			t := time.Now().Format("2006-01-02T15:04")
//			if err := store.I.History().SetHourlyValue(item.ID, t, item.Value); err != nil {
//				return errors.Wrap(err, "processObjectEvent")
//			}
//		}
//
//		if o.GetLogger().Level >= logrus.DebugLevel {
//			itemToLog, err := json.Marshal(item)
//			if err != nil {
//				return errors.Wrap(err, "processObjectEvent")
//			}
//
//			o.GetLogger().Debug("Server wrote to WS: " + string(itemToLog))
//		}
//	}
//
//	return nil
//}
//
//func (o *ServiceImpl) processNotification(msg messages.Message) error {
//	text, err := msg.GetStringValue("msg")
//	if err != nil {
//		return errors.Wrap(err, "processNotification")
//	}
//
//	notifyType, err := msg.GetStringValue("type")
//	if err != nil {
//		return errors.Wrap(err, "processNotification")
//	}
//
//	notification := &model.Notification{
//		Text: text,
//		Type: events.NotifyType(notifyType),
//		Date: time.Now().Format("2006-01-02T15:04:05"),
//	}
//
//	if err := store.I.Notifications().AddNotification(notification); err != nil {
//		return errors.Wrap(err, "processNotification")
//	}
//
//	// Отправка сообщения через вебсокет
//	ws.I.Send(notification)
//
//	// Отправка критических сообщений в пуш уведомления
//	if notification.Type == events.NotifyTypeCritical {
//		tokens, err := store.I.Notifications().GetPushTokens()
//		if err != nil {
//			return errors.Wrap(err, "processNotification")
//		}
//
//		msg := &model.PushNotification{
//			Title:  "Важное уведомление!",
//			Body:   notification.Text,
//			Tokens: tokens,
//		}
//
//		data, err := json.Marshal(msg)
//		if err != nil {
//			return errors.Wrap(err, "processNotification")
//		}
//
//		resp, err := http.Post(o.pushSenderAddress+"/push", "application/json", bytes.NewReader(data))
//		if err != nil {
//			return errors.Wrap(err, "processNotification")
//		}
//
//		o.GetLogger().Debugf("Send POST request: %s/push ===> %q response: %s", o.pushSenderAddress, notification.Text, resp.Status)
//	}
//
//	return nil
//}
