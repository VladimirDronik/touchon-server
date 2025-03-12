-- +goose Up
-- +goose StatementBegin
create table zones_dg_tmp (
    id        INTEGER not null primary key autoincrement,
    parent_id INTEGER references zones on update cascade on delete set null,
    name      TEXT    default '' not null,
    style     TEXT    default '' not null,
    sort      INTEGER default 0 not null,
    is_group  BOOL    default false not null
);

insert into zones_dg_tmp(id, parent_id, name, style, sort, is_group)
select id, parent_id, name, style, sort, is_group
from zones;

drop table zones;
alter table zones_dg_tmp rename to zones;

update zones set parent_id = null where parent_id = 0;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
update zones set parent_id = 0 where parent_id is null;

create table zones_dg_tmp (
    id        INTEGER not null primary key autoincrement,
    parent_id INTEGER default 0 references zones on update cascade on delete set default,
    name      TEXT    default '' not null,
    style     TEXT    default '' not null,
    sort      INTEGER default 0 not null,
    is_group  BOOL    default false not null
);

insert into zones_dg_tmp(id, parent_id, name, style, sort, is_group)
select id, parent_id, name, style, sort, is_group
from zones;

drop table zones;
alter table zones_dg_tmp rename to zones;
-- +goose StatementEnd
