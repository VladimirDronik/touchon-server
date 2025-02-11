-- +goose Up
-- +goose StatementBegin
create table view_items_dg_tmp (
    id            INTEGER not null primary key autoincrement,
    parent_id     INTEGER references view_items on update cascade on delete set null,
    zone_id       INTEGER references zones on update cascade on delete set null,
    type          TEXT    default '' not null,
    status        TEXT    default '' not null,
    icon          TEXT    default '' not null,
    title         TEXT    default '' not null,
    sort          INTEGER default 0 not null,
    params        JSON    default '{}' not null,
    color         TEXT    default '' not null,
    auth          TEXT    default '' not null,
    description   TEXT    default '' not null,
    position_left INTEGER default 0 not null,
    scene         INTEGER default 0 not null,
    position_top  INTEGER default 0 not null,
    enabled       bool    default true not null
);

insert into view_items_dg_tmp(
id, parent_id, zone_id, type, status, icon, title, sort, params,
color, auth, description, position_left, scene, position_top, enabled)
select id, parent_id, zone_id, type, status, icon, title, sort,
       params, color, auth, description, position_left, scene,
       position_top, enabled
from view_items;

drop table view_items;

alter table view_items_dg_tmp rename to view_items;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
CREATE TABLE view_items_dg_tmp_2 (
    id            INTEGER not null primary key autoincrement,
    parent_id     INTEGER not null default 0 references view_items on update cascade on delete set default,
    zone_id       INTEGER not null default 0 references zones on update cascade on delete set default,
    type          TEXT    not null default '',
    status        TEXT    not null default '',
    icon          TEXT    not null default '',
    title         TEXT    not null default '',
    sort          INTEGER not null default 0,
    params        JSON    not null default '{}',
    color         TEXT    not null default '',
    auth          TEXT    not null default '',
    description   TEXT    not null default '',
    position_left INTEGER not null default 0,
    scene         INTEGER not null default 0,
    position_top  INTEGER not null default 0,
    enabled       bool    not null default true
);

insert into view_items_dg_tmp_2 (
    id, parent_id, zone_id, type, status, icon, title, sort, params,
    color, auth, description, position_left, scene, position_top, enabled)
select id, parent_id, zone_id, type, status, icon, title, sort,
       params, color, auth, description, position_left, scene,
       position_top, enabled
from view_items;

drop table view_items;

alter table view_items_dg_tmp_2 rename to view_items;

-- +goose StatementEnd
