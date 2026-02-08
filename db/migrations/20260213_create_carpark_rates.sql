-- +goose Up
-- +goose StatementBegin
CREATE TABLE carpark_rates (
    id SERIAL PRIMARY KEY,
    carpark_code VARCHAR(10) NOT NULL,
    vehicle_type CHAR(1),
    parking_hrs TEXT,
    rate NUMERIC(6,2),
    data_source VARCHAR(20) NOT NULL,
    updated_at TIMESTAMP DEFAULT now(),
    UNIQUE(carpark_code, vehicle_type, parking_hrs, data_source),
    FOREIGN KEY (carpark_code, data_source) REFERENCES carparks(carpark_code, data_source) ON DELETE CASCADE
);

CREATE INDEX idx_carpark_rates_code ON carpark_rates(carpark_code);
CREATE INDEX idx_carpark_rates_vehicle ON carpark_rates(vehicle_type);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_carpark_rates_vehicle;
DROP INDEX IF EXISTS idx_carpark_rates_code;
DROP TABLE carpark_rates;
-- +goose StatementEnd
