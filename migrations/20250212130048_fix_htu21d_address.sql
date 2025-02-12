-- +goose Up
-- +goose StatementBegin
update om_props set value = '65;66' where object_id = 82 and code = 'address';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
update om_props set value = '0' where object_id = 82 and code = 'address';
-- +goose StatementEnd
