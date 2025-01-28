package service

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/VladimirDronik/touchon-server/events/service"
	"github.com/VladimirDronik/touchon-server/info"
	mqtt "github.com/VladimirDronik/touchon-server/mqtt/client"
	"github.com/VladimirDronik/touchon-server/mqtt/messages"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func New(client mqtt.Client, cfg map[string]string, bufferSize int, threads int, logger *logrus.Logger) (*Service, error) {
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
	client     mqtt.Client
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

func (o *Service) GetClient() mqtt.Client {
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

				travelTime := o.processTravelTime(m, maxTravelTime)

				switch o.logger.Level {
				case logrus.DebugLevel:
					o.logger.Debugf("mqtt.Service.Receive: [%s] QoS=%d travelTime=%s", m.GetTopic(), m.GetQoS(), travelTime)
				case logrus.TraceLevel:
					o.logger.Tracef("mqtt.Service.Receive: [%s] QoS=%d travelTime=%s %s", m.GetTopic(), m.GetQoS(), travelTime, m.String())
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

func (o *Service) processTravelTime(m messages.Message, maxTravelTime time.Duration) string {
	if m.GetSentAt().IsZero() {
		return ""
	}

	travelTime := m.GetReceivedAt().Sub(m.GetSentAt())

	if travelTime > maxTravelTime && o.GetLogger().Level >= logrus.DebugLevel {
		type TravelTimeTooLong struct {
			Duration string
			Message  messages.Message
		}
		msg := TravelTimeTooLong{
			Duration: travelTime.String(),
			Message:  m,
		}

		data, _ := json.Marshal(msg)

		if err := o.client.SendRaw("debug/travel_time_too_long/"+info.Name, messages.QoSNotGuaranteed, false, data); err != nil {
			o.GetLogger().Error(err)
		}
	}

	return travelTime.String()
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
