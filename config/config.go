package config

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
	"object-manager/internal/touchon-server/helpers"
)

// New Загружает настройки сервиса из toml-файла, переопределяет их из ENV,
// затем проверяет на валидность.
func New(defaults map[string]string) (map[string]string, error) {
	var tomlPath string
	flag.StringVar(&tomlPath, "config", "", "Path to config file")
	flag.Parse()

	envs := os.Environ()

	r := make(map[string]string, len(envs))
	for k, v := range defaults {
		r[strings.ToLower(k)] = v
	}

	for _, env := range envs {
		kv := strings.SplitN(env, "=", 2)
		if len(kv) != 2 {
			continue
		}

		if k, v := kv[0], kv[1]; v != "" {
			r[strings.ToLower(k)] = v
		}
	}

	if helpers.FileIsExists(tomlPath) {
		o := make(map[string]interface{})
		if _, err := toml.DecodeFile(tomlPath, &o); err != nil {
			return nil, errors.Wrap(err, "config.New")
		}

		for k, v := range o {
			k = strings.ToLower(k)

			switch v := v.(type) {
			case string:
				r[k] = v
			case int:
				r[k] = strconv.Itoa(v)
			case bool:
				r[k] = strconv.FormatBool(v)
			case float64:
				r[k] = fmt.Sprintf("%.1f", v)
			default:
				return nil, errors.Wrap(errors.Errorf("unexpected value %v type %T", v, v), "config.New")
			}
		}
	}

	return r, nil
}
