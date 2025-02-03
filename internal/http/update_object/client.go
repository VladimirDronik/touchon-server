package update_object

import (
	"touchon-server/internal/context"
	"touchon-server/internal/model"
	httpClient "touchon-server/lib/http/client"
)

// updateSensorCronTask отправляет в action-router запрос на добавление задачи и действия для крона
func updateSensorCronTask(req *Request) (bool, error) {
	cronAction := model.CronTask{}
	arBaseUrl := "http://" + context.Config["action_router_addr"]

	_, ok := req.Props["update_interval"].(string)
	if !ok {
		return false, nil
	}
	cronAction.Period = req.Props["update_interval"].(string)

	cronAction.Actions = append(cronAction.Actions,
		&model.CronAction{
			Enabled:    true,
			TargetType: "object",
			Type:       "method",
			TargetID:   req.ID,
			Name:       "check",
		})

	if _, err := httpClient.I.DoRequest("PUT", arBaseUrl+"/cron/task", nil, nil, cronAction); err != nil {
		context.Logger.Error(err)
		return true, err
	}

	return true, nil
}
