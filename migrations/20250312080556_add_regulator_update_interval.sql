-- +goose Up
-- +goose StatementBegin
insert into om_props(object_id, code, value)
select id, 'update_interval', '30s' from om_objects
where category = 'regulator' and type = 'regulator';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
delete from om_props where code = 'update_interval' and object_id in (
    select id from om_objects where category = 'regulator' and type = 'regulator'
);
-- +goose StatementEnd
