package messages

import (
	"sync"
	"time"

	"github.com/pkg/errors"
	"touchon-server/internal/context"
	"touchon-server/lib/interfaces"
	"touchon-server/lib/subscribers"
)

// Global instance
var I interfaces.MessagesService

func NewService(threads, queueCap int) (interfaces.MessagesService, error) {
	switch {
	case threads < 1:
		return nil, errors.Wrap(errors.New("threads < 1"), "messages.NewService")
	case queueCap < 1:
		return nil, errors.Wrap(errors.New("queueCap < 1"), "messages.NewService")
	}

	return &ServiceImpl{
		mu:          sync.Mutex{},
		threads:     threads,
		msgs:        make(chan interfaces.Message, queueCap),
		subscribers: subscribers.New(2000),
		wg:          sync.WaitGroup{},
		done:        make(chan struct{}),
	}, nil
}

type ServiceImpl struct {
	threads     int                      // Кол-во потоков обработки сообщений
	subscribers *subscribers.Subscribers // Механизм получения обработчиков сообщений по телу сообщения
	wg          sync.WaitGroup           // Контролирует завершение всех потоков
	done        chan struct{}            // Признак остановки сервиса

	mu   sync.Mutex              // Контролирует запись и закрытие канала
	msgs chan interfaces.Message // Канал сообщений
}

func (o *ServiceImpl) Start() error {
	o.wg.Add(o.threads)

	for i := 0; i < o.threads; i++ {
		go o.worker()
	}

	return nil
}

func (o *ServiceImpl) Shutdown() error {
	close(o.done)

	// Блокируем запись в закрытый канал msgs
	o.mu.Lock()
	defer o.mu.Unlock()

	close(o.msgs)
	o.wg.Wait()

	return nil
}

func (o *ServiceImpl) worker() {
	defer o.wg.Done()

	var msg interfaces.Message
	var ok bool

	for {
		// Ждем либо сообщение в очереди, либо завершение сервиса
		select {
		case <-o.done: // Если сервис останавливается, бросаем все и уходим
			return

		case msg, ok = <-o.msgs:
			if !ok { // false получим только в случае закрытого и пустого канала
				return
			}
		}

		handlers := o.subscribers.GetHandlers(msg)
		if len(handlers) == 0 {
			context.Logger.Warnf("Unhandled msg: [%s, %s, %s, %d]", msg.GetType(), msg.GetName(), msg.GetTargetType(), msg.GetTargetID())
			continue
		}

		for _, handler := range handlers {
			handler(msg)
		}
	}
}

func (o *ServiceImpl) Subscribe(msgType interfaces.MessageType, name string, targetType interfaces.TargetType, targetID *int, handler interfaces.MsgHandler) (int, error) {
	handlerID, err := o.subscribers.AddHandler(msgType, name, targetType, targetID, handler)
	if err != nil {
		return 0, errors.Wrap(err, "messages.ServiceImpl.Subscribe")
	}

	return handlerID, nil
}

func (o *ServiceImpl) Unsubscribe(handlerIDs ...int) {
	for _, handlerID := range handlerIDs {
		o.subscribers.DeleteHandler(handlerID)
	}
}

func (o *ServiceImpl) Send(msgs ...interfaces.Message) error {
	// Блокируем закрытие канала msgs
	o.mu.Lock()
	defer o.mu.Unlock()

	t := time.NewTimer(time.Second)
	defer t.Stop()

	for _, msg := range msgs {
		select {
		case <-o.done:
			return errors.Wrap(errors.New("shutting done started"), "messages.ServiceImpl.Send")
		case o.msgs <- msg:
		case <-t.C:
			return errors.Wrap(errors.New("queue is full"), "messages.ServiceImpl.Send")
		}
	}

	return nil
}

// TR

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
//	}
//
//	return nil
//}
