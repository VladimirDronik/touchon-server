package http

import (
	"encoding/json"
	"net/http"

	"github.com/VladimirDronik/touchon-server/helpers"
	"github.com/valyala/fasthttp"
	"translator/internal/model"
)

type getControlPanelResponse struct {
	ScenarioItems []*model.Scenario  `json:"scenario_items"`
	ZoneItems     []*model.GroupRoom `json:"room_items"`
}

// Получение элементов панели управления
// @Security TokenAuth
// @Summary Получение элементов панели управления
// @Tags Pages
// @Description Получение элементов панели управления
// @ID GetControlPanel
// @Produce json
// @Success      200 {object} Response[getControlPanelResponse]
// @Failure      400 {object} Response[any]
// @Failure      500 {object} Response[any]
// @Router /private/cp [get]
func (o *Server) getControlPanel(ctx *fasthttp.RequestCtx) (interface{}, int, error) {
	zoneItems, err := o.outputZoneItems()
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	scenarios, err := o.store.Items().GetScenarios()
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	r := getControlPanelResponse{
		ScenarioItems: scenarios,
		ZoneItems:     zoneItems,
	}

	// Backward compatibility
	for _, item := range r.ScenarioItems {
		if item.Enabled {
			item.Status = "on"
		} else {
			item.Status = "off"
		}
	}

	return r, http.StatusOK, nil
}

// Получение дашборда для приложения
// @Security TokenAuth
// @Summary Получение дашборда для приложения
// @Tags Pages
// @Description Получение дашборда для приложения
// @ID GetDashboard
// @Produce json
// @Success      200 {object} Response[any]
// @Failure      400 {object} Response[any]
// @Failure      500 {object} Response[any]
// @Router /private/dashboard [get]
func (o *Server) getDashboard(ctx *fasthttp.RequestCtx) (interface{}, int, error) {
	return nil, http.StatusNotFound, nil // (временно возвращаем 404)
}

// Создание помещения
// @Security TokenAuth
// @Summary Создание помещения
// @Tags Zones
// @Description Создание помещения
// @ID CreateRoom
// @Accept json
// @Produce json
// @Param room body model.Zone true "Помещение"
// @Success      200 {object} Response[model.Zone]
// @Failure      400 {object} Response[any]
// @Failure      500 {object} Response[any]
// @Router /private/room [post]
func (o *Server) handleCreateZone(ctx *fasthttp.RequestCtx) (interface{}, int, error) {
	zone := &model.Zone{}

	err := json.Unmarshal(ctx.Request.Body(), &zone)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	zoneID, err := o.store.Zones().CreateZone(zone)

	return zoneID, http.StatusOK, err
}

// Получение зон, в которых есть элементы
// @Security TokenAuth
// @Summary Получение зон, в которых есть элементы
// @Tags Zones
// @Description Получение зон, в которых есть элементы
// @ID GetZones
// @Produce json
// @Success      200 {object} Response[[]model.Zone]
// @Failure      400 {object} Response[any]
// @Failure      500 {object} Response[any]
// @Router /private/rooms-list [get]
func (o *Server) getZones(ctx *fasthttp.RequestCtx) (interface{}, int, error) {
	zones, err := o.store.Items().GetZones()
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	// Backward compatibility
	setIsGroup(zones)

	return zones, http.StatusOK, nil
}

// Backward compatibility
func setIsGroup(items []*model.Zone) {
	for _, item := range items {
		item.IsGroup = len(item.Children) > 0

		if len(item.Children) > 0 {
			setIsGroup(item.Children)
		}
	}
}

// Получение всех зон
// @Security TokenAuth
// @Summary Получение всех зон
// @Tags Zones
// @Description Получение списка всех помещений, независимо от того есть там итемы или нет
// @ID GetAllZones
// @Produce json
// @Success      200 {object} Response[[]model.Zone]
// @Failure      400 {object} Response[any]
// @Failure      500 {object} Response[any]
// @Router /private/rooms-list-all [get]
func (o *Server) getAllZones(ctx *fasthttp.RequestCtx) (interface{}, int, error) {
	zones, err := o.store.Zones().GetZoneTrees(0)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	// Backward compatibility
	setIsGroup(zones)

	return zones, http.StatusOK, nil
}

// Получение зоны
// @Security TokenAuth
// @Summary Получение зоны
// @Tags Zones
// @Description Получение зоны
// @ID GetZone
// @Produce json
// @Param id query int true "ID" default(1)
// @Success      200 {object} Response[model.GroupRoom]
// @Failure      400 {object} Response[any]
// @Failure      500 {object} Response[any]
// @Router /private/room [get]
func (o *Server) getZone(ctx *fasthttp.RequestCtx) (interface{}, int, error) {
	id, err := helpers.GetUintParam(ctx, "id")
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	zone, err := o.outputZoneItem(id)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return zone, http.StatusOK, nil
}

// Обновляет зоны рекурсивно или одно помещение только
// @Security TokenAuth
// @Summary Обновляет зоны рекурсивно или одно помещение только
// @Tags Zones
// @Description Обновляет зоны рекурсивно или одно помещение только
// @ID UpdateZones
// @Produce json
// @Param body body []model.Zone true "Body"
// @Success      200 {object} Response[any]
// @Failure      400 {object} Response[any]
// @Failure      500 {object} Response[any]
// @Router /private/rooms-list-all [patch]
func (o *Server) updateZones(ctx *fasthttp.RequestCtx) (interface{}, int, error) {
	var zones []*model.Zone

	if err := json.Unmarshal(ctx.Request.Body(), &zones); err != nil {
		return nil, http.StatusInternalServerError, err
	}

	if err := o.store.Zones().UpdateZones(zones); err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return nil, http.StatusOK, nil
}

// Установка порядка отображения зон
// @Security TokenAuth
// @Summary Установка порядка отображения зон
// @Tags Zones
// @Description Установка порядка отображения зон
// @ID SetZonesOrder
// @Produce json
// @Param body body []int true "Body"
// @Success      200 {object} Response[any]
// @Failure      400 {object} Response[any]
// @Failure      500 {object} Response[any]
// @Router /private/zones/order [patch]
func (o *Server) setZonesOrder(ctx *fasthttp.RequestCtx) (interface{}, int, error) {
	var zoneIDs []int

	if err := json.Unmarshal(ctx.Request.Body(), &zoneIDs); err != nil {
		return nil, http.StatusInternalServerError, err
	}

	if err := o.store.Zones().SetOrder(zoneIDs); err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return nil, http.StatusOK, nil
}

// Удаление помещения
// @Security TokenAuth
// @Summary Удаление помещения
// @Tags Zones
// @Description Удаление помещения
// @ID DeleteZone
// @Produce json
// @Param id query int true "ID"
// @Success      200 {object} Response[any]
// @Failure      400 {object} Response[any]
// @Failure      500 {object} Response[any]
// @Router /private/room [delete]
func (o *Server) handleDeleteZone(ctx *fasthttp.RequestCtx) (interface{}, int, error) {
	id, err := helpers.GetUintParam(ctx, "id")
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	err = o.store.Zones().DeleteZone(id)

	return nil, http.StatusOK, err
}
