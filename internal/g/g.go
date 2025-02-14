package g

import (
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
	"touchon-server/internal/model"
	"touchon-server/lib/interfaces"
)

var (
	Config     map[string]string
	Logger     *logrus.Logger
	Msgs       interfaces.MessagesService
	HttpServer interfaces.HttpServer
)

func GetAccessLevel(ctx *fasthttp.RequestCtx) (model.AccessLevel, error) {
	// TODO get access level for user
	return model.AccessLevelAllowed, nil
}
