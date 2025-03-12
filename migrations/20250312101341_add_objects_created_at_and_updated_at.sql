-- +goose Up
-- +goose StatementBegin
alter table om_objects add created_at datetime;
alter table om_objects add updated_at datetime;

CREATE TRIGGER om_objects_created_at AFTER INSERT ON om_objects FOR EACH ROW
BEGIN
    UPDATE om_objects SET created_at = current_timestamp WHERE id = new.id;
END;

CREATE TRIGGER om_objects_updated_at AFTER UPDATE ON om_objects FOR EACH ROW
BEGIN
    UPDATE om_objects SET updated_at = current_timestamp WHERE id = new.id;
END;

update om_objects set created_at = current_timestamp, updated_at = current_timestamp;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
