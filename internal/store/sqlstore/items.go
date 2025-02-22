package sqlstore

import (
	"encoding/json"
	"sort"

	"github.com/pkg/errors"
	"touchon-server/internal/model"
	"touchon-server/lib/interfaces"
)

type Items struct {
	store *Store
}

// SaveItem создает/обновляет элемент
func (o *Items) SaveItem(item *model.ViewItem) error {
	if item == nil {
		return errors.Wrap(errors.New("item is nil"), "SaveItem")
	}

	if item.ParentID != nil && *item.ParentID <= 0 {
		item.ParentID = nil
	}

	if item.ZoneID != nil && *item.ZoneID <= 0 {
		item.ZoneID = nil
	}

	if err := o.store.db.Save(item).Error; err != nil {
		return errors.Wrap(err, "SaveItem")
	}

	return nil
}

// GetItem получение итема по его ID
func (o *Items) GetItem(itemID int) (*model.ViewItem, error) {
	var r *model.ViewItem

	err := o.store.db.Table("view_items").
		Select("*").
		Where("id = ?", itemID).
		First(&r).Error

	if err != nil {
		return nil, errors.Wrap(err, "GetItem")
	}

	return r, nil

}

// UpdateItem обновление итема
func (o *Items) UpdateItem(ViewItem *model.ViewItem) error {

	if err := o.store.db.Table("view_items").Updates(ViewItem).Error; err != nil {
		return errors.Wrap(err, "UpdateItem")
	}

	return nil
}

// DeleteItem удаление итема
func (o *Items) DeleteItem(itemID int) error {

	if err := o.store.db.Table("view_items").
		Where("id = ?", itemID).Delete(&model.ViewItem{}).
		Error; err != nil {
		return errors.Wrap(err, "UpdateItem")
	}

	return nil
}

// GetScenarios Функция извлекает данные для главной комнтаты на дашборде и отдает их в формате модели
func (o *Items) GetScenarios(withDisabledItems bool) ([]*model.ViewItem, error) {
	r := make([]*model.ViewItem, 0)

	q := o.store.db.Where("zone_id IS NULL").Where("type IN ('switch','button')")

	if withDisabledItems == false {
		q = q.Where("enabled")
	}

	if err := q.Order("sort").Find(&r).Error; err != nil {
		return nil, errors.Wrap(err, "GetScenarios")
	}

	return r, nil
}

// GetZoneItems Функция извлекает данные для помещений дашборда (кроме главной комнаты) и отдает их в формате модели
// Сначала запрашиваем все группы, у которых есть какие-то элементы отображения, затем в полученном результирующем
// наборе для каждой группы запрашиваем элементы отображения, которые относятся к группе. Вместе с этими элементами
// отображения приходят помещения, в которых эти элементы располагаются. Из этих помещений формируем массив,
// который прикрепляем к группе параллельно элементам отображения
func (o *Items) GetZoneItems(withEmptyRooms bool, withDisabledItems bool) ([]*model.GroupRoom, error) {
	var r []*model.GroupRoom

	zones, err := o.GetZones(withEmptyRooms)
	if err != nil {
		return nil, errors.Wrap(err, "GetZoneItems")
	}

	for _, zone := range zones {
		groupRoom := &model.GroupRoom{
			ID:    zone.ID,
			Name:  zone.Name,
			Style: zone.Style,
		}

		groupRoom.Sensors, err = o.GetZoneSensors(zone.ID, withDisabledItems)
		if err != nil {
			return nil, errors.Wrap(err, "GetZoneItems")
		}

		groupRoom.Items, err = o.GetItems(zone.ID, withDisabledItems)
		if err != nil {
			return nil, errors.Wrap(err, "GetZoneItems")
		}

		r = append(r, groupRoom)
	}

	return r, nil
}

type getZonesRow struct {
	*model.Zone
	ItemsCount int            `json:"-"`
	Children   []*getZonesRow `gorm:"-"`
}

// GetZones Получение списка помещений в которых имеются итемы (либо без итемов с флагом withEmptyRooms)
func (o *Items) GetZones(withEmptyRooms bool) ([]*model.Zone, error) {
	var rows []*getZonesRow

	q := `select zones.id, zones.parent_id, zones.name, zones.style, count(view_items.id) items_count
	 from zones
	 left join view_items ON view_items.zone_id = zones.id
	 where view_items.id is null or view_items.enabled and view_items.type != 'group'
	 group by zones.id;`

	if err := o.store.db.Raw(q).Scan(&rows).Error; err != nil {
		return nil, errors.Wrap(err, "GetZones")
	}

	// Строит карту для быстрого поиска эл-ов
	m := make(map[int]*getZonesRow, len(rows))
	for _, row := range rows {
		m[row.ID] = row
	}

	roots := 0

	// Добавляем детей родителям
	for _, row := range m {
		if row.ParentID == 0 {
			roots += 1
			continue
		}

		parent, ok := m[row.ParentID]
		if !ok {
			return nil, errors.Wrap(errors.Errorf("can't find parent with ID %d", row.ParentID), "GetZones")
		}

		parent.Children = append(parent.Children, row)
	}

	// Собираем корневые эл-ты
	items := make([]*getZonesRow, 0, roots)
	for _, row := range m {
		if row.ParentID == 0 {
			items = append(items, row)
		}
	}

	// Рекурсивно удаляем пустые зоны и сортируем
	items = deleteEmptyZonesAndSort(items, withEmptyRooms)

	// Рекурсивно перекладываем данные
	r := zonesRowToZones(items)

	return r, nil
}

func deleteEmptyZonesAndSort(items []*getZonesRow, withEmptyRooms bool) []*getZonesRow {
	r := items[:0]

	for _, item := range items {
		if len(item.Children) > 0 {
			item.Children = deleteEmptyZonesAndSort(item.Children, withEmptyRooms)
		}

		if len(item.Children) > 0 || item.ItemsCount > 0 || withEmptyRooms {
			r = append(r, item)
		}
	}

	if len(r) > 0 {
		sort.Slice(r, func(i, j int) bool {
			switch {
			case r[i].Sort != r[j].Sort:
				return r[i].Sort < r[j].Sort
			default:
				return r[i].ID < r[j].ID
			}
		})
	}

	return r
}

func zonesRowToZones(items []*getZonesRow) []*model.Zone {
	r := make([]*model.Zone, 0, len(items))

	for _, item := range items {
		r = append(r, item.Zone)

		if len(item.Children) > 0 {
			item.Zone.Children = zonesRowToZones(item.Children)
		}
	}

	return r
}

// GetGroupElements Получение элементов группы
func (o *Items) GetGroupElements(groupID int) ([]*model.ViewItem, error) {
	var r []*model.ViewItem

	err := o.store.db.Table("view_items").
		Select("id, status, icon, type, auth, title, color").
		Where("type NOT IN (?, ?)", "sensor", "scenario").
		Where("enabled").
		Where("parent_id = ?", groupID).
		Order("sort").
		Find(&r).Error

	if err != nil {
		return nil, errors.Wrap(err, "GetGroupElements")
	}

	return r, nil
}

// GetCountersList Получение списка счетчиков
func (o *Items) GetCountersList() ([]*model.Counter, error) {
	var r []*model.Counter

	err := o.store.db.Table("counters").
		Select("*").
		Where("enabled").
		Order("sort").
		Find(&r).Error

	if err != nil {
		return nil, errors.Wrap(err, "GetCountersList")
	}

	return r, nil
}

// GetCounter Получение счетчика
func (o *Items) GetCounter(id int) (*model.Counter, error) {
	var r *model.Counter

	err := o.store.db.Table("counters").
		Select("*").
		Where("id = ?", id).
		First(&r).Error

	if err != nil {
		return nil, errors.Wrap(err, "GetCounter")
	}

	return r, nil
}

// ChangeItem обработка изменения итема
func (o *Items) ChangeItem(id int, status string) error {
	if err := o.store.db.Table("view_items").Where("id = ?", id).Update("status", status).Error; err != nil {
		return errors.Wrap(err, "ChangeItem")
	}

	return nil
}

func (o *Items) GetItemsForChange(targetType interfaces.TargetType, targetID int, eventName string) ([]*model.ItemForWS, error) {
	var items []*model.ItemForWS

	err := o.store.db.Table("events").Select("view_items.id, target_id, status, params, view_items.type, value AS EventValue").
		Where("target_type = ?", targetType).
		Where("target_id = ?", targetID).
		Where("event = ?", eventName).
		InnerJoins("INNER JOIN view_items ON view_items.id = events.item_id").
		Find(&items).Error

	if err != nil {
		return nil, errors.Wrap(err, "GetItemsForChange")
	}

	return items, nil
}

// GetZone Получение помещения
func (o *Items) GetZone(roomID int) (*model.GroupRoom, error) {
	var row *model.Zone

	err := o.store.db.Table("zones").
		Select("*").
		First(&row, "id = ?", roomID).Error

	if err != nil {
		return nil, errors.Wrap(err, "GetZone")
	}

	zone := &model.GroupRoom{
		ID:    row.ID,
		Name:  row.Name,
		Style: row.Style,
	}

	zone.Sensors, err = o.GetZoneSensors(zone.ID, false)
	if err != nil {
		return nil, errors.Wrap(err, "GetZone")
	}

	zone.Items, err = o.GetItems(zone.ID, false)
	if err != nil {
		return nil, errors.Wrap(err, "GetZone")
	}

	return zone, nil
}

func (o *Items) GetMenus(parentIDs ...int) ([]*model.Menu, error) {
	m := make(map[int]*model.Menu, 10)
	if err := o.getMenus(m, parentIDs...); err != nil {
		return nil, errors.Wrap(err, "GetMenus")
	}

	for _, item := range m {
		if json.Valid([]byte(item.Params)) {
			if err := json.Unmarshal([]byte(item.Params), &item.ParamsOutput); err != nil {
				return nil, errors.Wrap(err, "GetMenus")
			}
		}

		if item.ParentID == 0 {
			continue
		}

		if parent, ok := m[item.ParentID]; ok {
			parent.Children = append(parent.Children, item)
		}

	}

	r := make([]*model.Menu, 0, 10)
	for _, item := range m {
		if _, ok := m[item.ParentID]; !ok {
			r = append(r, item)
		}
	}

	sortMenus(r)

	return r, nil
}

func sortMenus(items []*model.Menu) {
	sort.Slice(items, func(i, j int) bool {
		switch {
		case items[i].Sort != items[j].Sort:
			return items[i].Sort < items[j].Sort
		default:
			return items[i].Title < items[j].Title
		}
	})

	for _, item := range items {
		if len(item.Children) > 0 {
			sortMenus(item.Children)
		}
	}
}

func (o *Items) getMenus(m map[int]*model.Menu, parentIDs ...int) error {
	var rows []*model.Menu

	err := o.store.db.Table("menus").
		Select("*").
		Where("parent_id in ?", parentIDs).
		Where("enabled").
		Find(&rows).Error

	if err != nil {
		return errors.Wrap(err, "getMenus")
	}

	if len(rows) == 0 {
		return nil
	}

	ids := make([]int, 0, len(rows))
	for _, row := range rows {
		ids = append(ids, row.ID)
		m[row.ID] = row
	}

	if err := o.getMenus(m, ids...); err != nil {
		return err
	}

	return nil
}

func (o *Items) GetZoneSensors(zoneID int, withDisabledItems bool) ([]*model.SensorItem, error) {
	zones, err := o.store.Zones().GetZoneTrees(zoneID)
	if err != nil {
		return nil, errors.Wrap(err, "GetZoneSensors")
	}

	zoneIDs := collectZoneIDs(nil, zones)

	r, err := o.getSensors(withDisabledItems, zoneIDs...)
	if err != nil {
		return nil, errors.Wrap(err, "GetZoneSensors")
	}

	return r, nil
}

func collectZoneIDs(ids []int, zones []*model.Zone) []int {
	for _, item := range zones {
		ids = append(ids, item.ID)

		if len(item.Children) > 0 {
			ids = append(ids, collectZoneIDs(ids, item.Children)...)
		}
	}

	return ids
}

func (o *Items) getSensors(withDisabledItems bool, zoneIDs ...int) ([]*model.SensorItem, error) {
	var sensorsStorage, sensors []*model.SensorItem

	q := o.store.db.Table("view_items").
		Select("view_items.id AS item_id, view_items.title AS title,"+
			" view_items.zone_id, sensors.adjustment AS adjustment,"+
			"view_items.icon AS icon, sensors.type AS type, view_items.auth AS auth, sensors.object_id AS object_id").
		Joins("INNER JOIN sensors ON view_items.id = sensors.view_item_id").
		Where("view_items.zone_id in ?", zoneIDs)

	if withDisabledItems == false {
		q = q.Where("view_items.enabled")
	}

	err := q.Order("view_items.sort").
		Find(&sensorsStorage).Error

	if err != nil {
		return nil, errors.Wrap(err, "getSensors")
	}

	if len(sensorsStorage) == 0 {
		return sensors, nil
	}

	for _, sensor := range sensorsStorage {
		err = o.store.db.Table("om_props").Select("value AS current").
			Where("object_id = ?", sensor.ObjectID).
			Where("code = 'value'").
			Find(&sensor.Current).Error
		if err != nil {
			return nil, errors.Wrap(err, "getSensors")
		}

		sensors = append(sensors, sensor)
	}

	return sensors, nil
}

// GetItems Получение устройств
func (o *Items) GetItems(zoneID int, withDisabledItems bool) ([]*model.ViewItem, error) {
	zones, err := o.store.Zones().GetZoneTrees(zoneID)
	if err != nil {
		return nil, errors.Wrap(err, "GetItems")
	}

	zoneIDs := collectZoneIDs(nil, zones)

	r, err := o.getItems(withDisabledItems, zoneIDs...)
	if err != nil {
		return nil, errors.Wrap(err, "GetItems")
	}

	return r, nil
}

func (o *Items) getItems(withDisabledItems bool, zoneIDs ...int) ([]*model.ViewItem, error) {
	var r []*model.ViewItem

	q := o.store.db.Table("view_items").
		Select("id, type, title, icon, auth, status, params, color")

	if withDisabledItems == false {
		q = q.Where("enabled")
	}

	err := q.Where("type NOT IN ('sensor', 'scenario')").Where("zone_id IN ?", zoneIDs).Order("sort").Find(&r).Error

	if err != nil {
		return nil, errors.Wrap(err, "getItems")
	}

	return r, nil
}

// SetOrder задает порядок сортировки
func (o *Items) SetOrder(itemIDs []int, zoneID int) error {
	for i, itemID := range itemIDs {
		if err := o.SetFieldValue(itemID, "sort", i+1, zoneID); err != nil {
			return errors.Wrap(err, "SetOrder")
		}
	}

	return nil
}

func (o *Items) SetFieldValue(itemID int, field string, value interface{}, zoneID int) error {
	err := o.store.db.
		Table("view_items").
		Where("id = ?", itemID).
		//	Where("zone_id = ?", zoneID).
		Update(field, value).
		Error

	if err != nil {
		return errors.Wrap(err, "SetFieldValue")
	}

	return nil
}
