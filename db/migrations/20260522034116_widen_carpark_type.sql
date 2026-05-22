-- +goose Up
-- +goose StatementBegin
ALTER TABLE carparks ALTER COLUMN carpark_type TYPE VARCHAR(50);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE carparks ALTER COLUMN carpark_type TYPE VARCHAR(30);
-- +goose StatementEnd
