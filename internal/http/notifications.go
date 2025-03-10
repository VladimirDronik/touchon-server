package http

import (
	"net/http"
	"sort"
	"time"

	"github.com/valyala/fasthttp"
	"touchon-server/internal/helpers"
	"touchon-server/internal/model"
	"touchon-server/internal/store"
)

// Получение кол-ва непрочитанных уведомлений
// @Security TokenAuth
// @Summary Получение кол-ва непрочитанных уведомлений
// @Tags Notifications
// @Description Получение кол-ва непрочитанных уведомлений
// @ID GetUnreadNotificationsCount
// @Produce json
// @Success      200 {object} Response[int]
// @Failure      400 {object} Response[any]
// @Failure      500 {object} Response[any]
// @Router /private/notifications/unread-count [get]
func (o *Server) getUnreadNotificationsCount(ctx *fasthttp.RequestCtx) (interface{}, int, error) {
	count, err := store.I.Notifications().GetUnreadNotificationsCount()
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	return count, http.StatusOK, nil
}

// Получение списка уведомлений
// @Security TokenAuth
// @Summary Получение списка уведомлений
// @Tags Notifications
// @Description Получение списка уведомлений, сгруппированных по дню
// @ID GetNotifications
// @Produce json
// @Param offset query int false "Offset" default(0)
// @Param limit  query int false "Limit" default(20)
// @Success      200 {object} Response[[]model.GroupedNotifications]
// @Failure      400 {object} Response[any]
// @Failure      500 {object} Response[any]
// @Router /private/notifications [get]
func (o *Server) getNotifications(ctx *fasthttp.RequestCtx) (interface{}, int, error) {
	offset, err := helpers.GetUintParam(ctx, "offset")
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	limit, err := helpers.GetUintParam(ctx, "limit")
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	if limit == 0 {
		limit = 20
	}

	notifications, err := store.I.Notifications().GetNotifications(offset, limit)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	return groupNotificationsByDate(notifications), http.StatusOK, nil
}

// Отмечает уведомление как прочитанное
// @Security TokenAuth
// @Summary Отмечает уведомление как прочитанное
// @Tags Notifications
// @Description Отмечает уведомление как прочитанное
// @ID SetNotificationIsRead
// @Produce json
// @Param notificationId query int true "ID" default(12)
// @Success      200 {object} Response[any]
// @Failure      400 {object} Response[any]
// @Failure      500 {object} Response[any]
// @Router /private/notification [patch]
func (o *Server) setNotificationIsRead(ctx *fasthttp.RequestCtx) (interface{}, int, error) {
	id, err := helpers.GetUintParam(ctx, "notificationId")
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	if err = store.I.Notifications().SetIsRead(id); err != nil {
		return nil, http.StatusBadRequest, err
	}

	return nil, http.StatusOK, nil
}

// groupNotificationsByDate функция преобразования массива уведомлений в список группированный по дню в убывающем по дате порядке
func groupNotificationsByDate(notifications []*model.Notification) []*model.GroupedNotifications {
	today := time.Now().Format("2006-01-02")
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	layout := "2006-01-02T15:04:05"

	groupedMap := make(map[string][]*model.Notification)

	for _, notification := range notifications {
		parsedDate, err := time.Parse(layout, notification.Date)
		if err != nil {
			continue
		}

		dateKey := parsedDate.Format("2006-01-02")

		var groupName string
		if dateKey == today {
			groupName = "today"
		} else if dateKey == yesterday {
			groupName = "yesterday"
		} else {
			groupName = dateKey
		}

		groupedMap[groupName] = append(groupedMap[groupName], notification)
	}

	var groupedNotifications []*model.GroupedNotifications
	for groupName, notifies := range groupedMap {
		groupedNotifications = append(groupedNotifications, &model.GroupedNotifications{
			Day:           groupName,
			Notifications: notifies,
		})
	}

	sort.SliceStable(groupedNotifications, func(i, j int) bool {
		if groupedNotifications[i].Day == "today" {
			return true
		}
		if groupedNotifications[j].Day == "today" {
			return false
		}
		if groupedNotifications[i].Day == "yesterday" {
			return true
		}
		if groupedNotifications[j].Day == "yesterday" {
			return false
		}
		return groupedNotifications[i].Day > groupedNotifications[j].Day
	})

	return groupedNotifications
}
