package messages

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/VladimirDronik/touchon-server/info"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/pkg/errors"
)

func NewFromMQTT(msg mqtt.Message) (Message, error) {
	m := &MessageImpl{}

	if err := json.Unmarshal(msg.Payload(), m); err != nil {
		return nil, errors.Wrap(err, "NewFromMQTT")
	}
	m.SetTopic(msg.Topic())
	m.SetQoS(QoS(msg.Qos()))

	return m, nil
}

func NewMessage(msgType MessageType, name string, targetType TargetType, targetID int, payload map[string]interface{}) (Message, error) {
	switch {
	case msgType != MessageTypeCommand && msgType != MessageTypeEvent:
		return nil, errors.Wrap(errors.Errorf("unknown message type %q", msgType), "NewMessage")
	case name == "":
		return nil, errors.Wrap(errors.New("name is empty"), "NewMessage")
	case targetID < 0:
		return nil, errors.Wrap(errors.New("object ID < 0"), "NewMessage")
	}

	if payload == nil {
		payload = make(map[string]interface{})
	}

	return &MessageImpl{
		publisher:  info.Name,
		msgType:    msgType,
		name:       name,
		targetID:   targetID,
		targetType: targetType,
		payload:    payload,
	}, nil
}

type MessageImpl struct {
	retained   bool
	publisher  string
	topic      string
	msgType    MessageType // event,command
	name       string      // onChange,check
	targetID   int         // 82
	targetType TargetType
	payload    map[string]interface{} //
	qos        QoS
	sentAt     time.Time
	receivedAt time.Time
}

func (o *MessageImpl) GetRetained() bool {
	return o.retained
}

func (o *MessageImpl) GetPublisher() string {
	return o.publisher
}

func (o *MessageImpl) GetTopic() string {
	return o.topic
}

func (o *MessageImpl) GetType() MessageType {
	return o.msgType
}

func (o *MessageImpl) GetName() string {
	return o.name
}

func (o *MessageImpl) GetTargetID() int {
	return o.targetID
}

func (o *MessageImpl) GetTargetType() TargetType {
	return o.targetType
}

func (o *MessageImpl) GetPayload() map[string]interface{} {
	return o.payload
}

func (o *MessageImpl) GetQoS() QoS {
	return o.qos
}

func (o *MessageImpl) GetTopicPublisher() string {
	items := strings.Split(o.GetTopic(), "/")
	if len(items) > 0 {
		return items[0]
	}

	return ""
}

func (o *MessageImpl) GetTopicType() string {
	items := strings.Split(o.GetTopic(), "/")
	if len(items) > 1 {
		return items[1]
	}

	return ""
}

func (o *MessageImpl) GetTopicAction() string {
	items := strings.Split(o.GetTopic(), "/")
	if len(items) > 2 {
		return items[2]
	}

	return ""
}

func (o *MessageImpl) SetRetained(v bool) {
	o.retained = v
}

func (o *MessageImpl) SetPublisher(v string) {
	o.publisher = v
}

func (o *MessageImpl) SetTopic(v string) {
	o.topic = v
}

func (o *MessageImpl) SetType(v MessageType) {
	o.msgType = v
}

func (o *MessageImpl) SetName(v string) {
	o.name = v
}

func (o *MessageImpl) SetTargetID(v int) {
	o.targetID = v
}

func (o *MessageImpl) SetTargetType(v TargetType) {
	o.targetType = v
}

func (o *MessageImpl) SetPayload(v map[string]interface{}) {
	if o.payload == nil {
		o.payload = make(map[string]interface{}, len(v))
	}

	for k, v := range v {
		o.payload[k] = v
	}
}

func (o *MessageImpl) SetQoS(v QoS) {
	o.qos = v
}

func (o *MessageImpl) GetFloatValue(name string) (float32, error) {
	v, ok := o.GetPayload()[name]
	if !ok {
		return 0, errors.Wrap(errors.Errorf("%s not found", name), "GetFloatValue")
	}

	switch v := v.(type) {
	case float32:
		return v, nil
	case float64:
		return float32(v), nil
	case int:
		return float32(v), nil
	default:
		return 0, errors.Wrap(errors.Errorf("unexpected data type %T", v), "GetFloatValue")
	}
}

func (o *MessageImpl) GetStringValue(name string) (string, error) {
	v, ok := o.GetPayload()[name]
	if !ok {
		return "", errors.Wrap(errors.Errorf("%s not found", name), "GetStringValue")
	}

	switch v := v.(type) {
	case string:
		return v, nil
	default:
		return "", errors.Wrap(errors.Errorf("unexpected data type %T", v), "GetStringValue")
	}
}

func (o *MessageImpl) GetIntValue(name string) (int, error) {
	v, ok := o.GetPayload()[name]
	if !ok {
		return 0, errors.Wrap(errors.Errorf("%s not found", name), "GetIntValue")
	}

	switch v := v.(type) {
	case float32:
		return int(v), nil
	case float64:
		return int(v), nil
	case int:
		return v, nil
	default:
		return 0, errors.Wrap(errors.Errorf("unexpected data type %T", v), "GetIntValue")
	}
}

func (o *MessageImpl) GetBoolValue(name string) (bool, error) {
	v, ok := o.GetPayload()[name]
	if !ok {
		return false, errors.Wrap(errors.Errorf("%s not found", name), "GetBoolValue")
	}

	switch v := v.(type) {
	case bool:
		return v, nil
	default:
		return false, errors.Wrap(errors.Errorf("unexpected data type %T", v), "GetBoolValue")
	}
}

func (o *MessageImpl) GetSentAt() time.Time {
	return o.sentAt
}

func (o *MessageImpl) SetSentAt(v time.Time) {
	o.sentAt = v
}

func (o *MessageImpl) GetReceivedAt() time.Time {
	return o.receivedAt
}

func (o *MessageImpl) SetReceivedAt(v time.Time) {
	o.receivedAt = v
}

const atFormat = "02.01.2006 15:04:05.000000 MST"

func (o *MessageImpl) MarshalJSON() ([]byte, error) {
	m := &message{
		Publisher:  o.GetPublisher(),
		Type:       o.GetType(),
		Name:       o.GetName(),
		TargetID:   o.GetTargetID(),
		TargetType: o.GetTargetType(),
		Payload:    o.GetPayload(),
	}

	if !o.GetSentAt().IsZero() {
		m.SentAt = o.GetSentAt().Format(atFormat)
	}

	if !o.GetReceivedAt().IsZero() {
		m.ReceivedAt = o.GetReceivedAt().Format(atFormat)
	}

	if len(m.Payload) == 0 {
		m.Payload = nil
	}

	return json.Marshal(m)
}

func (o *MessageImpl) UnmarshalJSON(data []byte) error {
	m := &message{}

	if err := json.Unmarshal(data, m); err != nil {
		return errors.Wrap(err, "MessageImpl.UnmarshalJSON")
	}

	o.SetPublisher(m.Publisher)
	o.SetType(m.Type)
	o.SetName(m.Name)
	o.SetTargetID(m.TargetID)
	o.SetTargetType(m.TargetType)
	o.SetPayload(m.Payload)

	sentAt, err := time.Parse(atFormat, m.SentAt)
	if err == nil {
		o.SetSentAt(sentAt)
	}

	receivedAt, err := time.Parse(atFormat, m.ReceivedAt)
	if err == nil {
		o.SetReceivedAt(receivedAt)
	}

	return nil
}

func (o *MessageImpl) String() string {
	data, _ := o.MarshalJSON()
	return string(data)
}

type message struct {
	Publisher  string                 `json:"publisher"`
	Type       MessageType            `json:"type"`
	Name       string                 `json:"name"`
	TargetID   int                    `json:"target_id,omitempty"`
	TargetType TargetType             `json:"target_type,omitempty"`
	Payload    map[string]interface{} `json:"payload,omitempty"`
	SentAt     string                 `json:"sent_at"`
	ReceivedAt string                 `json:"received_at"`
}
