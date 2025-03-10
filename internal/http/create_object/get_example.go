package create_object

import (
	"net/http"

	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
	"touchon-server/internal/model"
	"touchon-server/internal/objects"
	"touchon-server/lib/helpers"
	_ "touchon-server/lib/http/server"
)

// Возвращает пример тела запроса для создания объекта
// @Summary Возвращает пример тела запроса для создания объекта
// @Tags Service
// @Description Возвращает пример тела запроса для создания объекта
// @ID GetObjectExample
// @Produce json
// @Param category query string true "Категория объекта" Enums(controller,sensor,regulator,generic_input,relay,rs485,conditioner)
// @Param type query string true "Тип объекта" example(mega_d,htu21d,regulator,generic_input,relay,wb_mrm2_mini,onokom/hr_1_mb_b,bus)
// @Success      200 {object} server.Response[Request]
// @Failure      400 {object} server.Response[any]
// @Failure      500 {object} server.Response[any]
// @Router /_/objects/example [get]
func GetExample(ctx *fasthttp.RequestCtx) (_ interface{}, _ int, e error) {
	objCat := helpers.GetParam(ctx, "category")
	if objCat == "" {
		return nil, http.StatusBadRequest, errors.New("category is empty")
	}

	objType := helpers.GetParam(ctx, "type")
	if objType == "" {
		return nil, http.StatusBadRequest, errors.New("type is empty")
	}

	obj, err := objects.GetObjectModel(objCat, objType, false)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	resp := &Request{
		Object: struct {
			ParentID *int                   `json:"parent_id,omitempty"`
			ZoneID   *int                   `json:"zone_id,omitempty"`
			Category model.Category         `json:"category"`
			Type     string                 `json:"type"`
			Name     string                 `json:"name"`
			Props    map[string]interface{} `json:"props,omitempty"`
			Enabled  bool                   `json:"enabled"`
			Children []*Child               `json:"children,omitempty"`
		}{
			Category: obj.GetCategory(),
			Type:     obj.GetType(),
			Name:     obj.GetName(),
			Props:    make(map[string]interface{}, obj.GetProps().Len()),
			Enabled:  obj.GetEnabled(),
		},
	}

	for _, p := range obj.GetProps().GetAll().GetValueList() {
		switch {
		case p.GetValue() != nil:
			resp.Object.Props[p.Code] = p.GetValue()
		case p.DefaultValue != nil:
			resp.Object.Props[p.Code] = p.DefaultValue
		}
	}

	resp.Object.Children = getChildProps(obj)

	return resp, http.StatusOK, nil
}

func getChildProps(obj objects.Object) []*Child {
	r := make([]*Child, 0, obj.GetChildren().Len())

	for _, child := range obj.GetChildren().GetAll() {
		c := &Child{Props: map[string]interface{}{}}

		for _, p := range child.GetProps().GetAll().GetValueList() {
			switch {
			case p.GetValue() != nil:
				c.Props[p.Code] = p.GetValue()
			case p.DefaultValue != nil:
				c.Props[p.Code] = p.DefaultValue
			}
		}

		if child.GetChildren().Len() > 0 {
			c.Children = getChildProps(child)
		}

		r = append(r, c)
	}

	if len(r) > 0 {
		return r
	}

	return nil
}
