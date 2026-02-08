-- +goose Up
-- +goose StatementBegin
CREATE TABLE carpark_metadata (
    id SERIAL PRIMARY KEY,
    carpark_code VARCHAR(10) NOT NULL,
    car_park_decks INT,
    gantry_height NUMERIC(4,2),
    car_park_basement BOOLEAN,
    data_source VARCHAR(20) NOT NULL,
    updated_at TIMESTAMP DEFAULT now(),
    UNIQUE(carpark_code, data_source),
    FOREIGN KEY (carpark_code, data_source) REFERENCES carparks(carpark_code, data_source) ON DELETE CASCADE
);

CREATE INDEX idx_metadata_source ON carpark_metadata(data_source);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_metadata_source;
DROP TABLE carpark_metadata;
-- +goose StatementEnd
