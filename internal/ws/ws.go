package ws

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/VladimirDronik/touchon-server/helpers"
	"github.com/VladimirDronik/touchon-server/http/server"
	"github.com/fasthttp/websocket"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
	_ "translator/docs"
	"translator/internal/store"
	"translator/internal/token"
)

func New(cfg map[string]string, store store.Store, logger *logrus.Logger) (*Server, error) {
	switch {
	case cfg == nil:
		return nil, errors.Wrap(errors.New("cfg is nil"), "http.New")
	case logger == nil:
		return nil, errors.Wrap(errors.New("logger is nil"), "http.New")
	case store == nil:
		return nil, errors.Wrap(errors.New("store is nil"), "http.New")
	}

	baseServer, err := server.New("WS", cfg, nil, logger)
	if err != nil {
		return nil, errors.Wrap(err, "http.New")
	}

	o := &Server{
		Server:  baseServer,
		store:   store,
		clients: make(map[int]map[*websocket.Conn]bool),
	}

	o.GetServer().Handler = o.handler

	return o, nil
}

type Server struct {
	*server.Server
	store   store.Store
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
			//o.GetLogger().Infof("ws.Server: received pong from client %d", clientID)
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
				//o.GetLogger().Debugf("ws.Server: sent ping to client %d", clientID)
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

			go o.handleMessage(message)

			//err = ws.WriteMessage(mt, message)
			//if err != nil {
			//	log.Println("write:", err)
			//	break
			//}
		}
	})

	if err != nil {
		o.GetLogger().Error(err)
		//if _, ok := err.(websocket.HandshakeError); ok {
		//	log.Println(err)
		//}
	}
}

func (o *Server) Send(message interface{}, clientIDs ...int) {
	if len(clientIDs) > 0 {
		for _, clientID := range clientIDs {
			o.send(clientID, message)
		}
	} else {
		for clientID := range o.clients {
			o.send(clientID, message)
		}
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

func (o *Server) handleMessage(message []byte) {
	// todo
}
