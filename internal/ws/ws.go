package ws

import (
	"encoding/json"
	"github.com/google/uuid"
	"net/http"
	"time"

	"github.com/fasthttp/websocket"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
	"touchon-server/internal/g"
	"touchon-server/internal/helpers"
	"touchon-server/internal/token"
	"touchon-server/lib/http/server"
)

func New() (*Server, error) {
	baseServer, err := server.New("WS", g.Config, g.Logger)
	if err != nil {
		return nil, errors.Wrap(err, "http.New")
	}

	o := &Server{
		Server:  baseServer,
		clients: make(map[int]map[*websocket.Conn]bool),
	}

	o.GetServer().Handler = o.handler

	return o, nil
}

type Server struct {
	*server.Server
	clients map[int]map[*websocket.Conn]bool
}

//func (o *Server) Start() error {
//	if err := o.Server.Start(o.GetConfig()["ws_addr"]); err != nil {
//		return errors.Wrap(err, "Start")
//	}
//
//	return nil
//}

//func (o *Server) Shutdown() error {
//	if err := o.Server.Shutdown(); err != nil {
//		return errors.Wrap(err, "Shutdown")
//	}
//
//	return nil
//}

func (o *Server) handler(ctx *fasthttp.RequestCtx) {
	o.SetRequestID(ctx)
	helpers.DumpRequest(o.GetLogger(), ctx)
	defer helpers.DumpResponse(o.GetLogger(), ctx)

	deviceID := g.DisabledAuthDeviceID
	var err error

	if o.GetConfig()["token_secret"] != "disable_auth" {
		authToken := string(ctx.Request.Header.Peek("token"))
		if authToken == "" {
			ctx.Error(http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		// Проверяем и извлекаем ID клиента из токена
		if deviceID, err = token.KeysExtract(authToken, o.GetConfig()["token_secret"]); err != nil {
			ctx.Error(http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
	}

	var upgrader = websocket.FastHTTPUpgrader{
		CheckOrigin: func(ctx *fasthttp.RequestCtx) bool { return true }, // Пропускаем любой запрос
	}

	err = upgrader.Upgrade(ctx, func(ws *websocket.Conn) {
		defer ws.Close()

		ws.SetPongHandler(func(appData string) error { return nil })
		ws.SetPingHandler(func(appData string) error {
			return ws.WriteControl(websocket.PongMessage, []byte(appData), time.Now().Add(time.Second))
		})

		// Отправляем Ping сообщения каждые 30 секунд
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		go func() {
			for range ticker.C {
				if err := ws.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(time.Second)); err != nil {
					o.GetLogger().Error("Failed to send Ping: ", err)
					return
				}
			}
		}()

		if o.clients[deviceID] == nil {
			o.clients[deviceID] = map[*websocket.Conn]bool{}
		}

		o.GetLogger().Debugf("ws.Server: client %d connected", deviceID)

		o.clients[deviceID][ws] = true        // Сохраняем соединение, используя его как ключ
		defer delete(o.clients[deviceID], ws) // Удаляем соединение

		for {
			mt, message, err := ws.ReadMessage()
			if err != nil || mt == websocket.CloseMessage {
				o.GetLogger().Debugf("ws.Server: client %d disconnected", deviceID)
				break // Выходим из цикла, если клиент пытается закрыть соединение или связь прервана
			}

			o.GetLogger().Debugf("ws.Server.Receive: %s", string(message))

			//go func(message []byte) {
			//	// TODO
			//}(message)
		}
	})

	if err != nil {
		o.GetLogger().Error(err)
	}
}

func (o *Server) Send(event string, data interface{}) {
	msg := struct {
		Event     string      `json:"event"`
		UUID      string      `json:"uuid"`
		TimeStamp int64       `json:"timestamp"`
		Data      interface{} `json:"data"`
	}{
		Event:     event,
		UUID:      uuid.New().String(),
		TimeStamp: time.Now().Unix(),
		Data:      data,
	}

	for clientID := range o.clients {
		if o.GetLogger().IsLevelEnabled(logrus.DebugLevel) {
			data, _ := json.Marshal(data)
			o.GetLogger().Debugf("ws.Server.Send(to %d): %v", clientID, string(data))
		}

		for conn := range o.clients[clientID] {
			if err := conn.WriteJSON(msg); err != nil {
				o.GetLogger().Error(errors.Wrap(err, "send"))
			}
		}
	}
}
