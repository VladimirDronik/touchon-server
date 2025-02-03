package http

import (
	"strconv"

	"touchon-server/internal/context"
	httpClient "touchon-server/lib/http/client"
	"touchon-server/lib/mqtt/messages"
)

// deleteEvents удаляет события для объекта в action-router
func deleteEvent(objectID int) error {
	arBaseUrl := "http://" + context.Config["action_router_addr"]
	params := map[string]string{
		"target_type": string(messages.TargetTypeObject),
		"target_id":   strconv.Itoa(objectID),
		"event_name":  "all",
	}

	//Удаляем все возможные события объекта
	if _, err := httpClient.I.DoRequest("DELETE", arBaseUrl+"/events", params, nil, nil); err != nil {
		context.Logger.Error(err)
		return err
	}

	//Удаляем все действия для сторонних событий, где может фигурировать объект
	if _, err := httpClient.I.DoRequest("DELETE", arBaseUrl+"/events/all-actions", params, nil, nil); err != nil {
		context.Logger.Error(err)
		return err
	}

	//Удаляем все действия крона для объекта
	if _, err := httpClient.I.DoRequest("DELETE", arBaseUrl+"/cron/task", params, nil, nil); err != nil {
		context.Logger.Error(err)
		return err
	}

	return nil
}
