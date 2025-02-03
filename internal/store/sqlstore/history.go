package sqlstore

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/pkg/errors"
	"translator/internal/helpers"
	"translator/internal/model"
)

type History struct {
	store *Store
}

// GetHistory получение истории изменения значений
func (o *History) GetHistory(itemID int, itemType model.HistoryItemType, filter model.HistoryFilter) (*model.HistoryPoints, error) {
	var filters []model.HistoryFilter

	if filter != "" {
		filters = append(filters, filter)
	} else {
		filters = append(filters,
			model.HistoryFilterDay,
			model.HistoryFilterWeek,
			model.HistoryFilterMonth,
			model.HistoryFilterYear,
		)
	}

	r := &model.HistoryPoints{ServerDate: getCurrentDate()}

	for _, filter := range filters {
		points, err := o.getHistory(itemID, itemType, filter)
		if err != nil {
			return nil, errors.Wrap(err, "GetHistory")
		}

		switch filter {
		case model.HistoryFilterDay:
			r.DayPoints = points
		case model.HistoryFilterWeek:
			r.WeekPoints = points
		case model.HistoryFilterMonth:
			r.MonthPoints = points
		case model.HistoryFilterYear:
			r.YearPoints = points
		}
	}

	return r, nil
}

func (o *History) getHistory(itemID int, itemType model.HistoryItemType, filter model.HistoryFilter) ([]*model.HistoryPoint, error) {
	tableName, fieldIDName, err := getHistoryTableName(filter, itemType)
	if err != nil {
		return nil, nil
		//return nil, errors.Wrap(err, "GetHistory")
	}

	startDate, endDate, err := getDates(filter, itemType)
	if err != nil {
		return nil, errors.Wrap(err, "getHistory")
	}

	r := make([]*model.HistoryPoint, 0)
	err = o.store.db.Table(tableName).
		Select("*").
		Where(fieldIDName+" = ?", itemID).
		Where("datetime BETWEEN ? AND ?", startDate.Format("2006-01-02T15:04"), endDate.Format("2006-01-02T15:04")).
		Order("datetime").
		Find(&r).Error

	if err != nil {
		return nil, errors.Wrap(err, "getHistory")
	}

	if len(r) == 0 {
		return nil, nil
	}

	r, err = prepareHistoryPoints(r, *startDate, *endDate, filter)
	if err != nil {
		return nil, errors.Wrap(err, "getHistory")
	}

	return r, nil
}

func (o *History) GenerateHistory(itemID int, itemType model.HistoryItemType, filter model.HistoryFilter, startDate, endDate string, min, max float32) error {
	startTime, err := time.Parse("2006-01-02 15:04", startDate)
	if err != nil {
		return errors.Wrap(err, "GenerateHistory")
	}

	endTime, err := time.Parse("2006-01-02 15:04", endDate)
	if err != nil {
		return errors.Wrap(err, "GenerateHistory")
	}

	tableName, fieldIDName, err := getHistoryTableName(filter, itemType)
	if err != nil {
		return errors.Wrap(err, "GenerateHistory")
	}

	getValue := func() float32 {
		v := min + rand.Float32()*(max-min)
		return float32(math.Round(float64(v)*10) / 10)
	}

	var format string
	var step func(time.Time) time.Time

	switch filter {
	case model.HistoryFilterDay:
		format = "2006-01-02T15:04"
		step = func(t time.Time) time.Time { return t.Add(time.Hour) }

	case model.HistoryFilterWeek, model.HistoryFilterMonth:
		format = "2006-01-02"
		step = func(t time.Time) time.Time { return t.AddDate(0, 0, 1) }

	case model.HistoryFilterYear:
		format = "2006-01"
		step = func(t time.Time) time.Time { return t.AddDate(0, 1, 0) }

	default:
		return errors.Wrap(errors.Errorf("unknown filter %q", filter), "GenerateHistory")
	}

	q := fmt.Sprintf("INSERT INTO %s (datetime, %s, value) VALUES (?, ?, ?)", tableName, fieldIDName)

	for currentDate := startTime; currentDate.Before(endTime) || currentDate.Equal(endTime); currentDate = step(currentDate) {
		if err := o.store.db.Exec(q, currentDate.Format(format), itemID, getValue()).Error; err != nil {
			return errors.Wrap(err, "GenerateHistory")
		}
	}

	return nil
}

// SetHourlyValue Добавляет в график ежечасное значение
func (o History) SetHourlyValue(itemID int, dateTime string, value float32) error {
	type recordStruct struct {
		ViewItemID int `gorm:"view_item_id"`
		Datetime   string
		Value      float32
	}

	row := recordStruct{ViewItemID: itemID, Datetime: dateTime, Value: float32(math.Round(float64(value)*10)) / 10}

	r := o.store.db.Table("device_hourly_history").Where("view_item_id = ?", itemID).Where("datetime = ?", dateTime).Updates(&row)
	if r.Error != nil {
		return errors.Wrap(r.Error, "SetHourlyValue")
	}

	if r.RowsAffected == 0 {
		if err := o.store.db.Table("device_hourly_history").Create(&row).Error; err != nil {
			return errors.Wrap(err, "SetHourlyValue")
		}
	}

	return nil
}

// prepareHistoryPoints заполнение полей formatted_date
func prepareHistoryPoints(points []*model.HistoryPoint, startDate time.Time, endDate time.Time, filter model.HistoryFilter) ([]*model.HistoryPoint, error) {
	var err error
	r := make([]*model.HistoryPoint, 0, len(points))

	pointsMap := make(map[string]*model.HistoryPoint)
	for _, point := range points {
		pointsMap[point.Datetime] = point
	}

	switch filter {
	case model.HistoryFilterDay:
		startDate = startDate.Truncate(time.Hour).Add(time.Hour)
		for date := startDate; date.Before(endDate); date = date.Add(time.Hour) {
			dateString := date.Format("2006-01-02T15:04")

			point, ok := pointsMap[dateString]
			if !ok {
				point = &model.HistoryPoint{Datetime: dateString}
			}

			point.FormattedDate, err = getFormattedDate(date, filter)
			if err != nil {
				return nil, errors.Wrap(err, "prepareHistoryPoints")
			}

			r = append(r, point)
		}

	case model.HistoryFilterWeek, model.HistoryFilterMonth:
		startDate = startDate.Truncate(time.Hour*24).AddDate(0, 0, 1)

		for date := startDate; date.Before(endDate); date = date.AddDate(0, 0, 1) {
			dateString := date.Format("2006-01-02")

			point, ok := pointsMap[dateString]
			if !ok {
				point = &model.HistoryPoint{Datetime: dateString}
			}

			point.FormattedDate, err = getFormattedDate(date, filter)
			if err != nil {
				return nil, errors.Wrap(err, "prepareHistoryPoints")
			}

			r = append(r, point)
		}

	case model.HistoryFilterYear:
		startDate = time.Date(startDate.Year(), startDate.Month(), 1, 0, 0, 0, 0, startDate.Location()).AddDate(0, 1, 0)
		for date := startDate; date.Before(endDate); date = date.AddDate(0, 1, 0) {
			dateString := date.Format("2006-01")

			point, ok := pointsMap[dateString]
			if !ok {
				point = &model.HistoryPoint{Datetime: dateString}
			}

			point.FormattedDate, err = getFormattedDate(date, filter)
			if err != nil {
				return nil, errors.Wrap(err, "prepareHistoryPoints")
			}

			r = append(r, point)
		}
	}

	return r, nil
}

// getHistoryTableName получить названия таблицы, в которой хранится история значений для определенного типа сущности и фильтра
func getHistoryTableName(filter model.HistoryFilter, itemType model.HistoryItemType) (string, string, error) {
	switch {
	case filter == model.HistoryFilterDay && itemType == model.HistoryItemTypeDeviceObject:
		return "device_hourly_history", "view_item_id", nil
	case filter == model.HistoryFilterWeek && itemType == model.HistoryItemTypeDeviceObject:
		return "device_daily_history", "view_item_id", nil
	case filter == model.HistoryFilterMonth && itemType == model.HistoryItemTypeDeviceObject:
		return "device_daily_history", "view_item_id", nil
	case filter == model.HistoryFilterMonth && itemType == model.HistoryItemTypeCounterObject:
		return "counter_daily_history", "counter_id", nil
	case filter == model.HistoryFilterYear && itemType == model.HistoryItemTypeCounterObject:
		return "counter_monthly_history", "counter_id", nil
	default:
		return "", "", errors.Wrap(errors.Errorf("no table for this type (%s) and filter (%s)", itemType, filter), "getHistoryTableName")
	}
}

// getDates получить начало и конец временного промежутка для определенного фильтра и типа объекта
func getDates(filter model.HistoryFilter, itemType model.HistoryItemType) (*time.Time, *time.Time, error) {
	var startDate time.Time
	var endDate time.Time

	switch {
	case itemType == model.HistoryItemTypeDeviceObject && filter == model.HistoryFilterDay:
		endDate = time.Now()
		startDate = endDate.Add(-24 * time.Hour)

	case itemType == model.HistoryItemTypeDeviceObject && filter == model.HistoryFilterWeek:
		endDate = time.Now()
		startDate = endDate.AddDate(0, 0, -7)

	case itemType == model.HistoryItemTypeDeviceObject && filter == model.HistoryFilterMonth:
		endDate = time.Now()
		startDate = endDate.AddDate(0, -1, 0)

	case itemType == model.HistoryItemTypeCounterObject && filter == model.HistoryFilterMonth:
		now := time.Now()
		startDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		endDate = startDate.AddDate(0, 1, 0)
		startDate = startDate.AddDate(0, 0, -1)

	case itemType == model.HistoryItemTypeCounterObject && filter == model.HistoryFilterYear:
		now := time.Now()
		startDate = time.Date(now.Year(), time.January, 1, 0, 0, 0, 0, now.Location())
		startDate = startDate.AddDate(0, -1, 0)
		endDate = time.Date(now.Year(), time.December, 1, 0, 0, 0, 0, now.Location())
		endDate = endDate.AddDate(0, 1, 0)

	default:
		return nil, nil, errors.Wrap(errors.New("no data for this filter"), "getDates")
	}

	return &startDate, &endDate, nil
}

func getFormattedDate(date time.Time, filter model.HistoryFilter) (string, error) {
	switch filter {
	case model.HistoryFilterDay:
		return date.Format("15:04"), nil
	case model.HistoryFilterWeek:
		return helpers.DayOfWeekEnToRu(date.Weekday().String()), nil
	case model.HistoryFilterMonth:
		return date.Format("02.01"), nil
	case model.HistoryFilterYear:
		return date.Format("01"), nil
	}

	return "", errors.Wrap(errors.Errorf("unknown filter %q", filter), "getFormattedDate")
}

// getCurrentDate получить текущую дату сервера
func getCurrentDate() *model.ServerDate {
	now := time.Now().Truncate(time.Hour)

	return &model.ServerDate{
		Hour:  now.Format("15:04"),
		Day:   now.Format("02"),
		Month: helpers.MonthEnToRu(now.Month().String(), false),
		Year:  now.Format("2006"),
	}
}
