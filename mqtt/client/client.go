// Пакет для работы с MQTT. Реализует прием и отправку сообщений в шину.

package client

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/VladimirDronik/touchon-server/info"
	"github.com/VladimirDronik/touchon-server/mqtt/messages"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// Global instance
var I *Client

func New(clientID, connString string, timeout time.Duration, tries int, logger *logrus.Logger) (*Client, error) {
	o := &Client{
		clientID:       clientID,
		timeout:        timeout,
		tries:          tries,
		chans:          make(map[string][]chan mqtt.Message),
		logger:         logger,
		ignoreSelfMsgs: true,
	}

	var err error
	o.connString, err = url.Parse(connString)
	if err != nil {
		return nil, errors.Wrap(err, "Client.New")
	}

	password, _ := o.connString.User.Password()
	opts := mqtt.NewClientOptions().
		AddBroker("mqtt://" + o.connString.Host).
		SetUsername(o.connString.User.Username()).
		SetPassword(password).
		SetClientID(clientID).
		SetResumeSubs(true)

	opts.OnConnectionLost = func(client mqtt.Client, reason error) {
		o.logger.Debugf("mqtt.Client: connection lost: %v", reason)
	}

	opts.OnConnectAttempt = func(broker *url.URL, tlsCfg *tls.Config) *tls.Config {
		o.logger.Debugf("mqtt.Client: connection attempt: %v", broker)
		return tlsCfg
	}

	opts.OnConnect = func(client mqtt.Client) {
		o.logger.Debugf("mqtt.Client: connected = %t", client.IsConnected())
	}

	o.client = mqtt.NewClient(opts)

	token := o.client.Connect()
	if err := o.processToken(token); err != nil {
		return nil, errors.Wrap(err, "New")
	}

	return o, nil
}

type Client struct {
	clientID       string
	client         mqtt.Client
	timeout        time.Duration
	tries          int
	connString     *url.URL
	logger         *logrus.Logger
	ignoreSelfMsgs bool

	mu    sync.Mutex
	chans map[string][]chan mqtt.Message
}

func (o *Client) GetIgnoreSelfMsgs() bool {
	return o.ignoreSelfMsgs
}

func (o *Client) SetIgnoreSelfMsgs(v bool) {
	o.ignoreSelfMsgs = v
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
		if ok = token.WaitTimeout(o.timeout); !ok {
			continue
		}

		if err := token.Error(); err != nil {
			return errors.Wrap(err, "processToken")
		}

		break
	}

	if !ok {
		return errors.Wrap(errors.New("timeout"), "processToken")
	}

	return nil
}

// Subscribe Подписывает на топики
func (o *Client) Subscribe(topic string, bufferSize int) (<-chan mqtt.Message, error) {
	c := make(chan mqtt.Message, bufferSize)

	token := o.client.Subscribe(topic, 0, func(client mqtt.Client, msg mqtt.Message) {
		// Свои сообщения игнорируем
		if !o.ignoreSelfMsgs || !strings.HasPrefix(msg.Topic(), info.Name) {
			msgInfo := getMetaInfoFromRawMsg(msg.Payload())

			switch o.logger.Level {
			case logrus.DebugLevel:
				o.logger.Debugf("mqtt.Client.Receive: [%s]%s", msg.Topic(), msgInfo)
			case logrus.TraceLevel:
				o.logger.Tracef("mqtt.Client.Receive: [%s]%s %s", msg.Topic(), msgInfo, string(msg.Payload()))
			}

			// for non blocking
			go func() { c <- msg }()
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
func (o *Client) Send(msg messages.Message) error {
	msg.SetSentAt(time.Now())

	if err := o.SendRaw(msg.GetTopic(), msg.GetQoS(), msg.GetRetained(), msg); err != nil {
		return errors.Wrap(err, "Send")
	}

	return nil
}

var patternTargetType = regexp.MustCompile(`"target_type"\s*:\s*"([^"]+)"`)
var patternTargetID = regexp.MustCompile(`"target_id"\s*:\s*(\d+)`)
var patternName = regexp.MustCompile(`"name"\s*:\s*"([^"]+)"`)

func getMetaInfoFromRawMsg(data []byte) string {
	var targetType string
	if r := patternTargetType.FindStringSubmatch(string(data)); len(r) == 2 {
		targetType = r[1]
	}

	var targetID string
	if r := patternTargetID.FindStringSubmatch(string(data)); len(r) == 2 {
		targetID = r[1]
	}

	var name string
	if r := patternName.FindStringSubmatch(string(data)); len(r) == 2 {
		name = r[1]
	}

	return fmt.Sprintf(" [%s/%s/%s]", targetType, targetID, name)
}

// SendRaw Отправляет сообщения в топик
// sync - to track delivery of the message to the broker
func (o *Client) SendRaw(topic string, qos messages.QoS, retained bool, payload interface{}) error {
	if topic == "" {
		return errors.Wrap(errors.New("topic is empty"), "SendRaw")
	}

	switch v := payload.(type) {
	case []byte:
	default:
		var err error
		payload, err = json.Marshal(v)
		if err != nil {
			return errors.Wrap(err, "SendRaw")
		}
	}

	msgInfo := getMetaInfoFromRawMsg(payload.([]byte))

	switch o.logger.Level {
	case logrus.DebugLevel:
		o.logger.Debugf("mqtt.Client.Send: [%s]%s", topic, msgInfo)
	case logrus.TraceLevel:
		o.logger.Tracef("mqtt.Client.Send: [%s]%s %s", topic, msgInfo, string(payload.([]byte)))
	}

	token := o.client.Publish(topic, byte(qos), retained, payload)
	if err := o.processToken(token); err != nil {
		return errors.Wrap(err, "SendRaw")
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
