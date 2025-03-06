-- +goose Up
-- +goose StatementBegin
drop table counters;
drop table counter_daily_history;
drop table counter_monthly_history;
drop table device_daily_history;
drop table device_hourly_history;

create table hourly_history (
    id  INTEGER not null primary key autoincrement,
    object_id     INTEGER not null references om_objects on update cascade on delete cascade,
    datetime TEXT    not null default '',
    value TEXT not null default ''
);

create table daily_history (
    id  INTEGER not null primary key autoincrement,
    object_id     INTEGER not null references om_objects on update cascade on delete cascade,
    datetime TEXT    not null default '',
    value TEXT not null default ''
);

create table monthly_history (
    id  INTEGER not null primary key autoincrement,
    object_id     INTEGER not null references om_objects on update cascade on delete cascade,
    datetime TEXT    not null default '',
    value TEXT not null default ''
);

-- +goose StatementEnd

