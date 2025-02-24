-- +goose Up
-- +goose StatementBegin
alter table om_objects add enabled bool default true not null;

UPDATE om_objects SET enabled = false
FROM (SELECT object_id, value FROM om_props WHERE code = 'enable' and value = 'false') AS props
WHERE om_objects.id = props.object_id;

delete from om_props where code='enable';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
create table om_objects_dg_tmp (
    id        INTEGER not null primary key autoincrement,
    parent_id INTEGER references om_objects on update cascade on delete set null,
    zone_id   INTEGER references zones on update cascade on delete set null,
    category  TEXT    not null,
    type      TEXT    not null,
    internal  bool default false not null,
    name      TEXT    not null,
    status    TEXT default 'N/A' not null,
    tags      JSON default '{}' not null
);

insert into om_objects_dg_tmp(id, parent_id, zone_id, category, type, internal, name, status, tags)
select id, parent_id, zone_id, category, type, internal, name, status, tags from om_objects;

insert into om_props(object_id, code, value)
select id, 'enable', case when enabled then 'true' else 'false' end from om_objects
where category = 'conditioner' or type in ('wb_mrm2_mini', 'regulator') or (category = 'sensor' and type in ('motion', 'presence'));

drop table om_objects;

alter table om_objects_dg_tmp rename to om_objects;

create index parent_id on om_objects (parent_id);
create index tags on om_objects (tags);
create index zone_id on om_objects (zone_id);
-- +goose StatementEnd
