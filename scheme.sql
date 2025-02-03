-- -------------------------------------------------------------------------------------

CREATE TABLE om_objects
(
    id        INTEGER not null primary key autoincrement,
    parent_id INTEGER,
    zone_id   INTEGER,
    category  TEXT    not null,
    type      TEXT    not null,
    internal  bool default false not null,
    name      TEXT    not null,
    status    TEXT default 'N/A' not null,
    tags      JSON default '{}' not null,

    FOREIGN KEY(parent_id) REFERENCES om_objects(id) on update cascade on delete set null
);

create index parent_id  on om_objects (parent_id);
create index tags on om_objects (tags);
create index zone_id on om_objects (zone_id);


CREATE TABLE om_props (
    id        INTEGER not null primary key autoincrement,
    object_id INTEGER not null,
    code      TEXT not null,
    value     TEXT,

    FOREIGN KEY(object_id) REFERENCES om_objects(id) on update cascade on delete cascade
);

create unique index object_id_code on om_props (object_id, code);

create table om_scripts
(
    id             INTEGER primary key autoincrement,
    code           TEXT not null,
    name           TEXT not null,
    description    TEXT default '' not null,
    params         TEXT default '{}' not null,
    body           TEXT not null,
    name_lowercase TEXT default '' not null
);

create unique index code on om_scripts (code);

-- -------------------------------------------------------------------------------------

create table ar_cron_tasks
(
    id          INTEGER not null primary key autoincrement,
    name        TEXT not null,
    description TEXT,
    period      TEXT,
    enabled     INTEGER default 1
);

CREATE TABLE ar_cron_actions (
    id          integer not null primary key autoincrement,
    task_id     integer not null,
    target_type text not null default '',
    target_id   integer not null default 0,
    type        text not null,
    name        text not null,
    args        text not null default '{}',
    qos         integer not null default 0,
    enabled     integer default 1,
    sort        int not null default 0,
    comment     text not null default '',

    FOREIGN KEY(task_id) REFERENCES ar_cron_tasks(id) on update cascade on delete cascade
);

CREATE INDEX if not exists task_id ON ar_cron_actions(task_id);

CREATE TABLE ar_events (
    id          integer not null primary key autoincrement,
    target_type text not null,
    target_id   integer not null,
    event_name  text not null,
    enabled     integer default 1
);

CREATE UNIQUE INDEX if not exists tt_ti_en ON ar_events(target_type, target_id, event_name);

CREATE TABLE ar_event_actions (
    id          integer not null primary key autoincrement,
    event_id    integer not null,
    target_type text not null default '',
    target_id   integer not null default 0,
    type        text not null,
    name        text not null,
    args        text not null default '{}',
    qos         integer not null default 0,
    enabled     integer default 1,
    sort        int not null default 0,
    comment     text not null default '',

    FOREIGN KEY(event_id) REFERENCES ar_events(id) on update cascade on delete cascade
);

CREATE INDEX if not exists event_id ON ar_event_actions(event_id);

-- -------------------------------------------------------------------------------------