-- +goose Up
-- +goose StatementBegin
drop table conditioner_params;

CREATE TABLE IF NOT EXISTS "conditioner"
(
    id                    INTEGER           not null primary key autoincrement,
    view_item_id          INTEGER           not null unique references view_items on update cascade on delete cascade,
    object_id             INTEGER           INTEGER not null references om_objects on update cascade on delete cascade,
    min_threshold         REAL default 0    not null,
    max_threshold         REAL default 0    not null
);
-- +goose StatementEnd

