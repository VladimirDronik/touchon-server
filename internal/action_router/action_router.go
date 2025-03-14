package action_router

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"touchon-server/internal/g"
	"touchon-server/internal/model"
	"touchon-server/internal/store"
	"touchon-server/lib/events"
	"touchon-server/lib/events/item"
	"touchon-server/lib/interfaces"
	"touchon-server/lib/messages"
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
	handlerID, err := g.Msgs.Subscribe(interfaces.MessageTypeEvent, "", "", nil, o.actionRouter)
	if err != nil {
		return errors.Wrap(err, "action_router.Service.Start")
	}
	o.handlerIDs = append(o.handlerIDs, handlerID)

	handlerID, err = g.Msgs.Subscribe(interfaces.MessageTypeEvent, "", interfaces.TargetTypeObject, nil, o.processObjectEvent)
	if err != nil {
		return errors.Wrap(err, "action_router.Service.Start")
	}
	o.handlerIDs = append(o.handlerIDs, handlerID)

	handlerID, err = g.Msgs.Subscribe(interfaces.MessageTypeEvent, "", interfaces.TargetTypeItem, nil, o.processItemEvents)
	if err != nil {
		return errors.Wrap(err, "action_router.Service.Start")
	}
	o.handlerIDs = append(o.handlerIDs, handlerID)

	handlerID, err = g.Msgs.Subscribe(interfaces.MessageTypeEvent, "on_notify", "", nil, o.processNotification)
	if err != nil {
		return errors.Wrap(err, "action_router.Service.Start")
	}
	o.handlerIDs = append(o.handlerIDs, handlerID)

	return nil
}

func (o *Service) Shutdown() error {
	g.Msgs.Unsubscribe(o.handlerIDs...)

	return nil
}

func (o *Service) actionRouter(svc interfaces.MessageSender, msg interfaces.Message) {
	ev, ok := msg.(interfaces.Event)
	if !ok {
		g.Logger.Error(errors.Wrap(errors.Errorf("msg is not event: %T", msg), "action_router.Service.msgHandler"))
		return
	}

	storeEvent, err := store.I.EventsRepo().GetEvent(ev.GetTargetType(), ev.GetTargetID(), ev.GetName())
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return
		}

		g.Logger.Error(errors.Wrap(err, "action_router.Service.msgHandler"))
		return
	}

	actions, err := store.I.EventActionsRepo().GetActions(storeEvent.ID)
	if err != nil {
		g.Logger.Error(errors.Wrap(err, "action_router.Service.msgHandler"))
		return
	}

	for _, act := range actions[storeEvent.ID] {
		if !act.Enabled {
			continue
		}

		switch act.Type {
		case interfaces.ActionTypeDelay:
			v, ok := act.Args["duration"]
			if !ok {
				g.Logger.Error(errors.Wrap(errors.New("duration not found"), "action_router.Service.msgHandler"))
				return
			}

			s, ok := v.(string)
			if !ok {
				g.Logger.Error(errors.Wrap(errors.New("duration is not string"), "action_router.Service.msgHandler"))
				return
			}

			d, err := time.ParseDuration(s)
			if err != nil {
				g.Logger.Error(err, "action_router.Service.msgHandler")
				return
			}

			time.Sleep(d)

		case interfaces.ActionTypeMethod:
			msg, err := messages.NewCommand(act.Name, act.TargetType, act.TargetID, act.Args)
			if err != nil {
				g.Logger.Error(err, "action_router.Service.msgHandler")
				return
			}

			if err := svc.Send(msg); err != nil {
				g.Logger.Error(err, "action_router.Service.msgHandler")
				return
			}

		case interfaces.ActionTypeNotification:
			v, ok := act.Args["type"]
			if !ok {
				g.Logger.Error(errors.New("type not found"), "action_router.Service.msgHandler")
				return
			}

			notType, ok := v.(string)
			if !ok {
				g.Logger.Error(errors.New("type is not string"), "action_router.Service.msgHandler")
				return
			}

			v, ok = act.Args["text"]
			if !ok {
				g.Logger.Error(errors.New("text not found"), "action_router.Service.msgHandler")
				return
			}

			notText, ok := v.(string)
			if !ok {
				g.Logger.Error(errors.New("type is not string"), "action_router.Service.msgHandler")
				return
			}

			msg, err := events.NewNotification(notType, notText)
			if err != nil {
				g.Logger.Error(err, "action_router.Service.msgHandler")
				return
			}

			if err := svc.Send(msg); err != nil {
				g.Logger.Error(err, "action_router.Service.msgHandler")
				return
			}

		default:
			g.Logger.Error(errors.Errorf("unknown action type %q", act.Type), "action_router.Service.msgHandler")
			return
		}
	}
}

func (o *Service) processItemEvents(svc interfaces.MessageSender, msg interfaces.Message) {
	var state string

	switch msg.(type) {
	case item.OnChangeStateOn:
		state = "on"
	case item.OnChangeStateOff:
		state = "off"
	default:
		g.Logger.Warnf("unhandled item event %T [%s, %s, %s, %d]", msg, msg.GetType(), msg.GetName(), msg.GetTargetType(), msg.GetTargetID())
		return
	}

	if err := store.I.Items().ChangeItem(msg.GetTargetID(), state); err != nil {
		g.Logger.Error(errors.Wrap(err, "Service.processItemEvents"))
		return
	}

	g.WSServer.Send("item_status", &model.ItemForWS{ID: msg.GetTargetID(), Status: state})
}

func (o *Service) processObjectEvent(svc interfaces.MessageSender, msg interfaces.Message) {
	// Ищем в таблице событие, которое пришло в топике
	items, err := store.I.Items().GetItemsForChange(msg.GetTargetType(), msg.GetTargetID(), msg.GetName())
	if err != nil {
		g.Logger.Error(errors.Wrap(err, "processObjectEvent"))
		return
	}

	// Перебираем найденные итемы, чтобы произвести с ними действие
	for _, item := range items {

		args := make(map[string]interface{})
		err = json.Unmarshal(item.EventArgs, &args)
		if err != nil {
			g.Logger.Error(errors.Wrap(err, "processObjectEvent: error retrieving event arguments"))
			return
		}

		switch item.Type {

		case "button", "switch", "conditioner":
			item.Status, _ = msg.GetStringValue(args["param"].(string))
			if err := store.I.Items().ChangeItem(item.ID, item.Status); err != nil {
				g.Logger.Error(errors.Wrap(err, "processObjectEvent"))
				return
			}

			g.WSServer.Send("item_status", &model.ItemForWS{ID: item.ID, Status: item.Status})

		case "sensor":
			valueSensor, _ := msg.GetFloatValue(args["param"].(string))

			sensor, err := store.I.Devices().GetSensor(item.ID)
			if err != nil {
				g.Logger.Error(errors.Wrap(err, "processObjectEvent"))
				return
			}

			if sensor.Current != valueSensor {
				g.WSServer.Send("item_value", &model.ItemForWS{
					ID:     item.ID,
					Values: append([]model.Value{}, model.Value{Type: sensor.Type, Value: valueSensor}),
				})
			}

			// Обновление значения в таблице графиков для сенсора
			t := time.Now().Format("2006-01-02T15:04")
			if err := store.I.History().SetValue(item.ID, t, valueSensor, model.TableDailyHistory); err != nil {
				g.Logger.Error(errors.Wrap(err, "processObjectEvent"))
				return
			}
		}
	}
}

func (o *Service) processNotification(svc interfaces.MessageSender, msg interfaces.Message) {
	n, ok := msg.(*events.Notification)
	if !ok {
		g.Logger.Error(errors.Wrap(errors.Errorf("msg is not notification: %T", msg), "processNotification"))
		return
	}

	nType, err := n.GetStringValue("type")
	if err != nil {
		g.Logger.Error(errors.Wrap(err, "processNotification"))
		return
	}

	nText, err := n.GetStringValue("text")
	if err != nil {
		g.Logger.Error(errors.Wrap(err, "processNotification"))
		return
	}

	notification := &model.Notification{
		Type: nType,
		Text: nText,
		Date: time.Now().Format("2006-01-02T15:04:05"),
	}

	if err := store.I.Notifications().AddNotification(notification); err != nil {
		g.Logger.Error(errors.Wrap(err, "processNotification"))
		return
	}

	// Отправка сообщения через вебсокет
	g.WSServer.Send("send_notification", notification)

	// Отправка критических сообщений в пуш уведомления
	if notification.Type == interfaces.NotificationTypeCritical {
		tokens, err := store.I.Notifications().GetPushTokens()
		if err != nil {
			g.Logger.Error(errors.Wrap(err, "processNotification"))
			return
		}

		msg := &model.PushNotification{
			Title:  "Важное уведомление!",
			Body:   notification.Text,
			Tokens: tokens,
		}

		data, err := json.Marshal(msg)
		if err != nil {
			g.Logger.Error(errors.Wrap(err, "processNotification"))
			return
		}

		if _, err := http.Post(g.Config["push_sender_address"]+"/push", "application/json", bytes.NewReader(data)); err != nil {
			g.Logger.Error(errors.Wrap(err, "processNotification"))
			return
		}
	}
}
