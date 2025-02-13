package action_router

import (
	"time"

	"github.com/pkg/errors"
	"touchon-server/internal/context"
	"touchon-server/internal/model"
	"touchon-server/internal/store"
	"touchon-server/lib/events"
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
	handlerID int
}

func (o *Service) Start() error {
	var err error

	o.handlerID, err = msgs.I.Subscribe(interfaces.MessageTypeEvent, "", "", nil, o.msgHandler)
	if err != nil {
		return errors.Wrap(err, "scripts.Start")
	}

	return nil
}

func (o *Service) Shutdown() error {
	msgs.I.Unsubscribe(o.handlerID)

	return nil
}

func (o *Service) msgHandler(msg interfaces.Message) {
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
