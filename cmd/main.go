// Менеджер объектов.
// Содержит основные методы для работы с устройствами и методы для панели администратора.

// go install github.com/swaggo/swag/cmd/swag@latest
// ./bin $ swag init --dir=../cmd,../internal,../lib --output=../docs --outputTypes=go --parseDepth=1 --parseDependency --parseInternal
// mockery --dir=internal --all --inpackage --inpackage-suffix --with-expecter
// mockery --dir=lib --all --inpackage --inpackage-suffix --with-expecter
// ./bin $ go build -C ../cmd -o ../bin/cmd && MQTT_CONNECTION_STRING="mqtt://vn:1q2w3e4r@127.0.0.1:1883/#" LOG_LEVEL=debug HTTP_ADDR=localhost:8082 ./cmd
// ./bin $ go build -C ../cmd -o ../bin/cmd && MQTT_CONNECTION_STRING="mqtt://services:12345678@10.35.16.1:1883/#" LOG_LEVEL=debug HTTP_ADDR=localhost:8082 ./cmd
// docker build --progress=plain -t ts . && docker run --rm -it --network=host --name ts ts
// docker build --progress=plain -t ts . && docker run --rm -it --network=host -e MQTT_CONNECTION_STRING="mqtt://services:12345678@10.35.16.1:1883/#" --name ts ts
//
// Получение информации о сервисе по mqtt: {
//   "target_type": "service",
//   "type": "command",
//   "name": "info"
// }

package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/mattn/go-sqlite3"
	_ "touchon-server/docs"
	"touchon-server/internal/context"
	httpServer "touchon-server/internal/http"
	mqttService "touchon-server/internal/mqtt"
	"touchon-server/internal/objects"
	"touchon-server/internal/scripts"
	"touchon-server/internal/store"
	memStore "touchon-server/internal/store/memstore"
	"touchon-server/internal/store/sqlstore"
	httpClient "touchon-server/lib/http/client"
	mqttClient "touchon-server/lib/mqtt/client"
	"touchon-server/lib/service"

	_ "touchon-server/internal/object/PortMegaD"
	_ "touchon-server/internal/object/Regulator"
	_ "touchon-server/internal/object/Relay"
	_ "touchon-server/internal/object/Sensor/bh1750"
	_ "touchon-server/internal/object/Sensor/bme280"
	_ "touchon-server/internal/object/Sensor/cs"
	_ "touchon-server/internal/object/Sensor/ds18b20"
	_ "touchon-server/internal/object/Sensor/htu21d"
	_ "touchon-server/internal/object/Sensor/motion"
	_ "touchon-server/internal/object/Sensor/outdoor"
	_ "touchon-server/internal/object/Sensor/presence"
	_ "touchon-server/internal/object/Sensor/scd4x"
	_ "touchon-server/internal/object/SensorValue"
	_ "touchon-server/internal/object/WirenBoard/wb_mrm2_mini"
	// Объявляем объекты для их регистрации в реестре
	_ "touchon-server/internal/object/GenericInput"
	_ "touchon-server/internal/object/MegaD"
	_ "touchon-server/internal/object/Modbus"
	_ "touchon-server/internal/object/Onokom/Conditioner"
)

func init() {
	// Подключаем расширения для sqlite3
	path := fmt.Sprintf("sqlean/%s_%s/unicode", runtime.GOOS, runtime.GOARCH)
	sql.Register("sqlite3_with_extensions", &sqlite3.SQLiteDriver{Extensions: []string{path}})
}

var defaults = map[string]string{
	"http_addr":              "0.0.0.0:8082",
	"action_router_addr":     "0.0.0.0:8083",
	"database_url":           "./object-manager_db.sqlite?_foreign_keys=true",
	"server_key":             "c041d36e381a835afce48c91686370c8",
	"mqtt_connection_string": "mqtt://services:12345678@mqtt:1883/#",
	"log_level":              "fatal",
	"version":                "0.1",
	"service_name":           "object_manager",
	"mqtt_max_travel_time":   "50ms",
}

const banner = `
 ██████╗ ██████╗      ██╗███████╗ ██████╗████████╗              ███╗   ███╗ █████╗ ███╗   ██╗ █████╗  ██████╗ ███████╗██████╗ 
██╔═══██╗██╔══██╗     ██║██╔════╝██╔════╝╚══██╔══╝              ████╗ ████║██╔══██╗████╗  ██║██╔══██╗██╔════╝ ██╔════╝██╔══██╗
██║   ██║██████╔╝     ██║█████╗  ██║        ██║       █████╗    ██╔████╔██║███████║██╔██╗ ██║███████║██║  ███╗█████╗  ██████╔╝
██║   ██║██╔══██╗██   ██║██╔══╝  ██║        ██║       ╚════╝    ██║╚██╔╝██║██╔══██║██║╚██╗██║██╔══██║██║   ██║██╔══╝  ██╔══██╗
╚██████╔╝██████╔╝╚█████╔╝███████╗╚██████╗   ██║                 ██║ ╚═╝ ██║██║  ██║██║ ╚████║██║  ██║╚██████╔╝███████╗██║  ██║
 ╚═════╝ ╚═════╝  ╚════╝ ╚══════╝ ╚═════╝   ╚═╝                 ╚═╝     ╚═╝╚═╝  ╚═╝╚═╝  ╚═══╝╚═╝  ╚═╝ ╚═════╝ ╚══════╝╚═╝  ╚═╝

`

// Version заполняется компилятором
var Version string

// BuildAt заполняется компилятором
var BuildAt string

func main() {
	cfg, logger, rb, db, err := service.Prolog(banner, defaults, Version, BuildAt)
	check(err)
	context.Logger = logger
	context.Config = cfg

	store.I = sqlstore.New(db)

	// Инициализация клиента для MQTT
	mqttClient.I, err = mqttClient.New(cfg["service_name"], cfg["mqtt_connection_string"], 10*time.Second, 3, logger)
	check(err)
	mqttClient.I.SetIgnoreSelfMsgs(false)

	// Создаем скриптовый движок
	scripts.I = scripts.NewScripts(10*time.Second, objects.NewExecutor())

	mqttService.I, err = mqttService.New(1000, 4)
	check(err)

	// Загружает все объекты БД в память
	memStore.I, err = memStore.New()
	check(err)

	check(memStore.I.Start())

	check(mqttService.I.Start())

	httpClient.I = httpClient.New()

	httpServer.I, err = httpServer.New(rb)
	check(err)

	// Старт HTTP API сервера
	check(httpServer.I.Start(cfg["http_addr"]))

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-sig
	logger.Info("Получен сигнал на завершение...")

	if err := httpServer.I.Shutdown(); err != nil {
		logger.Error(err)
	}

	if err := mqttService.I.Shutdown(); err != nil {
		logger.Error(err)
	}

	if err := memStore.I.Shutdown(); err != nil {
		logger.Error(err)
	}
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
