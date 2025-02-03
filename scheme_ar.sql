CREATE TABLE cron_tasks (
    id          integer not null primary key autoincrement,
    name        text not null,
    description text,
    period      text,
    enabled     integer default 1
);

CREATE TABLE cron_actions (
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

    FOREIGN KEY(task_id) REFERENCES cron_tasks(id) on update cascade on delete cascade
);

CREATE INDEX if not exists task_id ON cron_actions(task_id);

CREATE TABLE events (
    id          integer not null primary key autoincrement,
    target_type text not null,
    target_id   integer not null,
    event_name  text not null,
    enabled     integer default 1
);

CREATE UNIQUE INDEX if not exists tt_ti_en ON events(target_type, target_id, event_name);

CREATE TABLE event_actions (
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

    FOREIGN KEY(event_id) REFERENCES events(id) on update cascade on delete cascade
);

CREATE INDEX if not exists event_id ON event_actions(event_id);

-- -------------------------------------------------------------------------------------

INSERT INTO cron_tasks (id, name, description, period, enabled)
VALUES (1, 'Test htu21d 5s', 'Проверка датчика htu21d', '5s', 1);

INSERT INTO cron_actions (id, task_id, name, type, target_id, target_type, qos, args, enabled, comment)
VALUES (1, 1, 'check', 'method', 82, 'object', 0, '{}', 1, '');

-- -----------------

INSERT INTO events (id, target_id, target_type, event_name, enabled)
VALUES (1, 1, 'object', 'on_click', 1);

INSERT INTO event_actions (id, event_id, name, type, target_id, target_type, qos, args, enabled, comment)
VALUES (1, 1, 'script1', 'script', 0, 'script', 0, '{"x":"1", "y":"2"}', 1, '');

INSERT INTO event_actions (id, event_id, name, type, target_id, target_type, qos, args, enabled, comment)
VALUES (2, 1, 'on', 'method', 3, 'object', 0, '{"x":"1", "y":"2"}', 1, '');

INSERT INTO event_actions (id, event_id, name, type, target_id, target_type, qos, args, enabled, comment)
VALUES (3, 1, 'toggle', 'method', 4, 'object', 0, '{}', 1, '');

-- -----------------

INSERT INTO events (id, target_id, target_type, event_name, enabled)
VALUES (2, 7, 'item', 'on_change_state_on', 1);

INSERT INTO event_actions (id, event_id, name, type, target_id, target_type, qos, args, enabled, comment)
VALUES (4, 2, 'on', 'method', 4, 'object', 0, '{}', 1, '');

-- -----------------

INSERT INTO events (id, target_id, target_type, event_name, enabled)
VALUES (3, 7, 'item', 'on_change_state_off', 1);

INSERT INTO event_actions (id, event_id, target_type, target_id, type, name, args, qos, enabled, comment)
VALUES (5, 3, 'object', 4, 'method', 'off', '{}', 0, 1, '');

-- -----------------

INSERT INTO events (id, target_id, target_type, event_name, enabled)
VALUES (4, 19, 'item', 'on_change_state_on', 1);

INSERT INTO event_actions (id, event_id, target_type, target_id, type, name, args, qos, enabled, comment)
VALUES (6, 4, 'item', 19, 'method', 'set_state', '{"state":"on"}', 0, 1, '');

-- -----------------

INSERT INTO events (id, target_id, target_type, event_name, enabled)
VALUES (5, 16, 'item', 'on_change_state_on', 1);

INSERT INTO event_actions (id, event_id, target_type, target_id, type, name, args, qos, enabled, comment)
VALUES (7, 5, 'object', 48, 'method', 'on', '{}', 0, 1, '');

INSERT INTO event_actions (id, event_id, target_type, target_id, type, name, args, qos, enabled, comment)
VALUES (8, 5, '', 0, 'delay', '', '{"duration":"1s"}', 0, 1, '');

INSERT INTO event_actions (id, event_id, target_type, target_id, type, name, args, qos, enabled, comment)
VALUES (9, 5, 'object', 48, 'method', 'off', '{}', 0, 1, '');

-- -----------------

INSERT INTO events (id, target_id, target_type, event_name, enabled)
VALUES (6, 226, 'item', 'on_change_state_on', 1);

INSERT INTO event_actions (id, event_id, target_type, target_id, type, name, args, qos, enabled, comment)
VALUES (10, 6, 'object', 47, 'method', 'on', '{}', 0, 1, '');

-- -----------------

INSERT INTO events (id, target_id, target_type, event_name, enabled)
VALUES (7, 226, 'item', 'on_change_state_off', 1);

INSERT INTO event_actions (id, event_id, target_type, target_id, type, name, args, qos, enabled, comment)
VALUES (11, 7, 'object', 47, 'method', 'off', '{}', 0, 1, '');

-- -----------------

INSERT INTO events (id, target_id, target_type, event_name, enabled)
VALUES (8, 230, 'item', 'on_click', 1);

INSERT INTO event_actions (id, event_id, target_type, target_id, type, name, args, qos, enabled, comment)
VALUES (12, 8, 'object', 4, 'method', 'on', '{}', 0, 1, '');

-- -----------------

INSERT INTO events (id, target_id, target_type, event_name, enabled)
VALUES (9, 1, 'item', 'on_change_state_on', 1);

INSERT INTO event_actions (id, event_id, target_type, target_id, type, name, args, qos, enabled, comment)
VALUES (13, 9, 'item', 4, 'method', 'set_state', '{"state":"on"}', 0, 1, '');

-- -----------------

INSERT INTO events (id, target_id, target_type, event_name, enabled)
VALUES (10, 1, 'item', 'on_change_state_off', 1);

INSERT INTO event_actions (id, event_id, target_type, target_id, type, name, args, qos, enabled, comment)
VALUES (14, 10, 'item', 4, 'method', 'set_state', '{"state":"off"}', 0, 1, '');
