package helpers

import "time"

// NewTimer создает таймер, после истечения которого выполняется функция.
// Таймер можно сбросить или остановить.
func NewTimer(dur time.Duration, f func()) *Timer {
	o := &Timer{
		dur:   dur,
		abort: make(chan struct{}),
		f:     f,
	}

	return o
}

type Timer struct {
	dur   time.Duration
	abort chan struct{}
	f     func()
}

func (o *Timer) Start() {
	o.Reset()
}

func (o *Timer) Reset() {
	o.Stop()
	o.setTimer()
}

func (o *Timer) Stop() {
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
