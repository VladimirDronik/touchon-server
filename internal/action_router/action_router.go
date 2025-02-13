package action_router

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"touchon-server/internal/context"
	"touchon-server/internal/model"
	"touchon-server/internal/store"
	"touchon-server/internal/ws"
	"touchon-server/lib/events"
	"touchon-server/lib/events/item"
	"touchon-server/lib/interfaces"
	"touchon-server/lib/messages"
	msgs "touchon-server/lib/messages"
)

// Global instance
var I *Service

func New() *Service {
	return &Service{}
}

type Service struct {
	handlerIDs []int
}

func (o *Service) Start() error {
	handlerID, err := msgs.I.Subscribe(interfaces.MessageTypeEvent, "", "", nil, o.actionRouter)
	if err != nil {
		return errors.Wrap(err, "action_router.Service.Start")
	}
	o.handlerIDs = append(o.handlerIDs, handlerID)

	handlerID, err = msgs.I.Subscribe(interfaces.MessageTypeEvent, "", interfaces.TargetTypeObject, nil, o.processObjectEvent)
	if err != nil {
		return errors.Wrap(err, "action_router.Service.Start")
	}
	o.handlerIDs = append(o.handlerIDs, handlerID)

	handlerID, err = msgs.I.Subscribe(interfaces.MessageTypeEvent, "", interfaces.TargetTypeItem, nil, o.processItemEvents)
	if err != nil {
		return errors.Wrap(err, "action_router.Service.Start")
	}
	o.handlerIDs = append(o.handlerIDs, handlerID)

	handlerID, err = msgs.I.Subscribe(interfaces.MessageTypeNotification, "", "", nil, o.processNotification)
	if err != nil {
		return errors.Wrap(err, "action_router.Service.Start")
	}
	o.handlerIDs = append(o.handlerIDs, handlerID)

	return nil
}

func (o *Service) Shutdown() error {
	msgs.I.Unsubscribe(o.handlerIDs...)

	return nil
}

func (o *Service) actionRouter(msg interfaces.Message) {
	ev, ok := msg.(interfaces.Event)
	if !ok {
		context.Logger.Error(errors.Wrap(errors.Errorf("msg [%s, %s, %s, %d] is not event", msg.GetType(), msg.GetName(), msg.GetTargetType(), msg.GetTargetID()), "action_router.Service.msgHandler"))
		return
	}

	storeEvent, err := store.I.EventsRepo().GetEvent(ev.GetTargetType(), ev.GetTargetID(), ev.GetEventCode())
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return
		}

		context.Logger.Error(errors.Wrap(err, "action_router.Service.msgHandler"))
		return
	}

	actions, err := store.I.EventActionsRepo().GetActions(storeEvent.ID)
	if err != nil {
		context.Logger.Error(errors.Wrap(err, "action_router.Service.msgHandler"))
		return
	}

	for _, act := range actions[storeEvent.ID] {
		if !act.Enabled {
			continue
		}

		switch act.Type {
		case model.ActionTypeDelay:
			v, ok := act.Args["duration"]
			if !ok {
				context.Logger.Error(errors.Wrap(errors.New("duration not found"), "action_router.Service.msgHandler"))
				return
			}

			s, ok := v.(string)
			if !ok {
				context.Logger.Error(errors.Wrap(errors.New("duration is not string"), "action_router.Service.msgHandler"))
				return
			}

			d, err := time.ParseDuration(s)
			if err != nil {
				context.Logger.Error(err, "action_router.Service.msgHandler")
				return
			}

			time.Sleep(d)

		case model.ActionTypeMethod:
			msg, err := messages.NewCommand(act.Name, act.TargetType, act.TargetID, act.Args)
			if err != nil {
				context.Logger.Error(err, "action_router.Service.msgHandler")
				return
			}

			if err := msgs.I.Send(msg); err != nil {
				context.Logger.Error(err, "action_router.Service.msgHandler")
				return
			}

		case model.ActionTypeNotification:
			v, ok := act.Args["type"]
			if !ok {
				context.Logger.Error(errors.New("type not found"), "action_router.Service.msgHandler")
				return
			}

			notType, ok := v.(string)
			if !ok {
				context.Logger.Error(errors.New("type is not string"), "action_router.Service.msgHandler")
				return
			}

			v, ok = act.Args["text"]
			if !ok {
				context.Logger.Error(errors.New("text not found"), "action_router.Service.msgHandler")
				return
			}

			notText, ok := v.(string)
			if !ok {
				context.Logger.Error(errors.New("type is not string"), "action_router.Service.msgHandler")
				return
			}

			msg, err := events.NewNotification(notType, notText)
			if err != nil {
				context.Logger.Error(err, "action_router.Service.msgHandler")
				return
			}

			if err := msgs.I.Send(msg); err != nil {
				context.Logger.Error(err, "action_router.Service.msgHandler")
				return
			}

		default:
			context.Logger.Error(errors.Errorf("unknown action type %q", act.Type), "action_router.Service.msgHandler")
			return
		}
	}
}

func (o *Service) processItemEvents(msg interfaces.Message) {
	var state string

	switch msg.(type) {
	case item.OnChangeStateOn:
		state = "on"
	case item.OnChangeStateOff:
		state = "off"
	default:
		context.Logger.Warnf("unhandled item event [%s, %s, %s, %d]", msg.GetType(), msg.GetName(), msg.GetTargetType(), msg.GetTargetID())
		return
	}

	if err := store.I.Items().ChangeItem(msg.GetTargetID(), state); err != nil {
		context.Logger.Error(errors.Wrap(err, "Service.processItemEvents"))
		return
	}

	ws.I.Send(&model.ViewItem{ID: msg.GetTargetID(), Status: state})
}

func (o *Service) processObjectEvent(msg interfaces.Message) {
	// Ищем в таблице событие, которое пришло в топике
	items, err := store.I.Items().GetItemsForChange(msg.GetTargetType(), msg.GetTargetID(), msg.GetName())
	if err != nil {
		context.Logger.Error(errors.Wrap(err, "processObjectEvent"))
		return
	}

	// Перебираем найденные итемы, чтобы произвести с ними действие
	for _, item := range items {
		switch item.Type {

		case "button", "switch", "conditioner":
			item.Status, _ = msg.GetStringValue(item.EventValue)
			if err := store.I.Items().ChangeItem(item.ID, item.Status); err != nil {
				context.Logger.Error(errors.Wrap(err, "processObjectEvent"))
				return
			}

			ws.I.Send(item)

		case "sensor":
			item.Value, _ = msg.GetFloatValue(item.EventValue)

			sensor, err := store.I.Devices().GetSensor(item.ID)
			if err != nil {
				context.Logger.Error(errors.Wrap(err, "processObjectEvent"))
				return
			}

			if sensor.Current != item.Value {
				ws.I.Send(item)
			}

			// Обновление значения в таблице сенсоров
			if err := store.I.Devices().SetSensorValue(item.ID, item.Value); err != nil {
				context.Logger.Error(errors.Wrap(err, "processObjectEvent"))
				return
			}

			// Обновление значения в таблице графиков для сенсора
			t := time.Now().Format("2006-01-02T15:04")
			if err := store.I.History().SetHourlyValue(item.ID, t, item.Value); err != nil {
				context.Logger.Error(errors.Wrap(err, "processObjectEvent"))
				return
			}
		}
	}
}

func (o *Service) processNotification(msg interfaces.Message) {
	n, ok := msg.(*events.Notification)
	if !ok {
		context.Logger.Error(errors.Wrap(errors.Errorf("msg is not notification: %T", msg), "processNotification"))
		return
	}

	notification := &model.Notification{
		Type: n.Type,
		Text: n.Text,
		Date: time.Now().Format("2006-01-02T15:04:05"),
	}

	if err := store.I.Notifications().AddNotification(notification); err != nil {
		context.Logger.Error(errors.Wrap(err, "processNotification"))
		return
	}

	// Отправка сообщения через вебсокет
	ws.I.Send(notification)

	// Отправка критических сообщений в пуш уведомления
	if notification.Type == interfaces.NotificationTypeCritical {
		tokens, err := store.I.Notifications().GetPushTokens()
		if err != nil {
			context.Logger.Error(errors.Wrap(err, "processNotification"))
			return
		}

		msg := &model.PushNotification{
			Title:  "Важное уведомление!",
			Body:   notification.Text,
			Tokens: tokens,
		}

		data, err := json.Marshal(msg)
		if err != nil {
			context.Logger.Error(errors.Wrap(err, "processNotification"))
			return
		}

		if _, err := http.Post(context.Config["push_sender_address"]+"/push", "application/json", bytes.NewReader(data)); err != nil {
			context.Logger.Error(errors.Wrap(err, "processNotification"))
			return
		}
	}
}
