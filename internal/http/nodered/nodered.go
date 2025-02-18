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
	"touchon-server/lib/messages"
	"touchon-server/lib/parallel"
)

func New() interfaces.NodeRed {
	return &NodeRedImpl{
		clients: make(map[interfaces.WebSocket]bool),
	}
}

type NodeRedImpl struct {
	mu      sync.Mutex
	clients map[interfaces.WebSocket]bool

	noderedHandlerID int
}

func (o *NodeRedImpl) Start() error {
	var err error

	o.noderedHandlerID, err = g.Msgs.Subscribe(interfaces.MessageTypeEvent, "", "", nil, func(svc interfaces.MessageSender, msg interfaces.Message) {
		o.sendAll(msg)
	})
	if err != nil {
		return errors.Wrap(err, "NodeRedImpl.Start")
	}

	return nil
}

func (o *NodeRedImpl) Shutdown() error {
	g.Msgs.Unsubscribe(o.noderedHandlerID)

	o.mu.Lock()
	for ws := range o.clients {
		if err := ws.Close(); err != nil {
			g.Logger.Error(errors.Wrap(err, "NodeRedImpl.Shutdown"))
		}
	}
	o.mu.Unlock()

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

		o.mu.Lock()
		o.clients[ws] = true // Сохраняем соединение
		o.mu.Unlock()

		ws.ReadMessages(func(ws interfaces.WebSocket, messageType int, message []byte) {
			cmd, err := messages.NewCommand("", interfaces.TargetTypeObject, 0, nil)
			if err != nil {
				g.Logger.Error(errors.Wrap(err, "ws.ReadMessages"))
				return
			}

			if err := json.Unmarshal(message, &cmd); err != nil {
				g.Logger.Error(errors.Wrap(err, "ws.ReadMessages"))
				return
			}

			if err := g.Msgs.Send(cmd); err != nil {
				g.Logger.Error(errors.Wrap(err, "ws.ReadMessages"))
				return
			}

			// TODO remove it after testing
			data, err := json.Marshal(cmd)
			if err != nil {
				g.Logger.Error(errors.Wrap(err, "ws.ReadMessages"))
				return
			}
			g.Logger.Debug(">>>> Received from NodeRed:", string(data))
		})

		if err := ws.Close(); err != nil {
			g.Logger.Error(errors.Wrap(err, "NodeRedImpl.Handler"))
		}

		o.mu.Lock()
		delete(o.clients, ws) // Удаляем соединение
		o.mu.Unlock()
	})

	if err != nil {
		g.Logger.Error(err)
	}
}

type nodeRedMsg struct {
	Payload interface{} `json:"payload"`
}

func (o *NodeRedImpl) sendAll(message interface{}) {
	if g.Logger.Level >= logrus.DebugLevel {
		data, _ := json.Marshal(message)
		g.Logger.Debug("<<<< Send to NodeRed:", string(data))
	}

	tasks := make([]parallel.Task, 0, len(o.clients))

	for ws := range o.clients {
		ws := ws // !
		tasks = append(tasks, func() {
			if err := ws.WriteJSON(nodeRedMsg{Payload: message}); err != nil {
				g.Logger.Error(errors.Wrap(err, "send"))
			}
		})
	}

	// Отправляем сообщение в несколько потоков
	if err := parallel.Do(runtime.NumCPU(), tasks, 10*time.Second); err != nil {
		g.Logger.Error(errors.Wrap(err, "NodeRedImpl.sendAll"))
	}
}
