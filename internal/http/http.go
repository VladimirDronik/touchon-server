package http

import (
	"fmt"

	"github.com/pkg/errors"
	_ "touchon-server/docs"
	"touchon-server/internal/context"
	"touchon-server/internal/http/create_object"
	"touchon-server/internal/http/send_to_mqtt"
	"touchon-server/internal/http/update_object"
	"touchon-server/internal/scripts"
	"touchon-server/internal/store"
	memStore "touchon-server/internal/store/memstore"
	"touchon-server/lib/http/server"
)

// Global instance
var I *Server

func New(ringBuffer fmt.Stringer) (*Server, error) {
	switch {
	case context.Config == nil:
		return nil, errors.Wrap(errors.New("cfg is nil"), "http.New")
	case scripts.I == nil:
		return nil, errors.Wrap(errors.New("scripts is nil"), "http.New")
	case memStore.I == nil:
		return nil, errors.Wrap(errors.New("memStore is nil"), "http.New")
	case context.Logger == nil:
		return nil, errors.Wrap(errors.New("Logger is nil"), "http.New")
	case store.I == nil:
		return nil, errors.Wrap(errors.New("Store is nil"), "http.New")
	}

	baseServer, err := server.New("API", context.Config, ringBuffer, context.Logger)
	if err != nil {
		return nil, errors.Wrap(err, "http.New")
	}

	o := &Server{Server: baseServer}

	// Служебные эндпоинты
	o.AddHandler("GET", "/_/sensors", o.handleGetSensors)               // получение значение датчиков
	o.AddHandler("GET", "/_/objects/example", create_object.GetExample) // получение примера json'а для создания объекта

	o.AddHandler("GET", "/mega", o.handleGetMegaD)                             // прием команд с megad
	o.AddHandler("GET", "/controllers/{id}/ports", o.handleGetControllerPorts) // получение портов контроллера

	o.AddHandler("GET", "/objects/types", o.handleGetObjectsTypes)    // получение категорий и типов объектов
	o.AddHandler("GET", "/objects/model", o.handleGetObjectModel)     // получение модели объекта
	o.AddHandler("GET", "/objects", o.handleGetObjects)               // получение объектов
	o.AddHandler("GET", "/objects/{id}", o.handleGetObject)           // получение объекта
	o.AddHandler("POST", "/objects", create_object.Handler)           // добавление объекта с методами
	o.AddHandler("PUT", "/objects", update_object.Handler)            // обновление объекта
	o.AddHandler("DELETE", "/objects/{id}", o.handleDeleteObject)     // удаление объекта
	o.AddHandler("GET", "/objects/tags", o.handleGetAllObjectsTags)   // получение всех тегов
	o.AddHandler("GET", "/objects/by_tags", o.handleGetObjectsByTags) // получение объектов по тегам

	o.AddHandler("GET", "/scripts/model", o.handleGetScriptModel)
	o.AddHandler("GET", "/scripts", o.handleGetScripts)
	o.AddHandler("GET", "/scripts/{id}", o.handleGetScript)
	o.AddHandler("POST", "/scripts", o.handleCreateScript)
	o.AddHandler("PUT", "/scripts", o.handleUpdateScript)
	o.AddHandler("DELETE", "/scripts/{id}", o.handleDeleteScript)
	o.AddHandler("POST", "/scripts/{id}/exec", o.handleExecScript)

	// Метод для тестирования
	o.AddHandler("POST", "/_/mqtt", send_to_mqtt.Handler)

	//o.AddHandler("POST", "/wizard/create_object", o.handleWizardCreateObject)

	return o, nil
}

type Server struct {
	*server.Server
}

//func (o *Server) Start() error {
//	if err := o.Server.Start(o.GetConfig()["bind_addr"]); err != nil {
//		return errors.Wrap(err, "Start")
//	}
//
//	return nil
//}

//func (o *Server) Shutdown() error {
//	if err := o.Server.Shutdown(); err != nil {
//		return errors.Wrap(err, "Shutdown")
//	}
//
//	return nil
//}
