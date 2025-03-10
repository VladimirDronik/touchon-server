package helpers

import (
	"strconv"
	"strings"
	"sync"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

// GetParam отдает любой параметр, который был получен из GET запроса
func GetParam(ctx *fasthttp.RequestCtx, paramName string) string {
	return strings.TrimSpace(string(ctx.QueryArgs().Peek(paramName)))
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

func GetBoolParam(ctx *fasthttp.RequestCtx, paramName string, defaultValue bool) (bool, error) {
	if s := GetParam(ctx, paramName); s != "" {
		v, err := strconv.ParseBool(s)
		if err != nil {
			return false, errors.Wrapf(err, "GetBoolParam(%s)", paramName)
		}

		return v, nil
	}

	return defaultValue, nil
}

func GetSliceParam(ctx *fasthttp.RequestCtx, paramName string) []string {
	v := GetParam(ctx, paramName)
	slice := strings.Split(v, ",")

	items := make([]string, 0, len(slice))
	for _, item := range slice {
		item = strings.TrimSpace(item)
		if item != "" {
			items = append(items, item)
		}
	}

	return items
}

func GetMapParam(ctx *fasthttp.RequestCtx, paramName string) (map[string]string, error) {
	s := GetSliceParam(ctx, paramName)

	m := make(map[string]string, len(s))
	for _, item := range s {
		kv := strings.Split(item, "=")
		if len(kv) != 2 {
			return nil, errors.Wrap(errors.New("bad map param"), "GetMapParam")
		}
		for i, item := range kv {
			kv[i] = strings.TrimSpace(item)
		}

		k, v := kv[0], kv[1]

		if k == "" {
			return nil, errors.Wrap(errors.New("bad map param"), "GetMapParam")
		}

		m[k] = v
	}

	return m, nil
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
