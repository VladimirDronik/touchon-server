package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
	_ "touchon-server/docs"
	"touchon-server/internal/context"
	"touchon-server/internal/http/create_object"
	"touchon-server/internal/http/send_to_mqtt"
	"touchon-server/internal/http/update_object"
	"touchon-server/internal/model"
	"touchon-server/internal/scripts"
	"touchon-server/internal/store"
	memStore "touchon-server/internal/store/memstore"
	"touchon-server/internal/token"
	httpClient "touchon-server/lib/http/client"
	"touchon-server/lib/http/server"
)

var (
	errTokenNotFound      = errors.New("can not find token in header")
	errParamsError        = errors.New("params error")
	errServerKeyIncorrect = errors.New("server key incorrect")
	errNoResult           = errors.New("no result")
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

	o := &Server{
		Server: baseServer,
		fasthttpClient: &fasthttp.Client{
			Name:                     "touchon-server",
			ReadTimeout:              5 * time.Second,
			WriteTimeout:             5 * time.Second,
			NoDefaultUserAgentHeader: false,
		},
		httpClient: httpClient.New(),
	}

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

	// AR

	o.AddHandler("GET", "/events/actions/count", o.handleGetEventsActionsCount)  // получение количества действий для событий
	o.AddHandler("GET", "/events/actions", o.handleGetEventsActions)             // получение действий для событий
	o.AddHandler("POST", "/events/actions", o.handleCreateEventAction)           // добавления действия для события
	o.AddHandler("PUT", "/events/actions", o.handleUpdateEventAction)            // обновление действия для события
	o.AddHandler("DELETE", "/events/actions/{id}", o.handleDeleteEventAction)    // удаление действия для события
	o.AddHandler("DELETE", "/events/all-actions", o.handleDeleteAllEventActions) // удаление всех событий по фильтру
	o.AddHandler("PUT", "/events/actions/order", o.handleOrderEventActions)      // смена порядка действий для события
	o.AddHandler("DELETE", "/events", o.handleDeleteEvent)                       // удаление события с действиями

	o.AddHandler("POST", "/cron/task", o.handleCreateTask)   // создание задания крона
	o.AddHandler("DELETE", "/cron/task", o.handleDeleteTask) // удаления задания крона
	o.AddHandler("PUT", "/cron/task", o.handleUpdateTask)    // изменение задания крона

	// TR

	o.AddHandler("GET", "/token", o.getToken)      // Запрос токена
	o.AddHandler("POST", "/token", o.refreshToken) // Рефреш

	// Вход в зону private
	private := o.addMiddleware("/private", o.authMiddleware)

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

	// Получение данных сенсора
	private("GET", "/sensor", o.getSensor)

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
	private("GET", "/item/dimer", o.getDimmer)          // Получение димера
	private("GET", "/item/thermostat", o.getThermostat) // Получение термостата

	private("POST", "/items", o.handleCreateItem)  // Создание элемента
	private("PUT", "/items", o.handleUpdateItem)   // Обновление элемента
	private("PATCH", "/item", o.updateItem)        // Изменение итема
	private("GET", "/item", o.getItem)             // Получение всей информации об итеме
	private("DELETE", "/item", o.handleDeleteItem) // Удаление итема
	private("POST", "/item-change", o.itemChange)
	private("PATCH", "/items/order", o.setItemsOrder) // Изменение порядка элементов

	// Датчики в помещении
	private("POST", "/item/sensor", o.handleCreateSensor) // Создание датчика
	//private("PATCH", "/item/sensor", o.handleUpdateSensor)
	private("DELETE", "/item/sensor", o.handleDeleteSensor)

	// Загрузка меню
	private("GET", "/menu", o.getMenu)

	// Мастера
	private("POST", "/wizard/create_item", o.handleWizardCreateItem)

	// Прокси для сервисов
	proxy := o.addMiddlewareRaw("/proxy", o.authMiddlewareRaw)
	proxy("*", "/{service}/{filepath:*}", o.proxy)

	return o, nil
}

type Server struct {
	*server.Server
	httpClient     *httpClient.Client
	fasthttpClient *fasthttp.Client
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

// Backward compatibility

// Заменяем в ответе поле data на response

// Response Ответ сервиса
type Response[T any] struct {
	Meta    server.Meta `json:"meta"` // Метаинформация о запросе/ответе
	Success bool        `json:"success"`
	Data    T           `json:"response,omitempty"` // Полезная нагрузка, зависит от запроса
	Error   string      `json:"error,omitempty"`    // Описание возвращенной ошибки
}

// JsonHandlerWrapper ответ в формате JSON оборачивает в единый формат и добавляет метаданные.
func JsonHandlerWrapper(f server.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		var r Response[any]
		const magic = "CoNtEnTLeNgTh"

		start := time.Now()
		data, status, err := f(ctx)
		r.Meta.Duration = float64(int(time.Since(start).Seconds()*1000)) / 1000
		r.Meta.ContentLength = magic
		ctx.Response.SetStatusCode(status)

		switch {
		case err != nil:
			r.Error = err.Error()
		case data != nil:
			r.Data = data
		}
		r.Success = err == nil

		var buf bytes.Buffer

		enc := json.NewEncoder(&buf)
		enc.SetIndent("", "  ")
		if err := enc.Encode(r); err != nil {
			buf.Reset()
			r.Data = nil
			r.Error = err.Error()
			_ = enc.Encode(r)
		}

		// Выставляем размер ответа
		contLength := buf.Len() - len(magic)
		body := strings.Replace(buf.String(), magic, strconv.Itoa(contLength/1024)+"K", 1)
		_, _ = ctx.WriteString(body)
	}
}

// Override AddHandler
func (o *Server) AddHandler(method, path string, handler server.RequestHandler) {
	o.GetRouter().Handle(method, path, JsonHandlerWrapper(handler))
}
