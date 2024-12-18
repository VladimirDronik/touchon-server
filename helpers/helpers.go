package helpers

import (
	"math"
	"os"
	"strconv"
	"sync"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// GetParam отдает любой параметр, который был получен из GET запроса
func GetParam(ctx *fasthttp.RequestCtx, paramName string) string {
	return string(ctx.QueryArgs().Peek(paramName))
}

func GetIntParam(ctx *fasthttp.RequestCtx, paramName string) (int, error) {
	if s := GetParam(ctx, paramName); s != "" {
		v, err := strconv.Atoi(s)
		if err != nil {
			return 0, errors.Wrapf(err, "GetIntParam(%s)", paramName)
		}

		return v, nil
	}

	return 0, nil
}

func GetUintParam(ctx *fasthttp.RequestCtx, paramName string) (int, error) {
	v, err := GetIntParam(ctx, paramName)
	if err != nil {
		return 0, errors.Wrap(err, "GetUintParam")
	}

	if v < 0 {
		return 0, errors.Errorf("GetUintParam(%s) < 0", paramName)
	}

	return v, nil
}

func GetBoolParam(ctx *fasthttp.RequestCtx, paramName string) (bool, error) {
	if s := GetParam(ctx, paramName); s != "" {
		v, err := strconv.ParseBool(s)
		if err != nil {
			return false, errors.Wrapf(err, "GetBoolParam(%s)", paramName)
		}

		return v, nil
	}

	return false, nil
}

func GetPathParam(ctx *fasthttp.RequestCtx, paramName string) string {
	s, _ := ctx.UserValue(paramName).(string)
	return s
}

func GetUintPathParam(ctx *fasthttp.RequestCtx, paramName string) (int, error) {
	s, _ := ctx.UserValue(paramName).(string)
	if s == "" {
		return 0, errors.Wrapf(errors.New("param not found"), "GetUintPathParam(%s)", paramName)
	}

	v, err := strconv.Atoi(s)
	if err != nil {
		return 0, errors.Wrapf(err, "GetUintPathParam(%s)", paramName)
	}

	if v < 0 {
		return 0, errors.Wrapf(err, "GetUintPathParam(%s) < 0", paramName)
	}

	return v, nil
}

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
func NewDB(connString string) (*gorm.DB, error) {
	d := sqlite.Dialector{DriverName: "sqlite3_with_extensions", DSN: connString}
	db, err := gorm.Open(d, &gorm.Config{})
	if err != nil {
		return nil, errors.Wrap(err, "NewDB")
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, errors.Wrap(err, "NewDB")
	}

	// TODO для чего здесь сетевые настройки для sqlite?
	//sqlDB.SetConnMaxLifetime(time.Second * config.MaxLifetime)
	//sqlDB.SetConnMaxIdleTime(time.Second * config.MaxIdleTime)
	//sqlDB.SetMaxOpenConns(config.MaxOpenConns)
	//sqlDB.SetMaxIdleConns(config.MaxIdleConns)

	if err := sqlDB.Ping(); err != nil {
		return nil, errors.Wrap(err, "NewDB")
	}

	return db, nil
}

func DumpRequestCtx(logger *logrus.Logger, ctx *fasthttp.RequestCtx) {
	DumpRequest(logger, ctx)
	DumpResponse(logger, ctx)
}

var dumpRequestMutex sync.Mutex

func DumpRequest(logger *logrus.Logger, ctx *fasthttp.RequestCtx) {
	dumpRequestMutex.Lock()
	defer dumpRequestMutex.Unlock()

	logger.Traceln()
	logger.Traceln("================================")
	logger.Debugf("REQUEST(%d) [%s] %s %s", GetRequestID(ctx), ctx.RemoteAddr().String(), string(ctx.Request.Header.Method()), string(ctx.Request.URI().FullURI()))
	ctx.Request.Header.VisitAll(func(k, v []byte) {
		logger.Tracef("HEADER: %s = %q", string(k), string(v))
	})
	if len(ctx.Request.Body()) > 0 {
		logger.Traceln(string(ctx.Request.Body()))
	}
}

var dumpResponseMutex sync.Mutex

func DumpResponse(logger *logrus.Logger, ctx *fasthttp.RequestCtx) {
	dumpResponseMutex.Lock()
	defer dumpResponseMutex.Unlock()

	logger.Traceln()
	logger.Traceln("---------------------------------")
	logger.Debugf("RESPONSE(%d) [%d]", GetRequestID(ctx), ctx.Response.StatusCode())
	ctx.Response.Header.VisitAll(func(k, v []byte) {
		logger.Tracef("HEADER: %s = %q", string(k), string(v))
	})
	if len(ctx.Response.Body()) > 0 {
		logger.Traceln(string(ctx.Response.Body()))
	}
}

func SetRequestID(ctx *fasthttp.RequestCtx, id uint64) {
	ctx.SetUserValue("request_id", id)
}
func GetRequestID(ctx *fasthttp.RequestCtx) uint64 {
	if v, ok := ctx.UserValue("request_id").(uint64); ok {
		return v
	}
	return 0
}
