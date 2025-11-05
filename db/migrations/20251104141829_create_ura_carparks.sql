-- +goose Up
-- +goose StatementBegin
CREATE TABLE ura_carparks (
    id SERIAL PRIMARY KEY,
    pp_code VARCHAR(5) UNIQUE NOT NULL,
    pp_name TEXT NOT NULL,
    parking_system VARCHAR(2),
    geom POINT,
    updated_at TIMESTAMP DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE ura_carparks;
-- +goose StatementEnd
