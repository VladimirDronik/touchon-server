-- +goose Up
-- +goose StatementBegin
update om_objects set category = 'rs485', type = 'bus', tags = '{"rs485":true}' where category = 'modbus' and type = 'modbus';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
update om_objects set category = 'modbus', type = 'modbus', tags = '{"modbus":true}' where category = 'rs485' and type = 'bus';
-- +goose StatementEnd
