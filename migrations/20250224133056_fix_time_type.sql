-- +goose Up
-- +goose StatementBegin
update om_props set value = value || 's' where code in ('sensor_value_ttl', 'timeout', 'update_interval', 'period');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
update om_props set value = rtrim(value, 's') where code in ('sensor_value_ttl', 'timeout', 'update_interval', 'period');
-- +goose StatementEnd
