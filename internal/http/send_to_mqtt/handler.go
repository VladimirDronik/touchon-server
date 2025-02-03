package send_to_mqtt

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/valyala/fasthttp"
	"touchon-server/lib/helpers"
	mqttClient "touchon-server/lib/mqtt/client"
	"touchon-server/lib/mqtt/messages"
)

// Отправить сообщение в MQTT
// @Summary Отправить сообщение в MQTT
// @Tags Service
// @Description Отправить сообщение в MQTT
// @ID ServiceExecMethod
// @Accept text/json
// @Param publisher query string false "Издатель" Default(swagger)
// @Param type query messages.MessageType false "Тип сообщения" Default(command)
// @Param name query string false "Название сообщения/метода"
// @Param target_type query messages.TargetType false "Тип сущности" Default(object)
// @Param target_id query int false "ID сущности" Default(0)
// @Param payload body map[string]interface{} false "Параметры сообщения" Default({})
// @Param retained query bool false "Retained" Default(false)
// @Param topic query string false "Topic" Default(swagger)
// @Param qos query int false "QoS" Enums(0,1,2) Default(0)
// @Produce json
// @Success      200 {object} http.Response[any]
// @Failure      500 {object} http.Response[any]
// @Router /_/mqtt [post]
func Handler(ctx *fasthttp.RequestCtx) (interface{}, int, error) {
	publisher := helpers.GetParam(ctx, "publisher")
	msgType := helpers.GetParam(ctx, "type")
	name := helpers.GetParam(ctx, "name")
	targetType := helpers.GetParam(ctx, "target_type")
	targetID, _ := helpers.GetUintParam(ctx, "target_id")
	payload := ctx.Request.Body()
	retained, _ := helpers.GetBoolParam(ctx, "retained")
	topic := helpers.GetParam(ctx, "topic")
	qos, _ := helpers.GetUintParam(ctx, "qos")

	var pl map[string]interface{}
	if len(payload) > 0 {
		pl = map[string]interface{}{}
		if err := json.Unmarshal(payload, &pl); err != nil {
			return nil, http.StatusBadRequest, err
		}
	}

	msg := &message{
		Publisher:  publisher,
		Type:       msgType,
		Name:       name,
		TargetID:   targetID,
		TargetType: targetType,
		Payload:    pl,
		SentAt:     time.Now().Format(messages.TimeLabelFormat),
	}

	if err := mqttClient.I.SendRaw(topic, qos, retained, msg); err != nil {
		return nil, http.StatusBadRequest, err
	}

	return nil, http.StatusOK, nil
}

type message struct {
	Publisher  string                 `json:"publisher"`
	Type       messages.MessageType   `json:"type"`
	Name       string                 `json:"name"`
	TargetID   int                    `json:"target_id,omitempty"`
	TargetType messages.TargetType    `json:"target_type,omitempty"`
	Payload    map[string]interface{} `json:"payload,omitempty"`
	SentAt     string                 `json:"sent_at"`
	ReceivedAt string                 `json:"received_at"`
}
