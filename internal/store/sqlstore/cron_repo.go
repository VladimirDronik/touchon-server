package sqlstore

import (
	"encoding/json"
	"sort"

	"action-router/internal/model"
	"github.com/pkg/errors"
)

type CronRepo struct {
	store *Store
}

// GetEnabledTasks получение периодических заданий
func (o *CronRepo) GetEnabledTasks() ([]*model.CronTask, error) {
	type R struct {
		model.CronTask
		model.CronAction
		Args string
	}
	rows := make([]*R, 0, 10)

	q := `
SELECT t.*, a.*
FROM cron_tasks as t
     inner join cron_actions as a on a.task_id = t.id
WHERE t.enabled = 1 AND a.enabled = 1
order by t.id, a.sort, a.id`

	if err := o.store.db.Raw(q).Scan(&rows).Error; err != nil {
		return nil, errors.Wrap(err, "GetEnabledTasks")
	}

	m := make(map[int]*model.CronTask, len(rows))
	for _, row := range rows {
		task, ok := m[row.CronTask.ID]
		if !ok {
			task = &model.CronTask{
				ID:          row.CronTask.ID,
				Name:        row.CronTask.Name,
				Description: row.CronTask.Description,
				Period:      row.CronTask.Period,
				Enabled:     row.CronTask.Enabled,
			}
			m[row.CronTask.ID] = task
		}

		if err := json.Unmarshal([]byte(row.Args), &row.CronAction.Args); err != nil {
			return nil, errors.Wrap(err, "GetEnabledTasks")
		}

		task.Actions = append(task.Actions, &model.CronAction{
			ID:         row.CronAction.ID,
			TaskID:     row.CronAction.TaskID,
			Name:       row.CronAction.Name,
			Type:       row.CronAction.Type,
			TargetID:   row.CronAction.TargetID,
			TargetType: row.CronAction.TargetType,
			QoS:        row.CronAction.QoS,
			Args:       row.CronAction.Args,
			Enabled:    row.CronAction.Enabled,
			Comment:    row.CronAction.Comment,
		})
	}

	r := make([]*model.CronTask, 0, len(m))
	for _, task := range m {
		r = append(r, task)
	}

	sort.Slice(r, func(i, j int) bool {
		return r[i].ID < r[j].ID
	})

	return r, nil
}

func (o *CronRepo) CreateTask(task *model.CronTask) (int, error) {
	err := o.store.db.Create(task).Error
	if err != nil {
		return 0, errors.Wrap(err, "CreateTask")
	}
	return task.ID, nil
}

func (o *CronRepo) UpdateTask(task *model.CronTask) error {
	var objectID int
	var targetType string

	for _, act := range task.Actions {
		objectID = act.TargetID
		targetType = act.TargetType
	}

	cronAction, err := o.GetCronAction(objectID, targetType)
	if err != nil {
		return errors.Wrap(err, "DeleteTask: select CRON action")
	}

	err = o.store.db.Where("id = ?", cronAction.TaskID).Updates(task).Error
	if err != nil {
		return errors.Wrap(err, "UpdateTask")
	}
	return nil
}

func (o *CronRepo) DeleteTask(objectID int, targetType string) error {
	cronAction, err := o.GetCronAction(objectID, targetType)
	if err != nil {
		return errors.Wrap(err, "DeleteTask: select CRON action")
	}

	err = o.store.db.
		Exec("DELETE FROM cron_actions WHERE target_id = ? AND target_type = ?", objectID, targetType).Error
	if err != nil {
		return errors.Wrap(err, "DeleteCronActions")
	}

	err = o.store.db.
		Exec("DELETE FROM cron_tasks WHERE id = ?", cronAction.TaskID).Error
	if err != nil {
		return errors.Wrap(err, "DeleteCronTask")
	}

	return nil
}

func (o *CronRepo) GetCronAction(objectID int, targetType string) (*model.CronAction, error) {
	var cronAction *model.CronAction

	err := o.store.db.
		Where("target_id = ?", objectID).
		Where("target_type = ?", targetType).
		Find(&cronAction).Error
	if err != nil {
		return cronAction, errors.Wrap(err, "DeleteCronActions, find task_id")
	}

	return cronAction, nil
}

func (o *CronRepo) CreateTaskAction(action *model.CronAction) error {
	err := o.store.db.Create(action).Error
	if err != nil {
		return errors.Wrap(err, "CreateTaskAction")
	}
	return nil
}
