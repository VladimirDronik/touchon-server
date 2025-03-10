package http

import (
	"net/http"

	"github.com/valyala/fasthttp"
	"touchon-server/internal/helpers"
	"touchon-server/internal/store"
)

// Получение пунктов меню
// @Security TokenAuth
// @Summary Получение пунктов меню
// @Tags Menus
// @Description Получение пунктов меню
// @ID GetMenu
// @Produce json
// @Param parent query int true "Parent ID" default(0)
// @Success      200 {object} Response[[]model.Menu]
// @Failure      400 {object} Response[any]
// @Failure      500 {object} Response[any]
// @Router /private/menu [get]
func (o *Server) getMenu(ctx *fasthttp.RequestCtx) (interface{}, int, error) {
	parent, err := helpers.GetUintParam(ctx, "parent")
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	menus, err := store.I.Items().GetMenus(parent)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	return menus, http.StatusOK, nil
}
