-- +goose Up
-- +goose StatementBegin
CREATE TABLE carpark_availability (
    id SERIAL PRIMARY KEY,
    carpark_code VARCHAR(10) NOT NULL REFERENCES carparks(carpark_code) ON DELETE CASCADE,
    vehicle_type CHAR(1),
    lots_available INT,
    total_lots INT,
    data_source VARCHAR(20),
    snapshot_time TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT now(),
    INDEX idx_code_time (carpark_code, snapshot_time DESC)
);

CREATE INDEX idx_availability_data_source ON carpark_availability(data_source);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_availability_data_source;
DROP INDEX IF EXISTS idx_code_time;
DROP TABLE carpark_availability;
-- +goose StatementEnd
