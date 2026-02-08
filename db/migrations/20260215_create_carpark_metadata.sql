-- +goose Up
-- +goose StatementBegin
CREATE TABLE carpark_metadata (
    id SERIAL PRIMARY KEY,
    carpark_code VARCHAR(10) NOT NULL UNIQUE REFERENCES carparks(carpark_code) ON DELETE CASCADE,
    car_park_decks INT,
    gantry_height NUMERIC(4,2),
    car_park_basement BOOLEAN,
    data_source VARCHAR(20),
    updated_at TIMESTAMP DEFAULT now()
);

CREATE INDEX idx_metadata_source ON carpark_metadata(data_source);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_metadata_source;
DROP TABLE carpark_metadata;
-- +goose StatementEnd
