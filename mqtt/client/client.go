// Пакет для работы с MQTT. Реализует прием и отправку сообщений в шину.

package client

import (
	"encoding/json"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/VladimirDronik/touchon-server/info"
	"github.com/VladimirDronik/touchon-server/mqtt/messages"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func New(clientID, connString string, timeout time.Duration, tries int, logger *logrus.Logger) (*Client, error) {
	o := &Client{
		clientID: clientID,
		timeout:  timeout,
		tries:    tries,
		chans:    make(map[string][]chan mqtt.Message),
		logger:   logger,
	}

	var err error
	o.connString, err = url.Parse(connString)
	if err != nil {
		return nil, errors.Wrap(err, "Client.New")
	}

	password, _ := o.connString.User.Password()
	opts := mqtt.NewClientOptions().
		AddBroker("tcp://" + o.connString.Host).
		SetUsername(o.connString.User.Username()).
		SetPassword(password).
		SetClientID(clientID)

	o.client = mqtt.NewClient(opts)

	token := o.client.Connect()
	if err := o.processToken(token); err != nil {
		return nil, errors.Wrap(err, "New")
	}

	return o, nil
}

type Client struct {
	clientID   string
	client     mqtt.Client
	timeout    time.Duration
	tries      int
	connString *url.URL
	logger     *logrus.Logger

	mu    sync.Mutex
	chans map[string][]chan mqtt.Message
}

func (o *Client) pushChan(topic string, ch chan mqtt.Message) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.chans[topic] = append(o.chans[topic], ch)
}

func (o *Client) popChans(topic string) []chan mqtt.Message {
	o.mu.Lock()
	defer o.mu.Unlock()
	chans := o.chans[topic]
	delete(o.chans, topic)
	return chans
}

// processToken Синхронно ожидает результата
func (o *Client) processToken(token mqtt.Token) error {
	var ok bool

	for i := 0; i < o.tries; i++ {
		if ok = token.WaitTimeout(o.timeout); ok {
			break
		}
	}

	if !ok {
		return errors.Wrap(errors.New("timeout"), "processToken")
	}

	if err := token.Error(); err != nil {
		return errors.Wrap(err, "processToken")
	}

	return nil
}

// Subscribe Подписывает на топики
func (o *Client) Subscribe(topic string, bufferSize int) (<-chan mqtt.Message, error) {
	c := make(chan mqtt.Message, bufferSize)

	token := o.client.Subscribe(topic, 0, func(client mqtt.Client, msg mqtt.Message) {
		// Свои сообщения игнорируем
		if !strings.HasPrefix(msg.Topic(), info.Name) {
			o.logger.Debugf("mqtt.Client.Receive: [%s]", msg.Topic())
			o.logger.Tracef("mqtt.Client.Receive: %s", string(msg.Payload()))

			c <- msg
		}
	})

	if err := o.processToken(token); err != nil {
		return nil, errors.Wrap(err, "Subscribe")
	}

	// Сохраняем канал, чтобы можно было его закрыть
	o.pushChan(topic, c)

	return c, nil
}

// Unsubscribe Отменяет подписку на топики
func (o *Client) Unsubscribe(topics ...string) error {
	token := o.client.Unsubscribe(topics...)

	if err := o.processToken(token); err != nil {
		return errors.Wrap(err, "Unsubscribe")
	}

	// Закрываем каналы, чтобы обработчики сообщений могли завершиться
	for _, topic := range topics {
		for _, ch := range o.popChans(topic) {
			close(ch)
		}
	}

	return nil
}

// Send Отправляет сообщения в топик
// sync - to track delivery of the message to the broker
func (o *Client) Send(msg messages.Message, sync ...bool) error {
	msg.SetSentAt(time.Now())

	if err := o.SendRaw(msg.GetTopic(), msg.GetQoS(), msg.GetRetained(), msg, sync...); err != nil {
		return errors.Wrap(err, "Send")
	}

	return nil
}

// SendRaw Отправляет сообщения в топик
// sync - to track delivery of the message to the broker
func (o *Client) SendRaw(topic string, qos messages.QoS, retained bool, payload interface{}, sync ...bool) error {
	if topic == "" {
		return errors.Wrap(errors.New("topic is empty"), "SendRaw")
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return errors.Wrap(err, "SendRaw")
	}

	o.logger.Debugf("mqtt.Client.Send: [%s]", topic)
	o.logger.Tracef("mqtt.Client.Send: [%s] %s", topic, string(data))

	token := o.client.Publish(topic, byte(qos), retained, payload)
	if len(sync) > 0 && sync[0] {
		if err := o.processToken(token); err != nil {
			return errors.Wrap(err, "SendRaw")
		}
	}

	return nil
}

func (o *Client) GetTopicFromConnectionString() string {
	topic := "#"
	if len(o.connString.Path) > 1 {
		topic = o.connString.Path[1:]
	}
	return topic
}

func (o *Client) Shutdown() error {
	// Получаем список топиков
	o.mu.Lock()
	topics := make([]string, 0, len(o.chans))
	for topic := range o.chans {
		topics = append(topics, topic)
	}
	o.mu.Unlock()

	errs := make([]error, 0, 5)

	// Отписываемся от всех топиков
	if err := o.Unsubscribe(topics...); err != nil {
		errs = append(errs, err)
	}

	// Отключаемся от шины
	o.client.Disconnect(uint(o.timeout.Milliseconds()))

	if len(errs) > 0 {
		return errors.Wrap(errs[0], "mqttClient.Shutdown")
	}

	return nil
}
