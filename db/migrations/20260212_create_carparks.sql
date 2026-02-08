-- +goose Up
-- +goose StatementBegin
CREATE TABLE carparks (
    id SERIAL PRIMARY KEY,
    carpark_code VARCHAR(10) NOT NULL,
    carpark_name TEXT NOT NULL,
    data_source VARCHAR(20) NOT NULL,
    carpark_type VARCHAR(50),
    parking_system VARCHAR(50),
    location_x FLOAT,
    location_y FLOAT,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now(),
    UNIQUE(carpark_code, data_source)
);

CREATE INDEX idx_carpark_code ON carparks(carpark_code);
CREATE INDEX idx_data_source ON carparks(data_source);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_data_source;
DROP INDEX IF EXISTS idx_carpark_code;
DROP TABLE carparks;
-- +goose StatementEnd
