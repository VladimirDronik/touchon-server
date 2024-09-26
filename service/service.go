package service

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"object-manager/internal/touchon-server/config"
	"object-manager/internal/touchon-server/helpers"
	"object-manager/internal/touchon-server/mqtt/messages"
)

func Prolog(banner string, configDefaults map[string]string, version, buildAt string) (map[string]string, *logrus.Logger, fmt.Stringer, *gorm.DB, error) {
	fmt.Print(banner)

	cfg, err := config.New(configDefaults)
	if err != nil {
		return nil, nil, nil, nil, errors.Wrap(err, "Prolog")
	}

	if version != "" {
		cfg["version"] = version
	}

	if buildAt != "" {
		cfg["build_at"] = buildAt
	}

	fmt.Print("Version: ", cfg["version"], "\n\n\n")

	messages.Publisher = cfg["service_name"]

	logger, err := helpers.NewLogger(cfg["log_level"])
	if err != nil {
		return nil, nil, nil, nil, errors.Wrap(err, "Prolog")
	}

	// Выводим логи в консоль и кольцевой буфер
	rb := helpers.NewRingBuffer(100 * 1024)
	logger.AddHook(rb)

	if logger.Level == logrus.DebugLevel {
		fmt.Printf("==========================================================================\n" +
			"=================== SERVICE IS RUNNING ON DEBUG MODE =====================\n" +
			"==========================================================================\n\n\n")

		logger.Infof("ENV: %#v", cfg)
	}

	db, err := helpers.NewDB(cfg["database_url"])
	if err != nil {
		return nil, nil, nil, nil, errors.Wrap(err, "Prolog")
	}

	return cfg, logger, rb, db, nil
}
