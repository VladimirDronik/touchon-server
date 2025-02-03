package models

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func NewLogger(logLevel string) (*logrus.Logger, error) {
	logger := logrus.New()

	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		return nil, errors.Wrap(err, "NewLogger")
	}

	logger.SetLevel(level)
	logger.SetFormatter(&LogFormatter{})

	return logger, nil
}
