-- +goose Up
-- +goose StatementBegin
ALTER TABLE carpark_availability
    ALTER COLUMN snapshot_time TYPE TIMESTAMPTZ USING snapshot_time AT TIME ZONE 'Asia/Singapore',
    ALTER COLUMN created_at    TYPE TIMESTAMPTZ USING created_at    AT TIME ZONE 'UTC';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE carpark_availability
    ALTER COLUMN snapshot_time TYPE TIMESTAMP USING snapshot_time AT TIME ZONE 'UTC',
    ALTER COLUMN created_at    TYPE TIMESTAMP USING created_at    AT TIME ZONE 'UTC';
-- +goose StatementEnd
