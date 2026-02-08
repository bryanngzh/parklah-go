-- +goose Up
-- +goose StatementBegin
CREATE TABLE carpark_availability (
    id SERIAL PRIMARY KEY,
    carpark_code VARCHAR(10) NOT NULL,
    vehicle_type CHAR(1),
    lots_available INT,
    total_lots INT,
    data_source VARCHAR(20) NOT NULL,
    snapshot_time TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT now(),
    FOREIGN KEY (carpark_code, data_source) REFERENCES carparks(carpark_code, data_source) ON DELETE CASCADE
);

CREATE INDEX idx_code_time ON carpark_availability(carpark_code, snapshot_time DESC);

CREATE INDEX idx_availability_data_source ON carpark_availability(data_source);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_availability_data_source;
DROP INDEX IF EXISTS idx_code_time;
DROP TABLE carpark_availability;
-- +goose StatementEnd
