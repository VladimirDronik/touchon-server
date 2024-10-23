package models

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

type LogFormatter struct{}

func (o *LogFormatter) Format(e *logrus.Entry) ([]byte, error) {
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

	return []byte(r), nil
}
