// TouchOn Server
//
// ./bin $ go generate -C ../cmd
// ./bin $ go build -C ../cmd -o ../bin/cmd && TOKEN_SECRET=disable_auth ./cmd
// docker build --progress=plain -t ts . && docker run --rm -it --network=host --name ts ts
//
// GOOSE_DRIVER=sqlite3 GOOSE_DBSTRING=bin/db.sqlite GOOSE_MIGRATION_DIR=migrations go tool goose status
// GOOSE_DRIVER=sqlite3 GOOSE_DBSTRING=bin/db.sqlite GOOSE_MIGRATION_DIR=migrations go tool goose create <migration_name> sql

package main

//go:generate find ../internal -name '*_mock.go' -delete
//go:generate find ../lib -name '*_mock.go' -delete
//go:generate go tool swag init --dir=../cmd,../internal,../lib --output=../docs --outputTypes=go --parseDepth=1 --parseDependency --parseInternal
//go:generate go tool mockery --dir=../internal --all --inpackage --inpackage-suffix --with-expecter
//go:generate go tool mockery --dir=../lib --all --inpackage --inpackage-suffix --with-expecter

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
	"touchon-server/internal/action_router"
	"touchon-server/internal/cron"
	"touchon-server/internal/g"
	httpServer "touchon-server/internal/http"
	"touchon-server/internal/http/nodered"
	_ "touchon-server/internal/object/GenericInput"
	_ "touchon-server/internal/object/ImpulseCounter"
	_ "touchon-server/internal/object/MegaD"
	_ "touchon-server/internal/object/Onokom/Conditioner"
	_ "touchon-server/internal/object/PortMegaD"
	_ "touchon-server/internal/object/RS485"
	_ "touchon-server/internal/object/Regulator"
	_ "touchon-server/internal/object/Relay"
	_ "touchon-server/internal/object/Sensor/bh1750"
	_ "touchon-server/internal/object/Sensor/bme280"
	_ "touchon-server/internal/object/Sensor/cs"
	_ "touchon-server/internal/object/Sensor/ds18b20"
	_ "touchon-server/internal/object/Sensor/htu21d"
	_ "touchon-server/internal/object/Sensor/htu31d"
	_ "touchon-server/internal/object/Sensor/motion"
	_ "touchon-server/internal/object/Sensor/outdoor"
	_ "touchon-server/internal/object/Sensor/presence"
	_ "touchon-server/internal/object/Sensor/scd4x"
	_ "touchon-server/internal/object/SensorValue"
	_ "touchon-server/internal/object/Server"
	_ "touchon-server/internal/object/WirenBoard/wb_mrm2_mini"
	"touchon-server/internal/objects"
	"touchon-server/internal/scripts"
	"touchon-server/internal/store"
	memStore "touchon-server/internal/store/memstore"
	"touchon-server/internal/store/sqlstore"
	"touchon-server/internal/ws"
	"touchon-server/lib/messages"
	"touchon-server/lib/service"
	_ "touchon-server/migrations"
)

func init() {
	// Подключаем расширения для sqlite3
	path := fmt.Sprintf("sqlean/%s_%s/unicode", runtime.GOOS, runtime.GOARCH)
	sql.Register("sqlite3_with_extensions", &sqlite3.SQLiteDriver{Extensions: []string{path}})
}

var defaults = map[string]string{
	"http_addr":    "0.0.0.0:8081",
	"database_url": "db.sqlite?_foreign_keys=true",
	"server_key":   "c041d36e381a835afce48c91686370c8",
	"log_level":    "debug",
	"version":      "0.1",

	"access_token_ttl":    "30m",
	"refresh_token_ttl":   "43200m",
	"token_secret":        "Alli80ed!",
	"ws_addr":             "0.0.0.0:8092",
	"push_sender_address": "http://localhost:8088",
}

const banner = `
████████╗ ██████╗ ██╗   ██╗ ██████╗██╗  ██╗         ███╗   ██╗    ███████╗███████╗██████╗ ██╗   ██╗███████╗██████╗ 
╚══██╔══╝██╔═══██╗██║   ██║██╔════╝██║  ██║ ██████╗ ████╗  ██║    ██╔════╝██╔════╝██╔══██╗██║   ██║██╔════╝██╔══██╗
   ██║   ██║   ██║██║   ██║██║     ███████║██║   ██║██╔██╗ ██║    ███████╗█████╗  ██████╔╝██║   ██║█████╗  ██████╔╝
   ██║   ██║   ██║██║   ██║██║     ██╔══██║██║   ██║██║╚██╗██║    ╚════██║██╔══╝  ██╔══██╗╚██╗ ██╔╝██╔══╝  ██╔══██╗
   ██║   ╚██████╔╝╚██████╔╝╚██████╗██║  ██║╚██████╔╝██║ ╚████║    ███████║███████╗██║  ██║ ╚████╔╝ ███████╗██║  ██║
   ╚═╝    ╚═════╝  ╚═════╝  ╚═════╝╚═╝  ╚═╝ ╚═════╝ ╚═╝  ╚═══╝    ╚══════╝╚══════╝╚═╝  ╚═╝  ╚═══╝  ╚══════╝╚═╝  ╚═╝

`

// Version заполняется компилятором
var Version string

// BuildAt заполняется компилятором
var BuildAt string

// @securityDefinitions.apikey TokenAuth
// @in header
// @name token
func main() {
	cfg, logger, rb, db, err := service.Prolog(banner, defaults, Version, BuildAt)
	check(err)
	g.Logger = logger
	g.Config = cfg

	// Создаем хранилище
	store.I = sqlstore.New(db)
	check(prepareDB())

	// Создаем экземпляр вебсокет сервера для мобильных приложений
	ws.I, err = ws.New()
	check(err)

	check(ws.I.Start(cfg["ws_addr"]))

	// Создаем шину сообщений
	g.Msgs, err = messages.NewService(runtime.NumCPU(), 2000)
	check(err)

	// Создаем скриптовый движок
	scripts.I = scripts.NewScripts(10*time.Second, objects.NewExecutor())

	// Загружает все объекты БД в память
	memStore.I, err = memStore.New()
	check(err)

	// Создаем штатные обработчики сообщений
	action_router.I = action_router.New()

	check(memStore.I.Start())

	check(g.Msgs.Start())

	check(scripts.I.Start())

	check(action_router.I.Start())

	// Создаем вебсокет-мост между NodeRed и шиной сообщений
	g.NodeRed = nodered.New()

	// Создаем основной API-сервер
	g.HttpServer, err = httpServer.New(rb)
	check(err)

	// Старт HTTP API сервера
	check(g.HttpServer.Start(cfg["http_addr"]))

	check(g.NodeRed.Start())

	// Создаем планировщик задач
	sch, err := cron.New()
	check(err)

	check(sch.Start())

	// Ждем от ОС сигнала на завершение работы сервиса
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-sig
	logger.Info("Получен сигнал на завершение...")

	if err := sch.Shutdown(); err != nil {
		logger.Error(err)
	}

	if err := g.NodeRed.Shutdown(); err != nil {
		logger.Error(err)
	}

	if err := g.HttpServer.Shutdown(); err != nil {
		logger.Error(err)
	}

	if err := action_router.I.Shutdown(); err != nil {
		logger.Error(err)
	}

	if err := scripts.I.Shutdown(); err != nil {
		logger.Error(err)
	}

	if err := g.Msgs.Shutdown(); err != nil {
		logger.Error(err)
	}

	if err := ws.I.Shutdown(); err != nil {
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
