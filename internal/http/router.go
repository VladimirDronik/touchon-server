package http

import (
	"net/http"

	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
	"touchon-server/internal/g"
	"touchon-server/internal/helpers"
	"touchon-server/internal/object/PortMegaD"
	"touchon-server/internal/objects"
	"touchon-server/internal/store"
	"touchon-server/lib/events"
	"touchon-server/lib/events/object/controller"
	"touchon-server/lib/interfaces"
)

// Получение команды от контроллера megaD
// @Summary Получение команды от контроллера megaD
// @Tags MegaD
// @Description Получение команды от контроллера megaD
// @ID Mega
// @Produce json
// @Param mdid  query string true  "Идентификатор контроллера (Настраивается здесь http://<controller_addr>/sec/?cf=2)" default(dev4)
// @Param st    query string false "Признак запуска контроллера" Enums(1)
// @Param pt    query string false "Номер сработавшего порта" default(0)
// @Param ext   query string false "Номер порта модуля"
// @Param click query string false "Клик одинарный/двойной" Enums(1,2)
// @Param m     query string false "При удержании передается 2, при отпускании 1" Enums(1,2)
// @Param v     query string false "Отправляется при срабатывании порта OUT"
// @Success      200 {object} http.Response[any]
// @Failure      400 {object} http.Response[any]
// @Failure      500 {object} http.Response[any]
// @Router /mega [get]
func (o *Server) handleGetMegaD(ctx *fasthttp.RequestCtx) (interface{}, int, error) {
	controllerID := helpers.GetParam(ctx, "mdid") // dev4

	// https://ab-log.ru/smart-house/ethernet/megad-2561#conf-cron
	// контроллер в момент своего включения отправляет на сервер сообщение с параметром "st=1"
	// /mega?st=1&mdid=dev4
	controllerStarted := helpers.GetParam(ctx, "st") // 1

	portNumber := helpers.GetParam(ctx, "pt") // номер сработавшего порта

	// https://ab-log.ru/smart-house/ethernet/megad-2561#conf-exp-pca
	extPortNumber := helpers.GetParam(ctx, "ext") // номер порта модуля

	clickCount := helpers.GetParam(ctx, "click") // одинарный (1) или двойной (2) клик
	holdRelease := helpers.GetParam(ctx, "m")    // при удержании передается 2, при отпускании 1
	value := helpers.GetParam(ctx, "v")          // Отправляется при срабатывании порта OUT
	countImpulse := helpers.GetParam(ctx, "c")   // Количество срабатываний порта IN

	var allMsgs []interfaces.Message

	switch {
	case controllerStarted == "1":
		obj, err := store.I.DeviceRepository().GetControllerByName(controllerID)
		if err != nil {
			return nil, http.StatusInternalServerError, err
		}

		msg, err := controller.NewOnLoad(obj.ID)
		if err != nil {
			return nil, http.StatusInternalServerError, err
		}

		allMsgs = append(allMsgs, msg)

	case portNumber != "":
		objectID, err := store.I.PortRepository().GetPortObjectID(controllerID, portNumber)
		if err != nil {
			return nil, http.StatusInternalServerError, err
		}

		obj, err := objects.LoadPort(objectID, false)
		if err != nil {
			return nil, http.StatusInternalServerError, err
		}

		port, ok := obj.(*PortMegaD.PortModel)
		if !ok {
			err := errors.New("MakeModel returns not PortModel")
			return nil, http.StatusInternalServerError, err
		}

		msgs, err := port.ResCommand(controllerID, portNumber, extPortNumber, clickCount, holdRelease, value, countImpulse)
		if err != nil {
			err = errors.Wrap(err, "ResCommand")
			g.Logger.Warn(err)

			msg, err := events.NewOnError(interfaces.TargetTypeObject, objectID, err.Error())
			if err != nil {
				return nil, http.StatusInternalServerError, err
			}

			allMsgs = append(allMsgs, msg)
		}

		allMsgs = append(allMsgs, msgs...)
	}

	go func() {
		for _, msg := range allMsgs {
			// TODO это костыль - надо переделать методы порта, возвращающие событие,
			// TODO на возврат событие+ошибка
			if msg == nil {
				continue
			}

			if err := g.Msgs.Send(msg); err != nil {
				o.GetLogger().Errorf("handleGetMegaD: %v", err)
			}
		}
	}()

	return nil, http.StatusOK, nil
}
