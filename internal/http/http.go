package http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
	_ "touchon-server/docs"
	"touchon-server/internal/g"
	"touchon-server/internal/http/create_object"
	"touchon-server/internal/http/update_object"
	"touchon-server/internal/model"
	"touchon-server/internal/scripts"
	"touchon-server/internal/store"
	memStore "touchon-server/internal/store/memstore"
	"touchon-server/internal/token"
	"touchon-server/lib/http/server"
)

var (
	errTokenNotFound      = errors.New("can not find token in header")
	errParamsError        = errors.New("params error")
	errServerKeyIncorrect = errors.New("server key incorrect")
	errNoResult           = errors.New("no result")
)

func New(ringBuffer fmt.Stringer) (*Server, error) {
	switch {
	case g.Config == nil:
		return nil, errors.Wrap(errors.New("cfg is nil"), "http.New")
	case scripts.I == nil:
		return nil, errors.Wrap(errors.New("scripts is nil"), "http.New")
	case memStore.I == nil:
		return nil, errors.Wrap(errors.New("memStore is nil"), "http.New")
	case g.Logger == nil:
		return nil, errors.Wrap(errors.New("Logger is nil"), "http.New")
	case store.I == nil:
		return nil, errors.Wrap(errors.New("Store is nil"), "http.New")
	}

	baseServer, err := server.New("API", g.Config, g.Logger)
	if err != nil {
		return nil, errors.Wrap(err, "http.New")
	}

	o := &Server{
		Server:     baseServer,
		ringBuffer: ringBuffer,
	}

	o.AddHandler("GET", "/token", o.getToken)      // Запрос токена
	o.AddHandler("POST", "/token", o.refreshToken) // Рефреш

	// Вход в зону private
	private := o.addMiddleware("/private", o.authMiddleware)

	// Служебные эндпоинты
	svc := o.addMiddleware("/_", o.authMiddleware)
	svc("GET", "/info", o.handleGetInfo)
	svc("GET", "/sensors", o.handleGetSensors)               // получение значение датчиков
	svc("GET", "/objects/example", create_object.GetExample) // получение примера json'а для создания объекта

	rawSvc := o.addRawMiddleware("/_", o.authMiddleware)
	rawSvc("GET", "/log", o.handleGetLog)

	o.AddHandler("GET", "/mega", o.handleGetMegaD) // прием команд с megad

	ctrl := o.addMiddleware("/controllers", o.authMiddleware)
	ctrl("GET", "/{id}/ports", o.handleGetControllerPorts) // получение портов контроллера

	objects := o.addMiddleware("/objects", o.authMiddleware)
	objects("GET", "/types", o.handleGetObjectsTypes)          // получение категорий и типов объектов
	objects("GET", "/model", o.handleGetObjectModel)           // получение модели объекта
	objects("GET", "/", o.handleGetObjects)                    // получение объектов
	objects("GET", "/{id}", o.handleGetObject)                 // получение объекта
	objects("POST", "/", create_object.Handler)                // добавление объекта с методами
	objects("PUT", "/", update_object.Handler)                 // обновление объекта
	objects("DELETE", "/{id}", o.handleDeleteObject)           // удаление объекта
	objects("GET", "/tags", o.handleGetAllObjectsTags)         // получение всех тегов
	objects("GET", "/by_tags", o.handleGetObjectsByTags)       // получение объектов по тегам
	objects("POST", "/{id}/exec/{method}", o.handleExecMethod) // Запуск метода объекта

	scripts := o.addMiddleware("/scripts", o.authMiddleware)
	scripts("GET", "/model", o.handleGetScriptModel)
	scripts("GET", "/", o.handleGetScripts)
	scripts("GET", "/{id}", o.handleGetScript)
	scripts("POST", "/", o.handleCreateScript)
	scripts("PUT", "/", o.handleUpdateScript)
	scripts("DELETE", "/{id}", o.handleDeleteScript)
	scripts("POST", "/{id}/exec", o.handleExecScript)

	// AR

	events := o.addMiddleware("/events", o.authMiddleware)
	events("GET", "/actions/count", o.handleGetEventsActionsCount)  // получение количества действий для событий
	events("GET", "/actions", o.handleGetEventsActions)             // получение действий для событий
	events("POST", "/actions", o.handleCreateEventAction)           // добавления действия для события
	events("PUT", "/actions", o.handleUpdateEventAction)            // обновление действия для события
	events("DELETE", "/actions/{id}", o.handleDeleteEventAction)    // удаление действия для события
	events("DELETE", "/all-actions", o.handleDeleteAllEventActions) // удаление всех событий по фильтру
	events("PUT", "/actions/order", o.handleOrderEventActions)      // смена порядка действий для события
	events("DELETE", "/", o.handleDeleteEvent)                      // удаление события с действиями

	cron := o.addMiddleware("/cron", o.authMiddleware)
	cron("POST", "/task", o.handleCreateTask)   // создание задания крона
	cron("DELETE", "/task", o.handleDeleteTask) // удаления задания крона
	cron("PUT", "/task", o.handleUpdateTask)    // изменение задания крона

	// TR

	// Users
	private("PATCH", "/users/link-token", o.linkDeviceToken)
	private("GET", "/users", o.handleGetAllUsers)
	private("POST", "/users", o.handleCreateUser)
	private("DELETE", "/users", o.handleDeleteUser)

	private("GET", "/checkserver", o.checkServer) // Проверка доступности локального сервера

	// Получение дашборда и панели управления
	private("GET", "/cp", o.getControlPanel)     // Запрос панели управления
	private("GET", "/dashboard", o.getDashboard) // Запрос дашборда

	// Помещения
	private("POST", "/room", o.handleCreateZone)       // Создание помещения
	private("GET", "/rooms-list", o.getZones)          // Запрос помещений, где есть итемы
	private("GET", "/rooms-list-all", o.getAllZones)   // Запрос всех помещений
	private("PATCH", "/rooms-list-all", o.updateZones) // Изменение помещений (одного или рекурсивно)
	private("PATCH", "/zones/order", o.setZonesOrder)  // Изменение сортировки помещения
	private("GET", "/room", o.getZone)                 // Запрос помещения
	private("DELETE", "/room", o.handleDeleteZone)     // Удаление помещения

	// История значений
	private("GET", "/history", o.getObjectHistory)
	private("GET", "/generate-history", o.generateHistory)

	// Свет
	private("GET", "/light", o.getLight)
	private("PATCH", "/light/hsv", o.setLightHSVColor)
	private("PATCH", "/light/cct", o.setLightCCTColor)
	private("PATCH", "/light/brightness", o.setLightBrightness)

	// Шторы
	private("GET", "/curtain", o.getCurtain)
	private("PATCH", "/curtain/open-percent", o.setCurtainOpenPercent)

	// Кондиционер
	private("GET", "/conditioner", o.getConditioner)
	private("PATCH", "/conditioner/temp", o.setConditionerTemperature)
	private("PATCH", "/conditioner/mode", o.setConditionerMode)
	private("PATCH", "/conditioner/operating-mode", o.setConditionerOperatingMode)
	private("PATCH", "/conditioner/fan-speed", o.setConditionerFanSpeed)
	private("PATCH", "/conditioner/air-direction", o.setConditionerDirection)
	private("PATCH", "/conditioner/extra-mode", o.setConditionerExtraMode)

	// Котел
	private("GET", "/boiler", o.getBoiler)
	private("PATCH", "/boiler/outline-status", o.setBoilerOutlineStatus)
	private("PATCH", "/boiler/heating-mode", o.setBoilerHeatingMode)
	private("PATCH", "/boiler/heating-temp", o.setBoilerHeatingTemperature)
	private("PUT", "/boiler/presets", o.updateBoilerPresets)

	// Счетчики
	private("GET", "/counters-list", o.getCountersList)
	private("GET", "/counter", o.getCounter)

	// Уведомления
	private("GET", "/notifications/unread-count", o.getUnreadNotificationsCount)
	private("GET", "/notifications", o.getNotifications)
	private("PATCH", "/notification", o.setNotificationIsRead)

	// Items
	// Получение элементов, которые требуют дополнительных данных
	//private("GET", "/item/dimer", o.getDimmer)          // Получение димера, убрать
	//private("GET", "/item/thermostat", o.getThermostat) // Получение термостата, убарть

	private("POST", "/items", o.handleCreateItem)  // Создание элемента
	private("PUT", "/items", o.handleUpdateItem)   // Обновление элемента
	private("PATCH", "/item", o.updateItem)        // Изменение итема
	private("GET", "/item", o.getItem)             // Получение всей информации об итеме
	private("DELETE", "/item", o.handleDeleteItem) // Удаление итема
	private("POST", "/item-change", o.itemChange)
	private("PATCH", "/items/order", o.setItemsOrder) // Изменение порядка элементов

	// Итемы датчиков
	private("GET", "/item/sensor", o.getSensor)           //получение данных датчика
	private("POST", "/item/sensor", o.handleCreateSensor) // Создание датчика
	private("PATCH", "/item/sensor", o.handleUpdateSensor)
	private("DELETE", "/item/sensor", o.handleDeleteSensor)
	private("PATCH", "/item/sensor/value", o.handleSetTargetSensor) //Установка значение target для датчика

	// Загрузка меню
	private("GET", "/menu", o.getMenu)

	// Мастера
	private("POST", "/wizard/create_item", o.handleWizardCreateItem)

	return o, nil
}

type Server struct {
	*server.Server
	ringBuffer fmt.Stringer
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

func (o *Server) createSession(deviceID int) (*model.Tokens, error) {
	tokenJWT := token.New(o.GetConfig()["token_secret"])

	accessTokenTTL, err := time.ParseDuration(o.GetConfig()["access_token_ttl"])
	if err != nil {
		return nil, errors.Wrap(err, "createSession")
	}

	refreshTokenTTL, err := time.ParseDuration(o.GetConfig()["refresh_token_ttl"])
	if err != nil {
		return nil, errors.Wrap(err, "createSession")
	}

	accessToken, err := tokenJWT.NewJWT(deviceID, accessTokenTTL)
	if err != nil {
		return nil, errors.Wrap(err, "createSession")
	}

	refreshToken, err := tokenJWT.NewRefreshToken()
	if err != nil {
		return nil, errors.Wrap(err, "createSession")
	}

	tokens := &model.Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	if err := store.I.Users().AddRefreshToken(deviceID, tokens.RefreshToken, refreshTokenTTL); err != nil {
		return nil, errors.Wrap(err, "createSession")
	}

	return tokens, nil
}

func (o *Server) setCookie(ctx *fasthttp.RequestCtx, name string, value string, httpOnly bool) {
	c := &fasthttp.Cookie{}
	c.SetKey(name)
	c.SetValue(value)
	c.SetHTTPOnly(httpOnly)

	if httpOnly {
		c.SetPath("/token")
	}

	ctx.Response.Header.SetCookie(c)
}

// checkLocal проверка доступности локального сервера, возвращает его настройки в случае успешного выполнения
func (o *Server) checkServer(ctx *fasthttp.RequestCtx) (interface{}, int, error) {
	type Server struct {
		Version string `json:"version"`
	}

	return &Server{Version: o.GetConfig()["version"]}, http.StatusOK, nil
}
