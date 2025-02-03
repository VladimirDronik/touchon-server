package context

import (
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
	"touchon-server/internal/model"
)

var (
	Config map[string]string
	Logger *logrus.Logger
)

func GetAccessLevel(ctx *fasthttp.RequestCtx) (model.AccessLevel, error) {
	// TODO get access level for user
	return model.AccessLevelAllowed, nil
}
