-- +goose Up
-- +goose StatementBegin
delete from ar_cron_actions where id = 1;
delete from ar_cron_tasks where id = 1;
update om_objects set enabled = 0 where id in (85, 89, 93, 78, 97, 91);
UPDATE om_props SET value = '120s' where object_id = 304 and code = 'sensor_value_ttl';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
INSERT INTO ar_cron_tasks (id, name, description, period, enabled) VALUES (1, 'Test htu21d 5s', 'Проверка датчика htu21d', '5s', 1);
INSERT INTO ar_cron_actions (id, task_id, target_type, target_id, type, name, args, qos, enabled, sort, comment) VALUES (1, 1, 'object', 82, 'method', 'check', '{}', 0, 1, 0, '');
update om_objects set enabled = 1 where id in (85, 89, 93, 78, 97, 91);
UPDATE om_props SET value = '30s' where object_id = 304 and code = 'sensor_value_ttl';
-- +goose StatementEnd
