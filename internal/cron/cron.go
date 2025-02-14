package cron

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/pkg/errors"
	"touchon-server/internal/g"
	"touchon-server/internal/store"
	"touchon-server/lib/interfaces"
	"touchon-server/lib/messages"
)

func New() (*Scheduler, error) {
	ctx, cancel := context.WithCancel(context.Background())

	return &Scheduler{
		ctx:    ctx,
		cancel: cancel,
	}, nil
}

type Scheduler struct {
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

func (o *Scheduler) Start() error {
	o.wg.Add(1)

	go func() {
		defer o.wg.Done()

		tasks, err := store.I.CronRepo().GetEnabledTasks()
		if err != nil {
			g.Logger.Error(err)
		}

		for {
			select {
			case <-o.ctx.Done():
				return
			default:
				// Раз в минуту запрашиваем задачи из базы
				if _, _, sec := time.Now().Clock(); sec == 0 {
					tasks, err = store.I.CronRepo().GetEnabledTasks()
					if err != nil {
						g.Logger.Error(err)
						continue
					}
				}

				m := make(map[string][]*interfaces.CronTask, len(tasks))
				for _, task := range tasks {
					m[task.Period] = append(m[task.Period], task)
				}

				if err := o.task(m); err != nil {
					g.Logger.Error(err)
					continue
				}
			}

			t1 := time.Now()
			t2 := t1.Truncate(time.Second).Add(time.Second) // обрезаем миллисекунды и прибавляем секунду
			time.Sleep(t2.Sub(t1))                          // спим ровно до начала следующей секунды
		}
	}()

	return nil
}

var taskConditions = map[string]func(h, m, s int) bool{
	"1s":  func(h, m, s int) bool { return true },
	"5s":  func(h, m, s int) bool { return s%5 == 0 },
	"10s": func(h, m, s int) bool { return s%10 == 0 },
	"15s": func(h, m, s int) bool { return s%15 == 0 },
	"20s": func(h, m, s int) bool { return s%20 == 0 },
	"30s": func(h, m, s int) bool { return s%30 == 0 },
	"1m":  func(h, m, s int) bool { return s == 0 },
	"5m":  func(h, m, s int) bool { return s == 0 && m%5 == 0 },
	"10m": func(h, m, s int) bool { return s == 0 && m%10 == 0 },
	"15m": func(h, m, s int) bool { return s == 0 && m%15 == 0 },
	"20m": func(h, m, s int) bool { return s == 0 && m%20 == 0 },
	"30m": func(h, m, s int) bool { return s == 0 && m%30 == 0 },
	"1h":  func(h, m, s int) bool { return s == 0 && m == 0 },
	"2h":  func(h, m, s int) bool { return s == 0 && m == 0 && h%2 == 0 },
	"3h":  func(h, m, s int) bool { return s == 0 && m == 0 && h%3 == 0 },
	"4h":  func(h, m, s int) bool { return s == 0 && m == 0 && h%4 == 0 },
	"6h":  func(h, m, s int) bool { return s == 0 && m == 0 && h%6 == 0 },
	"12h": func(h, m, s int) bool { return s == 0 && m == 0 && h%12 == 0 },
}

func (o *Scheduler) task(tasks map[string][]*interfaces.CronTask) error {
	h, m, s := time.Now().Clock()

	for period, cond := range taskConditions {
		if cond(h, m, s) {
			if err := o.doAction(tasks[period]); err != nil {
				return errors.Wrap(err, "task")
			}
		}

	}

	return nil
}

func (o *Scheduler) doAction(tasks []*interfaces.CronTask) error {
	for _, task := range tasks {
		for _, act := range task.Actions {
			switch act.Type {
			case interfaces.ActionTypeDelay:
				v, ok := act.Args["duration"]
				if !ok {
					return errors.Wrap(errors.New("duration not found"), "doAction")
				}

				s, ok := v.(string)
				if !ok {
					return errors.Wrap(errors.New("duration is not string"), "doAction")
				}

				d, err := time.ParseDuration(s)
				if err != nil {
					return errors.Wrap(err, "doAction")
				}

				time.Sleep(d)

			case interfaces.ActionTypeMethod:
				msg, err := messages.NewCommand(act.Name, act.TargetType, act.TargetID, act.Args)
				if err != nil {
					return errors.Wrap(err, "doAction")
				}

				if err := g.Msgs.Send(msg); err != nil {
					return errors.Wrap(err, "doAction")
				}

			default:
				return errors.Wrap(errors.Errorf("unknown action type %q", act.Type), "doAction")
			}
		}
	}

	return nil
}

type payloadStruct struct {
	IdObject int    `json:"object_id,omitempty"`
	IdItem   int    `json:"item_id,omitempty"`
	Method   string `json:"method"`
	State    string `json:"state,omitempty"`
	Params   []byte `json:"params,omitempty"`
}

func GetTopicAndParams(target string, value string) (string, interface{}) {
	var data struct {
		ObjectID int    `json:"object_id,omitempty"`
		ItemID   int    `json:"item_id,omitempty"`
		Method   string `json:"method"`
		State    string `json:"state,omitempty"`
		Params   []byte `json:"params,omitempty"`
	}

	if err := json.Unmarshal([]byte(value), &data); err != nil {
		println(err)
	}

	params, err := json.Marshal(data.Params)
	if err != nil {
		println(err)
	}

	payload := payloadStruct{IdObject: data.ObjectID,
		IdItem: data.ItemID,
		State:  data.State,
		Params: params,
		Method: data.Method,
	}
	payloadResp, _ := json.Marshal(payload)

	topic := "action_router/" + target + "/method"
	return topic, payloadResp
}

func (o *Scheduler) Shutdown() error {
	o.cancel()
	o.wg.Wait()
	return nil
}
