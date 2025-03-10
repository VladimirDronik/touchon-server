package helpers

import (
	"database/sql"
	"math"
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/pressly/goose/v3"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"touchon-server/migrations"
)

func FileIsExists(path string) bool {
	s, err := os.Stat(path)
	return !errors.Is(err, os.ErrNotExist) && !s.IsDir()
}

func Round(v float32) float32 {
	return float32(math.Round(float64(v)*10)) / 10
}

// Add into main.go:
//
//	func init() {
//		var exts = map[string]string{"darwin": ".dylib", "linux": ".so", "windows": ".dll"}
//		path := fmt.Sprintf("sqlean/%s_%s/unicode%s", runtime.GOOS, runtime.GOARCH, exts[runtime.GOOS])
//		sql.Register("sqlite3_with_extensions", &sqlite3.SQLiteDriver{Extensions: []string{path}})
//	}
func NewDB(connString string, logger *logrus.Logger) (*gorm.DB, error) {
	// Запускаем миграции (с отключенной опцией _foreign_keys=true)
	path := strings.Split(connString, "?")
	if len(path) == 0 {
		return nil, errors.Wrap(errors.New("connString is empty"), "NewDB")
	}

	sqlDB, err := sql.Open("sqlite3", path[0])

	goose.SetBaseFS(migrations.EmbedMigrations)
	goose.SetLogger(logger)

	if err := goose.SetDialect(string(goose.DialectSQLite3)); err != nil {
		return nil, errors.Wrap(err, "NewDB")
	}

	logger.Info("Текущая версия данных:")
	if err := goose.Version(sqlDB, "."); err != nil {
		return nil, errors.Wrap(err, "NewDB")
	}

	logger.Info("Запускаем миграции...")
	if err := goose.Up(sqlDB, "."); err != nil {
		return nil, errors.Wrap(err, "NewDB")
	}

	// Подключаемся к БД с указанными (в строке подключения) параметрами
	db, err := gorm.Open(sqlite.Dialector{DriverName: "sqlite3_with_extensions", DSN: connString}, &gorm.Config{})
	if err != nil {
		return nil, errors.Wrap(err, "NewDB")
	}

	return db, nil
}
