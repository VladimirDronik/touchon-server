package messages

import (
	"sync"
	"time"

	"github.com/pkg/errors"
	"touchon-server/internal/g"
	"touchon-server/lib/interfaces"
	"touchon-server/lib/subscribers"
)

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

	wg := sync.WaitGroup{}
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
		g.Logger.Debugf("messages.ServiceImpl: %T [%s, %s, %s, %d, %v] - %d handlers", msg, msg.GetType(), msg.GetName(), msg.GetTargetType(), msg.GetTargetID(), msg.GetPayload(), len(handlers))

		if len(handlers) == 0 {
			g.Logger.Warnf("Unhandled msg: [%s, %s, %s, %d]", msg.GetType(), msg.GetName(), msg.GetTargetType(), msg.GetTargetID())
			continue
		}

		wg.Add(len(handlers))

		// Одновременно запускаем все обработчики
		for _, handler := range handlers {
			go func(handler interfaces.MsgHandler) {
				defer wg.Done()
				handler(o, msg)
			}(handler)
		}

		// Ждем завершения всех обработчиков
		wg.Wait()
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

	select {
	case <-o.done:
		return errors.Wrap(errors.New("service is done"), "messages.ServiceImpl.Send")
	default:
	}

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
