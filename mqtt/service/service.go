package service

import (
	"sync"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"object-manager/internal/touchon-server/mqtt/client"
	"object-manager/internal/touchon-server/mqtt/messages"
)

func New(client *client.Client, bufferSize int, threads int, logger *logrus.Logger) (*Service, error) {
	o := &Service{
		client:     client,
		bufferSize: bufferSize,
		topic:      client.GetTopicFromConnectionString(),
		logger:     logger,
		threads:    threads,
		wg:         &sync.WaitGroup{},
	}

	return o, nil
}

type Service struct {
	client     *client.Client
	bufferSize int
	logger     *logrus.Logger
	topic      string
	threads    int
	wg         *sync.WaitGroup
	handler    func(mqtt.Message) error
}

func (o *Service) SetHandler(handler func(mqtt.Message) error) {
	o.handler = handler
}

func (o *Service) GetLogger() *logrus.Logger {
	return o.logger
}

func (o *Service) GetClient() *client.Client {
	return o.client
}

func (o *Service) Start() error {
	if o.handler == nil {
		return errors.Wrap(errors.New("handler is nil"), "Start")
	}

	msgs, err := o.client.Subscribe(o.topic, o.bufferSize)
	if err != nil {
		return errors.Wrap(err, "Start")
	}

	o.wg.Add(o.threads)

	// Запускаем воркеров
	for i := 0; i < o.threads; i++ {
		go func() {
			defer o.wg.Done()

			for msg := range msgs {
				if err := o.handler(msg); err != nil {
					o.logger.Error(err)
				}
			}
		}()
	}

	o.logger.Info("MQTT: сервис запущен")

	return nil
}

func (o *Service) Shutdown() error {
	o.logger.Info("MQTT: Останавливаем сервис")

	if err := o.client.Shutdown(); err != nil {
		return errors.Wrap(err, "mqttService.Shutdown")
	}

	o.logger.Info("MQTT: Ждем остановки всех воркеров")

	o.wg.Wait()

	return nil
}

// Send отправка сообщения в топик
func (o *Service) Send(msg messages.Message) error {
	if err := o.client.Send(msg); err != nil {
		return errors.Wrap(err, "Send")
	}

	return nil
}
