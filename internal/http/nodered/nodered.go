package nodered

import (
	"encoding/json"
	"runtime"
	"sync"
	"time"

	"github.com/fasthttp/websocket"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
	"touchon-server/internal/g"
	"touchon-server/lib/interfaces"
	"touchon-server/lib/parallel"
)

func New() interfaces.NodeRed {
	return &NodeRedImpl{
		clients: make(map[interfaces.WebSocket]bool),
	}
}

type NodeRedImpl struct {
	mu               sync.Mutex
	clients          map[interfaces.WebSocket]bool
	noderedHandlerID int
}

func (o *NodeRedImpl) Start() error {
	var err error

	o.noderedHandlerID, err = g.Msgs.Subscribe("", "", "", nil, func(svc interfaces.MessageSender, msg interfaces.Message) {
		o.sendAll(msg)
	})
	if err != nil {
		return errors.Wrap(err, "NodeRedImpl.Start")
	}

	return nil
}

func (o *NodeRedImpl) Shutdown() error {
	g.Msgs.Unsubscribe(o.noderedHandlerID)

	for ws := range o.clients {
		if err := ws.Close(); err != nil {
			g.Logger.Error(errors.Wrap(err, "NodeRedImpl.Shutdown"))
		}
	}

	return nil
}

func (o *NodeRedImpl) Handler(ctx *fasthttp.RequestCtx) {
	upgrader := websocket.FastHTTPUpgrader{
		CheckOrigin: func(ctx *fasthttp.RequestCtx) bool { return true }, // Пропускаем любой запрос
	}

	err := upgrader.Upgrade(ctx, func(conn *websocket.Conn) {
		ws, err := NewWS(conn, time.Duration(30)*time.Second)
		if err != nil {
			g.Logger.Error(errors.Wrap(err, "NodeRedImpl.Handler"))
			return
		}

		o.clients[ws] = true  // Сохраняем соединение
		delete(o.clients, ws) // Удаляем соединение

		ws.ReadMessages(func(ws interfaces.WebSocket, messageType int, message []byte) {
			g.Logger.Debug(string(message))
		})

		if err := ws.Close(); err != nil {
			g.Logger.Error(errors.Wrap(err, "NodeRedImpl.Handler"))
		}
	})

	if err != nil {
		g.Logger.Error(err)
	}
}

func (o *NodeRedImpl) sendAll(message interface{}) {
	tasks := make([]parallel.Task, 0, len(o.clients))

	for ws := range o.clients {
		ws := ws // !
		tasks = append(tasks, func() {
			if g.Logger.Level >= logrus.DebugLevel {
				data, _ := json.Marshal(message)
				g.Logger.Debugf("Send to NodeRed: %v", string(data))
			}

			if err := ws.WriteJSON(message); err != nil {
				g.Logger.Error(errors.Wrap(err, "send"))
			}
		})
	}

	// Отправляем сообщение в несколько потоков
	if err := parallel.Do(runtime.NumCPU(), tasks, 10*time.Second); err != nil {
		g.Logger.Error(errors.Wrap(err, "NodeRedImpl.sendAll"))
	}
}
