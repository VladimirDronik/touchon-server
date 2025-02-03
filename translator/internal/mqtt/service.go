// Пакет для работы с MQTT. Реализует прием и отправку сообщений в шину.

package mqtt

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"github.com/VladimirDronik/touchon-server/events"
	"github.com/VladimirDronik/touchon-server/mqtt/client"
	"github.com/VladimirDronik/touchon-server/mqtt/messages"
	"github.com/VladimirDronik/touchon-server/mqtt/service"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"translator/internal/model"
	"translator/internal/store"
	"translator/internal/store/sqlstore"
	"translator/internal/ws"
)

func New(client *client.Client, cfg map[string]string, bufferSize int, threads int, store *sqlstore.Store, wsServer *ws.Server, pushSenderAddress string, logger *logrus.Logger) (*Service, error) {
	baseService, err := service.New(client, cfg, bufferSize, threads, logger)
	if err != nil {
		return nil, errors.Wrap(err, "service.New")
	}

	o := &Service{
		Service:           baseService,
		store:             store,
		wsServer:          wsServer,
		pushSenderAddress: pushSenderAddress,
	}

	o.Service.SetHandler(o.errorWrapper)

	return o, nil
}

type Service struct {
	*service.Service
	store             store.Store
	wsServer          *ws.Server
	pushSenderAddress string
}

func (o *Service) errorWrapper(msg messages.Message) error {
	if err := o.messageHandler(msg); err != nil {
		msg, err := events.NewOnErrorMessage("translator/error", msg.GetTargetType(), msg.GetTargetID(), err.Error())
		if err != nil {
			return errors.Wrap(err, "errorWrapper")
		}

		if err := o.GetClient().Send(msg); err != nil {
			return errors.Wrap(err, "errorWrapper")
		}

		return errors.Wrap(err, "errorWrapper")
	}

	return nil
}

// messageHandler Обработка сообщений, которые пришли в MQTT
func (o *Service) messageHandler(msg messages.Message) error {
	switch {
	case msg.GetType() == messages.MessageTypeEvent && msg.GetName() == "on_notify":
		if err := o.processNotification(msg); err != nil {
			return errors.Wrap(err, "messageHandler")
		}

	case msg.GetTargetType() == messages.TargetTypeObject && msg.GetType() == messages.MessageTypeEvent:
		if err := o.processObjectEvent(msg); err != nil {
			return errors.Wrap(err, "messageHandler")
		}

	case msg.GetTargetType() == messages.TargetTypeItem && msg.GetType() == messages.MessageTypeCommand:
		switch msg.GetName() {
		case "set_state":
			state, ok := msg.GetPayload()["state"].(string)
			if !ok {
				return errors.Wrap(errors.New("state is not string"), "messageHandler")
			}

			if err := o.store.Items().ChangeItem(msg.GetTargetID(), state); err != nil {
				return err
			}

			o.wsServer.Send(&model.ViewItem{
				ID:     msg.GetTargetID(),
				Status: state,
			})
		}

	default:
		topic := msg.GetTopic()
		msg.SetTopic("")
		o.GetLogger().Infof("MQTT: unhandled msg [%s] %s", topic, msg)
	}

	return nil
}

// processObjectEvent запуск метода объекта
func (o *Service) processObjectEvent(msg messages.Message) error {
	// Ищем в таблице событие, которое пришло в топике
	items, err := o.store.Items().GetItemsForChange(msg.GetTargetType(), msg.GetTargetID(), msg.GetName())
	if err != nil {
		return errors.Wrap(err, "processObjectEvent")
	}

	// Перебираем найденные итемы, чтобы произвести с ними действие
	for _, item := range items {
		switch item.Type {

		case "button", "switch", "conditioner":
			item.Status, _ = msg.GetStringValue(item.EventValue)
			if err := o.store.Items().ChangeItem(item.ID, item.Status); err != nil {
				return errors.Wrap(err, "processObjectEvent")
			}

			o.wsServer.Send(item)

		case "sensor":
			item.Value, _ = msg.GetFloatValue(item.EventValue)

			sensor, err := o.store.Devices().GetSensor(item.ID)
			if err != nil {
				return errors.Wrap(err, "processObjectEvent")
			}

			if sensor.Current != item.Value {
				o.wsServer.Send(item)
			}

			// Обновление значения в таблице сенсоров
			if err := o.store.Devices().SetSensorValue(item.ID, item.Value); err != nil {
				return errors.Wrap(err, "processObjectEvent")
			}

			// Обновление значения в таблице графиков для сенсора
			t := time.Now().Format("2006-01-02T15:04")
			if err := o.store.History().SetHourlyValue(item.ID, t, item.Value); err != nil {
				return errors.Wrap(err, "processObjectEvent")
			}
		}

		if o.GetLogger().Level >= logrus.DebugLevel {
			itemToLog, err := json.Marshal(item)
			if err != nil {
				return errors.Wrap(err, "processObjectEvent")
			}

			o.GetLogger().Debug("Server wrote to WS: " + string(itemToLog))
		}
	}

	return nil
}

func (o *Service) processNotification(msg messages.Message) error {
	text, err := msg.GetStringValue("msg")
	if err != nil {
		return errors.Wrap(err, "processNotification")
	}

	notifyType, err := msg.GetStringValue("type")
	if err != nil {
		return errors.Wrap(err, "processNotification")
	}

	notification := &model.Notification{
		Text: text,
		Type: events.NotifyType(notifyType),
		Date: time.Now().Format("2006-01-02T15:04:05"),
	}

	if err := o.store.Notifications().AddNotification(notification); err != nil {
		return errors.Wrap(err, "processNotification")
	}

	// Отправка сообщения через вебсокет
	o.wsServer.Send(notification)

	// Отправка критических сообщений в пуш уведомления
	if notification.Type == events.NotifyTypeCritical {
		tokens, err := o.store.Notifications().GetPushTokens()
		if err != nil {
			return errors.Wrap(err, "processNotification")
		}

		msg := &model.PushNotification{
			Title:  "Важное уведомление!",
			Body:   notification.Text,
			Tokens: tokens,
		}

		data, err := json.Marshal(msg)
		if err != nil {
			return errors.Wrap(err, "processNotification")
		}

		resp, err := http.Post(o.pushSenderAddress+"/push", "application/json", bytes.NewReader(data))
		if err != nil {
			return errors.Wrap(err, "processNotification")
		}

		o.GetLogger().Debugf("Send POST request: %s/push ===> %q response: %s", o.pushSenderAddress, notification.Text, resp.Status)
	}

	return nil
}
