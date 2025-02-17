package ws

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/fasthttp/websocket"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
	"touchon-server/internal/g"
	"touchon-server/internal/token"
	"touchon-server/lib/helpers"
	"touchon-server/lib/http/server"
)

// Global instance
var I *Server

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

	authToken := string(ctx.Request.Header.Peek("token"))
	if authToken == "" {
		ctx.Error(http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	// Проверяем и извлекаем ID клиента из токена
	clientID, err := token.KeysExtract(authToken, o.GetConfig()["token_secret"])
	if err != nil {
		ctx.Error(http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	var upgrader = websocket.FastHTTPUpgrader{
		CheckOrigin: func(ctx *fasthttp.RequestCtx) bool { return true }, // Пропускаем любой запрос
	}

	err = upgrader.Upgrade(ctx, func(ws *websocket.Conn) {
		defer ws.Close()

		// Устанавливаем PongHandler
		ws.SetPongHandler(func(appData string) error {
			return nil
		})

		// Устанавливаем обработчик Ping сообщений
		ws.SetPingHandler(func(appData string) error {
			//o.GetLogger().Debugf("ws.Server: received ping from client %d", clientID)
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

		if o.clients[clientID] == nil {
			o.clients[clientID] = map[*websocket.Conn]bool{}
		}
		o.GetLogger().Debugf("ws.Server: client %d connected", clientID)
		o.clients[clientID][ws] = true        // Сохраняем соединение, используя его как ключ
		defer delete(o.clients[clientID], ws) // Удаляем соединение

		for {
			mt, message, err := ws.ReadMessage()
			if err != nil || mt == websocket.CloseMessage {
				o.GetLogger().Debugf("ws.Server: client %d disconnected", clientID)
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

func (o *Server) Send(message interface{}) {
	for clientID := range o.clients {
		o.send(clientID, message)
	}
}

func (o *Server) send(clientID int, message interface{}) {
	for conn := range o.clients[clientID] {
		if o.GetLogger().Level >= logrus.DebugLevel {
			data, _ := json.Marshal(message)
			o.GetLogger().Debugf("ws.Server.Send(to %d): %v", clientID, string(data))
		}

		if err := conn.WriteJSON(message); err != nil {
			o.GetLogger().Error(errors.Wrap(err, "send"))
		}
	}
}
