package helpers

import (
	"sync/atomic"
	"time"
)

// NewTimer создает таймер, после истечения которого выполняется функция.
// Таймер можно сбросить или остановить.
func NewTimer(dur time.Duration, f func()) *Timer {
	o := &Timer{
		dur:     dur,
		abort:   make(chan struct{}),
		f:       f,
		stopped: atomic.Bool{},
	}

	o.stopped.Store(true)

	return o
}

type Timer struct {
	dur     time.Duration
	abort   chan struct{}
	f       func()
	stopped atomic.Bool
}

func (o *Timer) Start() {
	o.stopped.Store(false)
	o.Reset()
}

func (o *Timer) Reset() {
	// Если таймер остановлен, то не перезапускаем его.
	if o.stopped.Load() {
		return
	}

	select {
	case o.abort <- struct{}{}:
	default:
	}

	o.setTimer()
}

func (o *Timer) Stop() {
	o.stopped.Store(true)

	select {
	case o.abort <- struct{}{}:
	default:
	}
}

func (o *Timer) setTimer() {
	go func() {
		t := time.NewTimer(o.dur)
		defer t.Stop()

		select {
		case <-t.C:
			o.f()
		case <-o.abort:
		}
	}()
}

func (o *Timer) GetDuration() time.Duration {
	return o.dur
}
