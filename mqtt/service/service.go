package service

import (
	"sync"
	"time"

	"github.com/VladimirDronik/touchon-server/events/service"
	"github.com/VladimirDronik/touchon-server/info"
	"github.com/VladimirDronik/touchon-server/mqtt/client"
	"github.com/VladimirDronik/touchon-server/mqtt/messages"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func New(client *client.Client, cfg map[string]string, bufferSize int, threads int, logger *logrus.Logger) (*Service, error) {
	o := &Service{
		client:     client,
		config:     cfg,
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
	config     map[string]string
	bufferSize int
	logger     *logrus.Logger
	topic      string
	threads    int
	wg         *sync.WaitGroup
	handler    func(messages.Message) error
}

func (o *Service) SetHandler(handler func(messages.Message) error) {
	o.handler = handler
}

func (o *Service) GetLogger() *logrus.Logger {
	return o.logger
}

func (o *Service) GetClient() *client.Client {
	return o.client
}

func (o *Service) GetConfig() map[string]string {
	return o.config
}

func (o *Service) Start() error {
	if o.handler == nil {
		return errors.Wrap(errors.New("handler is nil"), "Start")
	}

	msgs, err := o.client.Subscribe(o.topic, o.bufferSize)
	if err != nil {
		return errors.Wrap(err, "Start")
	}

	maxTravelTime := 0xFFFF * time.Hour
	if v := o.config["mqtt_max_travel_time"]; v != "" {
		maxTravelTime, err = time.ParseDuration(v)
		if err != nil {
			return errors.Wrap(err, "Start")
		}
	}

	o.wg.Add(o.threads)

	// Запускаем воркеров
	for i := 0; i < o.threads; i++ {
		go func() {
			defer o.wg.Done()

			for msg := range msgs {
				m, err := messages.NewFromMQTT(msg)
				if err != nil {
					o.logger.Error(err)
					continue
				}

				m.SetReceivedAt(time.Now())
				travelTime := m.GetReceivedAt().Sub(m.GetSentAt())

				o.GetLogger().Debugln()
				o.GetLogger().Debugf("MQTT: [%s] QoS=%d travelTime=%s", m.GetTopic(), m.GetQoS(), travelTime)
				o.GetLogger().Debug(m.String())
				if travelTime > maxTravelTime && o.GetLogger().Level >= logrus.DebugLevel {
					type TravelTimeTooLong struct {
						Duration string
						Message  messages.Message
					}
					msg := TravelTimeTooLong{
						Duration: travelTime.String(),
						Message:  m,
					}

					if err := o.client.SendRaw("debug/travel_time_too_long/"+info.Name, messages.QoSNotGuaranteed, false, msg); err != nil {
						o.GetLogger().Error(err)
					}
				}

				if m.GetTargetType() == messages.TargetTypeService && m.GetType() == messages.MessageTypeCommand && m.GetName() == "info" {
					m, err := service.NewOnInfoMessage("service/info")
					if err != nil {
						o.logger.Error(err)
						continue
					}

					if err := o.client.Send(m); err != nil {
						o.logger.Error(err)
						continue
					}

					continue
				}

				if err := o.handler(m); err != nil {
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
