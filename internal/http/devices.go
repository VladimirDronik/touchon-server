package http

import (
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
	"touchon-server/internal/helpers"
	"touchon-server/internal/model"
	"touchon-server/internal/object/MegaD"
	"touchon-server/internal/store"
)

type GetControllerPortsResponseItem struct {
	Group string  `json:"group"`
	Ports []*Port `json:"ports"`
}

type Port struct {
	ObjectID int      `json:"object_id"`
	Number   int      `json:"number"`
	Type     string   `json:"type"`
	Mode     string   `json:"mode"`
	Objects  []string `json:"objects,omitempty"`
}

// Получение портов контроллера
// @Summary Получение портов контроллера
// @Tags Devices
// @Description Получение портов контроллера
// @ID GetControllerPorts
// @Produce json
// @Param id path int true "ID объекта" default(2)
// @Param group query string false "Группа порта" example(inputs,digital)
// @Param type query string false "Тип порта"
// @Success      200 {object} http.Response[[]GetControllerPortsResponseItem]
// @Failure      400 {object} http.Response[any]
// @Failure      500 {object} http.Response[any]
// @Router /controllers/{id}/ports [get]
func (o *Server) handleGetControllerPorts(ctx *fasthttp.RequestCtx) (interface{}, int, error) {
	objectID, err := helpers.GetUintPathParam(ctx, "id")
	if err != nil {
		return nil, http.StatusBadRequest, err
	}
	typePort := helpers.GetParam(ctx, "type")
	groupPort := strings.Split(strings.TrimSpace(helpers.GetParam(ctx, "group")), ",")

	children, err := store.I.ObjectRepository().GetObjectChildren(objectID)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	childrenIDs := make([]int, 0, len(children))
	for _, row := range children {
		childrenIDs = append(childrenIDs, row.ID)
	}

	allProps, err := store.I.ObjectRepository().GetPropsByObjectIDs(childrenIDs)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	// Обрабатываем порты
	ports := make(map[int]*Port, 50)
	groups := make(map[string]*GetControllerPortsResponseItem, 3)

	for _, child := range children {
		if child.Category != model.CategoryPort {
			continue
		}

		port := &Port{}

		props := allProps[child.ID]

		port.ObjectID = child.ID

		if p, ok := props["number"]; ok {
			port.Number, _ = strconv.Atoi(p.Value)
		}

		if p, ok := props["type"]; ok {
			t, err := MegaD.PortTypes.Get(p.Value)
			if err != nil {
				return nil, http.StatusBadRequest, err
			}

			port.Type = t.Name

			if p, ok := props["mode"]; ok {
				m, err := t.Modes.Get(p.Value)
				if err != nil {
					return nil, http.StatusBadRequest, err
				}

				port.Mode = m.Name
			}
		}

		ports[child.ID] = port

		g, ok := props["group"]
		if !ok {
			return nil, http.StatusInternalServerError, errors.Errorf("prop group not found for object %d", child.ID)
		}

		group, ok := groups[g.Value]
		if !ok {
			group = &GetControllerPortsResponseItem{Group: g.Value}
			groups[g.Value] = group
		}

		for _, gp := range groupPort {
			if (typePort == port.Type || typePort == "") && (gp == group.Group || gp == "") {
				group.Ports = append(group.Ports, port)
			}
		}

	}

	// Добавляем к портам имена привязанных объектов
	for _, child := range children {
		if child.Category == model.CategoryPort {
			continue
		}

		props := allProps[child.ID]
		addr := props["address"]

		var items []string
		if strings.Contains(addr.Value, ";") {
			items = strings.Split(addr.Value, ";")
		} else {
			items = append(items, addr.Value)
		}

		for _, item := range items {
			id, err := strconv.Atoi(strings.TrimSpace(item))
			if err != nil {
				return nil, http.StatusBadRequest, err
			}

			p, ok := ports[id]
			if !ok {
				return nil, http.StatusBadRequest, errors.Errorf("port %d not found", id)
			}

			p.Objects = append(p.Objects, child.Name)
		}
	}

	for _, group := range groups {
		sort.Slice(group.Ports, func(i, j int) bool {
			return group.Ports[i].Number < group.Ports[j].Number
		})
	}

	r := make([]*GetControllerPortsResponseItem, 0, 3)
	for _, item := range MegaD.Groups.GetKeyValueList() {
		if group, ok := groups[item.Key]; ok && group.Ports != nil {
			r = append(r, groups[item.Key])
		}
	}

	return r, http.StatusOK, nil
}
