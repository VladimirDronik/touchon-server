package service

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"touchon-server/lib/config"
	"touchon-server/lib/helpers"
	"touchon-server/lib/info"
	"touchon-server/lib/models"
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

	info.Config = cfg

	logger, err := models.NewLogger(cfg["log_level"])
	if err != nil {
		return nil, nil, nil, nil, errors.Wrap(err, "Prolog")
	}

	// Выводим логи в консоль и кольцевой буфер
	rb := models.NewRingBuffer(100*1024, &models.LogFormatter{})
	logger.AddHook(rb)

	logger.Debugf("\n==========================================================================\n" +
		"=================== SERVICE IS RUNNING ON DEBUG MODE =====================\n" +
		"==========================================================================\n\n\n")

	logger.Debugf("ENV: %#v", cfg)

	db, err := helpers.NewDB(cfg["database_url"], logger)
	if err != nil {
		return nil, nil, nil, nil, errors.Wrap(err, "Prolog")
	}

	return cfg, logger, rb, db, nil
}
