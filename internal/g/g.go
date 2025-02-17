package g

import (
	"github.com/sirupsen/logrus"
	"touchon-server/lib/interfaces"
)

var (
	Config     map[string]string
	Logger     *logrus.Logger
	Msgs       interfaces.MessagesService
	HttpServer interfaces.HttpServer
	NodeRed    interfaces.NodeRed
)
