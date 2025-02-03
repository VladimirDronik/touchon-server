-- --------------------------------------------------------------------------------------------------------
-- --------------------------------------------------------------------------------------------------------

create table zones (
    id        INTEGER not null primary key autoincrement,
    parent_id INTEGER not null default 0 references zones on update cascade on delete set default,
    name      TEXT    not null default '',
    style     TEXT    not null default '',
    sort      INTEGER not null default 0
);

-- --------------------------------------------------------------------------------------------------------

create table view_items (
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

-- --------------------------------------------------------------------------------------------------------

create table boilers (
    id                   INTEGER not null primary key autoincrement,
    heating_status       TEXT    not null default '',
    water_status         TEXT    not null default '',
    heating_current_temp REAL    not null default 0,
    heating_optimal_temp REAL    not null default 0,
    water_current_temp   REAL    not null default 0,
    heating_mode         TEXT    not null default '',
    indoor_temp          REAL    not null default 0,
    outdoor_temp         REAL    not null default 0,
    min_threshold        REAL    not null default 0,
    max_threshold        REAL    not null default 0,
    icon                 TEXT    not null default '',
    title                TEXT    not null default '',
    color                TEXT    not null default '',
    auth                 TEXT    not null default ''
);

-- --------------------------------------------------------------------------------------------------------

create table boiler_presets (
    id           INTEGER not null primary key autoincrement,
    boiler_id    INTEGER not null references boilers(id) on update cascade on delete cascade,
    temp_out     REAL    not null default 0,
    temp_coolant REAL    not null default 0
);

-- --------------------------------------------------------------------------------------------------------

create table boiler_properties (
    id         INTEGER not null primary key autoincrement,
    boiler_id  INTEGER not null references boilers(id) on update cascade on delete cascade,
    title      TEXT    not null default '',
    image_name TEXT    not null default '',
    value      TEXT    not null default '',
    status     TEXT    not null default ''
);

-- --------------------------------------------------------------------------------------------------------

create table conditioner_params (
    id                    INTEGER not null primary key autoincrement,
    view_item_id          INTEGER not null unique references view_items on update cascade on delete cascade,

    inside_temp           REAL    not null default 0,
    outside_temp          REAL    not null default 0,
    current_temp          REAL    not null default 0,
    optimal_temp          REAL    not null default 0,
    min_threshold         REAL    not null default 0,
    max_threshold         REAL    not null default 0,

    silent_mode           bool    not null default false,
    eco_mode              bool    not null default false,
    turbo_mode            bool    not null default false,
    sleep_mode            bool    not null default false,

    fan_speeds            JSON    not null default '[]',
    fan_speed             TEXT    not null default '',

    vertical_directions   JSON    not null default '[]',
    vertical_direction    TEXT    not null default '',

    horizontal_directions JSON    not null default '[]',
    horizontal_direction  TEXT    not null default '',

    operating_modes       JSON    not null default '[]',
    operating_mode        TEXT    not null default '',

    ionisation            bool    not null default false,
    self_cleaning         bool    not null default false,
    anti_mold             bool    not null default false,
    sound                 bool    not null default false,
    on_duty_heating       bool    not null default false,
    soft_top              bool    not null default false
);

-- --------------------------------------------------------------------------------------------------------

create table counters (
    id             INTEGER not null primary key autoincrement,
    name           TEXT    not null default '',
    type           TEXT    not null default '',
    unit           TEXT    not null default '',
    today_value    REAL    not null default 0,
    week_value     REAL    not null default 0,
    month_value    REAL    not null default 0,
    year_value     REAL    not null default 0,
    price_for_unit REAL    not null default 0,
    impulse        REAL    not null default 0,
    sort           INTEGER not null default 0,
    enabled        bool    not null default true
);

-- --------------------------------------------------------------------------------------------------------

create table counter_daily_history (
    id         INTEGER  not null primary key autoincrement,
    counter_id INTEGER  not null references counters on update cascade on delete cascade,
    datetime   TEXT     not null default '',
    value      REAL     not null default 0,

    unique (counter_id, datetime)
);

-- --------------------------------------------------------------------------------------------------------

create table counter_monthly_history (
    id         INTEGER not null primary key autoincrement,
    counter_id INTEGER not null references counters on update cascade on delete cascade,
    datetime   TEXT    not null default '',
    value      REAL    not null default 0,

    unique (counter_id, datetime)
);

-- --------------------------------------------------------------------------------------------------------

create table curtain_params (
    id           INTEGER not null primary key autoincrement,
    view_item_id INTEGER not null unique references view_items on update cascade on delete cascade,
    type         TEXT    not null default '',
    control_type TEXT    not null default '',
    open_percent REAL    not null default 0
);

-- --------------------------------------------------------------------------------------------------------

create table device_daily_history (
    id           INTEGER not null primary key autoincrement,
    view_item_id INTEGER not null references view_items on update cascade on delete cascade,
    datetime     TEXT    not null default '',
    value        REAL    not null default 0,

    unique (view_item_id, datetime)
);

-- --------------------------------------------------------------------------------------------------------

create table device_hourly_history (
    id           INTEGER not null primary key autoincrement,
    view_item_id INTEGER not null references view_items on update cascade on delete cascade,
    datetime     TEXT    not null default '',
    value        REAL    not null default 0,

    unique (view_item_id, datetime)
);

-- --------------------------------------------------------------------------------------------------------

create table dimmers (
    id           INTEGER not null primary key autoincrement,
    view_item_id INTEGER not null unique references view_items on update cascade on delete cascade,
    name         TEXT    not null default '',
    value        INTEGER not null default 0,
    enabled      bool    not null default false
);

-- --------------------------------------------------------------------------------------------------------

create table events (
    id          INTEGER not null primary key autoincrement,
    target_type TEXT    not null default '',
    target_id   INTEGER not null default 0,
    event       TEXT    not null default '',
    value       TEXT    not null default ''
);

-- --------------------------------------------------------------------------------------------------------

create table light_params (
    id           INTEGER not null primary key autoincrement,
    view_item_id INTEGER not null unique references view_items on update cascade on delete cascade,
    hue          INTEGER not null default 0,
    saturation   REAL    not null default 0,
    brightness   REAL    not null default 0,
    cct          INTEGER not null default 0
);

-- --------------------------------------------------------------------------------------------------------

create table local_users (
    id       INTEGER not null primary key autoincrement ,
    name     TEXT    not null default '',
    password TEXT    not null default ''
);

INSERT INTO local_users (id, name, password) VALUES (1, 'web', '12345');

-- --------------------------------------------------------------------------------------------------------

create table menus (
    id        INTEGER not null primary key autoincrement,
    parent_id INTEGER not null default 0 references menus on update cascade on delete cascade,
    page      TEXT    not null default '',
    title     TEXT    not null default '',
    image     TEXT    not null default '',
    sort      INTEGER not null default 0,
    params    JSON    not null default '{}',
    enabled   bool    not null default true
);

-- --------------------------------------------------------------------------------------------------------

create table notifications (
    id      INTEGER not null primary key autoincrement,
    type    TEXT    not null default '',
    date    TEXT    not null default '',
    text    TEXT    not null default '',
    is_read bool    not null default false
);

-- --------------------------------------------------------------------------------------------------------

create table scenarios (
    id           INTEGER not null primary key autoincrement,
    view_item_id INTEGER not null unique references view_items on update cascade on delete cascade,
    type         TEXT    not null default '',
    description  TEXT    not null default '',
    icon         TEXT    not null default '',
    title        TEXT    not null default '',
    sort         INTEGER not null default 0,
    color        TEXT    not null default '',
    auth         TEXT    not null default '',
    enabled      bool    not null default false
);

-- --------------------------------------------------------------------------------------------------------

create table sensors (
    id            INTEGER not null primary key autoincrement,
    view_item_id  INTEGER not null unique references view_items on update cascade on delete cascade,
    zone_id       INTEGER not null references zones on update cascade on delete set default,
    type          TEXT    not null default '',
    name          TEXT    not null default '',
    current       REAL    not null default 0,
    optimal       REAL    not null default 0,
    min_threshold REAL    not null default 0,
    max_threshold REAL    not null default 0,
    icon          TEXT    not null default '',
    position_left INTEGER not null default 0,
    position_top  INTEGER not null default 0,
    sort          INTEGER not null default 0,
    auth          TEXT    not null default '',
    enabled       BOOLEAN not null default false
);

-- --------------------------------------------------------------------------------------------------------

create table temp_presets (
    id      INTEGER not null primary key autoincrement,
    zone_id INTEGER not null unique references zones on update cascade on delete cascade,
    normal  REAL    not null default 0,
    night   REAL    not null default 0,
    eco     REAL    not null default 0,
    sort    INTEGER not null default 0
);

-- --------------------------------------------------------------------------------------------------------

create table users (
    id            INTEGER  not null primary key autoincrement,
    login         TEXT     not null default '',
    password      TEXT     not null default '',
    role          TEXT     not null default '',
    send_push     bool     not null default true,
    refresh_token TEXT     not null default '',
    token_expired datetime not null default '',
    device_id     INTEGER  not null default 0,
    device_type   TEXT     not null default '',
    device_token  TEXT     not null default '',
    comment       TEXT     not null default ''
);

-- --------------------------------------------------------------------------------------------------------
-- --------------------------------------------------------------------------------------------------------

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

-- --------------------------------------------------------------------------------------------------------

CREATE TABLE om_props (
    id        INTEGER not null primary key autoincrement,
    object_id INTEGER not null,
    code      TEXT not null,
    value     TEXT,

    FOREIGN KEY(object_id) REFERENCES om_objects(id) on update cascade on delete cascade
);

create unique index object_id_code on om_props (object_id, code);

-- --------------------------------------------------------------------------------------------------------

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

-- ---------------------------------------------------------------------------------------------------------
-- --------------------------------------------------------------------------------------------------------

create table ar_cron_tasks
(
    id          INTEGER not null primary key autoincrement,
    name        TEXT not null,
    description TEXT,
    period      TEXT,
    enabled     INTEGER default 1
);

-- --------------------------------------------------------------------------------------------------------

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

-- --------------------------------------------------------------------------------------------------------

CREATE TABLE ar_events (
    id          integer not null primary key autoincrement,
    target_type text not null,
    target_id   integer not null,
    event_name  text not null,
    enabled     integer default 1
);

CREATE UNIQUE INDEX if not exists tt_ti_en ON ar_events(target_type, target_id, event_name);

-- --------------------------------------------------------------------------------------------------------

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

-- --------------------------------------------------------------------------------------------------------
-- --------------------------------------------------------------------------------------------------------
