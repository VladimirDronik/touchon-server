package http

import (
	"net/http"
	"sort"

	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
	"touchon-server/internal/helpers"
	"touchon-server/internal/model"
	"touchon-server/internal/objects"
	memStore "touchon-server/internal/store/memstore"
	libHelpers "touchon-server/lib/helpers"
)

type GetObjectsParams struct {
	FilterByTags     []string
	FilterByID       int
	FilterByParentID int
	FilterByZoneID   int
	FilterByCategory string
	FilterByType     string
	FilterByName     string
	FilterByStatus   string

	AllTypes    bool
	Generations int
	TypeStruct  string
	WithMethods bool
	WithTags    bool
	WithParents bool

	Offset int
	Limit  int
}

func parseGetObjectsParams(ctx *fasthttp.RequestCtx) (*GetObjectsParams, error) {
	r := &GetObjectsParams{}
	var err error

	r.FilterByTags = libHelpers.PrepareTags(helpers.GetParam(ctx, "tags"), ",")

	r.FilterByID, err = helpers.GetUintParam(ctx, "filter_by_id")
	if err != nil {
		return nil, errors.Wrap(err, "parseGetObjectsParams")
	}

	r.FilterByParentID, err = helpers.GetUintParam(ctx, "filter_by_parent_id")
	if err != nil {
		return nil, errors.Wrap(err, "parseGetObjectsParams")
	}

	r.FilterByZoneID, err = helpers.GetUintParam(ctx, "filter_by_zone_id")
	if err != nil {
		return nil, errors.Wrap(err, "parseGetObjectsParams")
	}

	r.FilterByCategory = helpers.GetParam(ctx, "filter_by_category")
	if r.FilterByCategory != "" {
		if _, ok := objects.GetCategoriesAndTypes()[r.FilterByCategory]; !ok || r.FilterByCategory == model.CategoryPort || r.FilterByCategory == model.CategorySensorValue {
			return nil, errors.Wrap(errors.Errorf("bad category value"), "parseGetObjectsParams")
		}
	}

	r.FilterByType = helpers.GetParam(ctx, "filter_by_type")
	r.FilterByName = helpers.GetParam(ctx, "filter_by_name")
	r.FilterByStatus = helpers.GetParam(ctx, "filter_by_status")

	r.Generations, err = helpers.GetUintParam(ctx, "children")
	if err != nil {
		return nil, errors.Wrap(err, "parseGetObjectsParams")
	}

	r.TypeStruct = helpers.GetParam(ctx, "type_struct")
	typeChildren := helpers.GetParam(ctx, "type_children")
	r.AllTypes = typeChildren == "all" || typeChildren == "internal"

	r.WithMethods, err = helpers.GetBoolParam(ctx, "with_methods", false)
	if err != nil {
		return nil, errors.Wrap(err, "parseGetObjectsParams")
	}

	r.WithTags, err = helpers.GetBoolParam(ctx, "with_tags", true)
	if err != nil {
		return nil, errors.Wrap(err, "parseGetObjectsParams")
	}

	simpleTree, err := helpers.GetBoolParam(ctx, "simple_tree", false)
	if err != nil {
		return nil, errors.Wrap(err, "parseGetObjectsParams")
	}
	r.WithParents = !simpleTree

	r.Offset, err = helpers.GetUintParam(ctx, "offset")
	if err != nil {
		return nil, errors.Wrap(err, "parseGetObjectsParams")
	}

	r.Limit, err = helpers.GetUintParam(ctx, "limit")
	if err != nil {
		return nil, errors.Wrap(err, "parseGetObjectsParams")
	}

	if r.Limit == 0 {
		r.Limit = 20
	}

	return r, nil
}

// Фильтруем объекты
func filterObjects(params *GetObjectsParams) ([]objects.Object, error) {
	items, err := memStore.I.Search(func(items map[int]objects.Object) ([]objects.Object, error) {
		r := make([]objects.Object, 0, 10)

		for _, item := range items {
			switch {
			case len(params.FilterByTags) > 0 && !compareTags(params.FilterByTags, item.GetTags()):
			case params.FilterByID > 0 && item.GetID() != params.FilterByID:
			case params.FilterByParentID > 0 && (item.GetParentID() == nil || *item.GetParentID() != params.FilterByParentID):
			case params.FilterByZoneID > 0 && (item.GetZoneID() == nil || *item.GetZoneID() != params.FilterByZoneID):
			case params.FilterByCategory != "" && item.GetCategory() != params.FilterByCategory:
			case params.FilterByType != "" && item.GetType() != params.FilterByType:
			case params.FilterByName != "" && item.GetName() != params.FilterByName:
			case params.FilterByStatus != "" && item.GetStatus() != params.FilterByStatus:
			case !params.AllTypes && item.GetFlags().Has(objects.Hidden):
			default:
				r = append(r, item)
			}
		}

		sort.Slice(r, func(i, j int) bool {
			return r[i].GetID() < r[j].GetID()
		})

		return r, nil
	})

	if err != nil {
		return nil, errors.Wrap(err, "filterObjects")
	}

	return items, nil
}

type GetObjectsResponse struct {
	Total int                       `json:"total"`
	List  []*GetObjectsResponseItem `json:"list"`
}

// Получение списка объектов
// @Summary Получение объектов
// @Tags Objects
// @Description Получение объектов
// @ID GetObjects
// @Produce json
// @Param filter_by_id query string false "ID объекта"
// @Param filter_by_parent_id query string false "ID родительского объекта"
// @Param filter_by_zone_id query string false "ID зоны"
// @Param filter_by_category query string false "Категория" Enums(controller,sensor,regulator,generic_input, sensor_value)
// @Param filter_by_type query string false "Тип" example(mega_d, htu21d, regulator, generic_input)
// @Param filter_by_name query string false "Название"
// @Param filter_by_status query string false "Статус" Enums(ON,OFF,Enable,N/A)
// @Param tags query string false "Тэги"
// @Param children query string false "Возвращать дочерние объекты (0-без детей, 1-дети, 2-дети+внуки и т.д.)" default(1)
// @Param type_children query string false "Тип выводимых дочерних элементов (all - все, internal - только внутр., external - только внешние)" Enums(all, internal, external)
// @Param type_struct query string false "Тип структуры в ответе" Enums(easy, full)
// @Param with_methods query string false "Добавить методы в структуру" Enums(true, false)
// @Param with_tags query string false "Добавлять тэги в структуру" Enums(true, false)
// @Param simple_tree query string false "Дерево объектов будет строиться от отфильтрованного объекта вниз, но не вверх" Enums(true, false)
// @Param offset query string false "Смещение" default(0)
// @Param limit query string false "Лимит" default(20)
// @Success      200 {object} http.Response[any]
// @Failure      400 {object} http.Response[any]
// @Failure      500 {object} http.Response[any]
// @Router /objects [get]
func (o *Server) handleGetObjects(ctx *fasthttp.RequestCtx) (interface{}, int, error) {
	params, err := parseGetObjectsParams(ctx)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	items, err := filterObjects(params)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	total := len(items)

	if params.Offset >= len(items) {
		return GetObjectsResponse{Total: total, List: []*GetObjectsResponseItem{}}, http.StatusOK, nil
	}
	items = items[params.Offset:]

	if params.Limit < len(items) {
		items = items[:params.Limit]
	}

	rows := make([]*GetObjectsResponseItem, 0, len(items))
	for _, item := range items {
		rows = append(rows, NewGetObjectsResultItem(item, params.AllTypes, params.WithParents, params.WithMethods, params.WithTags, params.Generations))
	}

	// Если упрощенный вывод включен, то сразу выводим плоскую модель без вложенности
	if params.TypeStruct == "easy" {
		for _, row := range rows {
			row.ParentID = nil
			row.ZoneID = nil
			row.Category = ""
			row.Status = ""
			row.Enabled = nil
			row.Children = nil
		}

		return GetObjectsResponse{
			Total: total,
			List:  rows,
		}, http.StatusOK, nil
	}

	// Рекурсивно загружаем родителей
	if params.WithParents {
		// Создаем общий список всех эл-ов
		m := make(map[int]*GetObjectsResponseItem, len(rows))
		for _, row := range rows {
			m[row.ID] = row
		}

		if err := loadParents(m, rows, params.WithMethods, params.WithTags); err != nil {
			return nil, http.StatusInternalServerError, err
		}

		// Детей добавляем к родителям
		for _, item := range m {
			if item.ParentID != nil {
				parent, ok := m[*item.ParentID]
				if !ok {
					return nil, http.StatusInternalServerError, errors.Errorf("Parent(ID=%d) not found for object(ID=%d)", *item.ParentID, item.ID)
				}

				parent.Children = append(parent.Children, item)
			}
		}

		// Выбираем эл-ты верхнего уровня
		rows = rows[:0]
		for _, item := range m {
			if item.ParentID == nil {
				rows = append(rows, item)
			}
		}
	}

	sortObjectsTree(rows)

	return GetObjectsResponse{
		Total: total,
		List:  rows,
	}, http.StatusOK, nil
}

func loadParents(m map[int]*GetObjectsResponseItem, items []*GetObjectsResponseItem, withMethods, withTags bool) error {
	parents := make([]*GetObjectsResponseItem, 0, len(items))
	for _, item := range items {
		if item.ParentID == nil {
			continue
		}

		// Если родителя уже загрузили, переходим к следующему
		if _, ok := m[*item.ParentID]; ok {
			continue
		}

		obj, err := memStore.I.GetObject(*item.ParentID)
		if err != nil {
			return errors.Wrap(err, "loadParents")
		}

		parent := NewGetObjectsResultItem(obj, true, true, withMethods, withTags, 0)
		parents = append(parents, parent)
		m[parent.ID] = parent
	}

	if len(parents) == 0 {
		return nil
	}

	// В списке новых родителей могут быть их родители отсутствующие в списке
	return loadParents(m, parents, withMethods, withTags)
}

func sortObjectsTree(items []*GetObjectsResponseItem) {
	sort.Slice(items, func(i, j int) bool {
		switch {
		case items[i].Category != items[j].Category:
			return items[i].Category < items[j].Category
		case items[i].Type != items[j].Type:
			return items[i].Type < items[j].Type
		default:
			return items[i].Name < items[j].Name
		}
	})

	for _, item := range items {
		if len(item.Children) > 0 {
			sortObjectsTree(item.Children)
		}

		item.ParentID = nil
	}
}

func NewGetObjectsResultItem(obj objects.Object, allTypes, withParents, withMethods, withTags bool, generations int) *GetObjectsResponseItem {
	enabled := obj.GetEnabled()

	r := &GetObjectsResponseItem{
		ID:       obj.GetID(),
		ZoneID:   obj.GetZoneID(),
		Category: obj.GetCategory(),
		Type:     obj.GetType(),
		Name:     obj.GetName(),
		Status:   obj.GetStatus(),
		Enabled:  &enabled,
	}

	if withParents {
		r.ParentID = obj.GetParentID()
	}

	if withMethods {
		for _, method := range obj.GetMethods().GetAll() {
			r.Methods = append(r.Methods, &GetObjectsResponseItemMethod{
				Name:        method.Name,
				Description: method.Description,
			})
		}
	}

	if withTags {
		r.Tags = obj.GetTags()
	}

	if generations < 1 {
		return r
	}

	for _, child := range obj.GetChildren().GetAll() {
		// Если нет флага получения всех объектов и объект имеет флаг "Скрытый", то его пропускаем
		if !allTypes && child.GetFlags().Has(objects.Hidden) {
			continue
		}

		r.Children = append(r.Children, NewGetObjectsResultItem(child, allTypes, true, withMethods, withTags, generations-1))
	}

	return r
}

type GetObjectsResponseItem struct {
	ID       int                             `json:"id"`                  // ID объекта
	ParentID *int                            `json:"parent_id,omitempty"` // ID родительского объекта
	ZoneID   *int                            `json:"zone_id,omitempty"`   // ID зоны, в которой размещен объект
	Category string                          `json:"category,omitempty"`  // Категория объекта
	Type     string                          `json:"type,omitempty"`      // Тип объекта
	Name     string                          `json:"name,omitempty"`      // Название объекта
	Status   string                          `json:"status,omitempty"`    // Состояние объекта
	Tags     []string                        `json:"tags,omitempty"`      //
	Enabled  *bool                           `json:"enabled,omitempty"`   // Включает методы Start/Shutdown
	Methods  []*GetObjectsResponseItemMethod `json:"methods,omitempty"`   //

	Children []*GetObjectsResponseItem `json:"children,omitempty"` // Дочерние объекты
}

type GetObjectsResponseItemMethod struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}
