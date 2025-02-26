-- +goose Up
-- +goose StatementBegin
create table om_objects_dg_tmp
(
    id        INTEGER not null primary key autoincrement,
    parent_id INTEGER references om_objects on update cascade on delete set null,
    zone_id   INTEGER references zones on update cascade on delete set null,
    category  TEXT    not null,
    type      TEXT    not null,
    name      TEXT    not null,
    status    TEXT default 'N/A' not null,
    tags      JSON default '{}' not null,
    enabled   bool default true not null
);

insert into om_objects_dg_tmp(id, parent_id, zone_id, category, type, name, status, tags, enabled)
select id,parent_id,zone_id,category,type,name,status,tags,enabled
from om_objects;

drop table om_objects;
alter table om_objects_dg_tmp rename to om_objects;

create index parent_id on om_objects (parent_id);
create index tags on om_objects (tags);
create index zone_id on om_objects (zone_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table om_objects add internal bool default false not null;
update om_objects set internal = true
  where category = 'port' and type = 'port_mega_d'
     or category = 'sensor_value'
     or category = 'regulator' and type = 'regulator';
-- +goose StatementEnd
