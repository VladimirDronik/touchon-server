package models

import (
	"sync"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// NewRingBuffer создает новый кольцевой буфер.
func NewRingBuffer(bufSize int, formatter logrus.Formatter) *RingBuffer {
	return &RingBuffer{
		buf:          make([]byte, 0, bufSize),
		size:         bufSize,
		logLevel:     logrus.TraceLevel,
		logFormatter: formatter,
	}
}

// RingBuffer - кольцевой буфер, содержит только
// последние size байт записанных в него данных.
type RingBuffer struct {
	mu   sync.RWMutex
	buf  []byte
	size int

	logLevel     logrus.Level
	logFormatter logrus.Formatter
}

func (o *RingBuffer) Levels() []logrus.Level {
	return logrus.AllLevels[:o.logLevel+1]
}

func (o *RingBuffer) Fire(e *logrus.Entry) error {
	data, err := o.logFormatter.Format(e)
	if err != nil {
		return errors.Wrap(err, "RingBuffer.Fire")
	}

	if _, err := o.Write(data); err != nil {
		return errors.Wrap(err, "RingBuffer.Fire")
	}

	return nil
}

func (o *RingBuffer) String() string {
	o.mu.RLock()
	defer o.mu.RUnlock()
	return string(o.buf)
}

func (o *RingBuffer) Write(p []byte) (n int, err error) {
	o.mu.Lock()
	defer o.mu.Unlock()

	switch {
	// Размер сообщения превышает размер буфера,
	// обрезаем начало сообщения
	case len(p) >= o.size:
		copy(o.buf, p[len(p)-o.size:])

	// Сообщение полностью помещается в буфер
	// вместе с имеющимися данными в буфере
	case len(o.buf)+len(p) <= o.size:
		o.buf = append(o.buf, p...)

	// Сообщение меньше размера буфера, но не помещается
	// вместе с имеющимися данными в буфере,
	// обрезаем начало имеющихся данных
	default:
		trimSize := (len(o.buf) + len(p)) - o.size
		copy(o.buf, o.buf[trimSize:])
		o.buf = o.buf[:len(o.buf)-trimSize]
		o.buf = append(o.buf, p...)
	}

	return len(p), nil
}
