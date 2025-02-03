package mdns

import (
	"net"
	"strconv"

	"github.com/grandcat/zeroconf"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func NewServer(instance, service, domain, connIface, serverID string, apiPort, wsPort, cctvPort, webPort int, logger *logrus.Logger) (*Server, error) {
	switch {
	case apiPort < 1 || 65535 < apiPort:
		return nil, errors.Wrap(errors.New("apiPort is bad"), "mdns.NewServer")
	case wsPort < 1 || 65535 < wsPort:
		return nil, errors.Wrap(errors.New("wsPort is bad"), "mdns.NewServer")
	case cctvPort < 1 || 65535 < cctvPort:
		return nil, errors.Wrap(errors.New("cctvPort is bad"), "mdns.NewServer")
	case webPort < 1 || 65535 < webPort:
		return nil, errors.Wrap(errors.New("webPort is bad"), "mdns.NewServer")
	}

	o := &Server{
		instance: instance,
		service:  service,
		domain:   domain,
		webPort:  webPort,
		txtRecords: []string{
			serverID,
			"apiPort=" + strconv.Itoa(apiPort),
			"wsPort=" + strconv.Itoa(wsPort),
			"cctvPort=" + strconv.Itoa(cctvPort),
		},
		logger: logger,
	}

	items, err := net.Interfaces()
	if err != nil {
		return nil, errors.Wrap(err, "mdns.NewServer")
	}

	for _, item := range items {
		if item.Name == connIface {
			multicastAddrs, err := item.MulticastAddrs()
			if err != nil {
				o.logger.Warnf("Ошибка получения мультикастовых адресов для %s: %v", item.Name, err)
				continue
			}

			for _, addr := range multicastAddrs {
				// Проверяем, является ли адрес IPv4
				if ipAddr, ok := addr.(*net.IPAddr); ok && ipAddr.IP.To4() != nil {
					o.interfaces = append(o.interfaces, item)
					break
				}
			}
		}
	}

	return o, nil
}

type Server struct {
	instance   string
	service    string
	domain     string
	webPort    int
	interfaces []net.Interface
	txtRecords []string
	logger     *logrus.Logger
	server     *zeroconf.Server
}

// Start начать трансляцию mDNS сервиса TouchON
func (o *Server) Start() error {
	var err error
	o.server, err = zeroconf.Register(o.instance, o.service, o.domain, o.webPort, o.txtRecords, o.interfaces)
	if err != nil {
		return errors.Wrap(err, "Start")
	}

	return nil
}

func (o *Server) Shutdown() error {
	o.server.Shutdown()
	return nil
}
