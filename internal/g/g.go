package g

import (
	"github.com/sirupsen/logrus"
	"touchon-server/lib/interfaces"
)

const DisabledAuthDeviceID = 10

var (
	Config     map[string]string
	Logger     *logrus.Logger
	Msgs       interfaces.MessagesService
	HttpServer interfaces.HttpServer
	WSServer   interfaces.WSServer
	NodeRed    interfaces.NodeRed
)
