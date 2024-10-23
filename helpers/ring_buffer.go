package helpers

import (
	"fmt"
	"sync"

	"github.com/sirupsen/logrus"
)

// NewRingBuffer создает новый кольцевой буфер.
func NewRingBuffer(bufSize int) *RingBuffer {
	return &RingBuffer{
		buf:      make([]byte, 0, bufSize),
		size:     bufSize,
		logLevel: logrus.TraceLevel,
		//logFormatter: &logrus.JSONFormatter{
		//	TimestampFormat:   "02.01.2006 15:04:05",
		//	DisableTimestamp:  false,
		//	DisableHTMLEscape: true,
		//	DataKey:           "",
		//	FieldMap: logrus.FieldMap{
		//		logrus.FieldKeyTime:  "ts",
		//		logrus.FieldKeyLevel: "lvl",
		//		logrus.FieldKeyMsg:   "msg",
		//	},
		//	CallerPrettyfier: nil,
		//	PrettyPrint:      false,
		//},
	}
}

// RingBuffer - кольцевой буфер, содержит только
// последние size байт записанных в него данных.
type RingBuffer struct {
	mu   sync.RWMutex
	buf  []byte
	size int

	logLevel logrus.Level
}

func (o *RingBuffer) Levels() []logrus.Level {
	return logrus.AllLevels[:o.logLevel+1]
}

func (o *RingBuffer) Fire(e *logrus.Entry) error {
	levels := map[logrus.Level]string{
		logrus.TraceLevel: "[TRS]",
		logrus.DebugLevel: "[DBG]",
		logrus.InfoLevel:  "[NFO]",
		logrus.WarnLevel:  "[WRN]",
		logrus.ErrorLevel: "[ERR]",
		logrus.FatalLevel: "[FTL]",
		logrus.PanicLevel: "[PNC]",
	}

	r := fmt.Sprintf("%s %s %s\n", e.Time.Format("02.01.2006 15:04:05.000"), levels[e.Level], e.Message)

	if _, err := o.Write([]byte(r)); err != nil {
		return err
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
