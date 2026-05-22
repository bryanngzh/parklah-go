-- +goose Up
-- +goose StatementBegin
CREATE TABLE carpark_short_term_rates (
    id             SERIAL PRIMARY KEY,
    carpark_code   VARCHAR(10)  NOT NULL,
    data_source    VARCHAR(10)  NOT NULL,
    vehicle_type   CHAR(1)      NOT NULL CHECK (vehicle_type IN ('C', 'M', 'H')),
    day_type       VARCHAR(15)  NOT NULL CHECK (day_type IN ('weekday', 'saturday', 'sunday_ph', 'all')),
    start_time     TIME,
    end_time       TIME,
    rate_per_30min NUMERIC(5,2) NOT NULL,
    min_duration   VARCHAR(20),
    updated_at     TIMESTAMP    DEFAULT now(),
    UNIQUE (carpark_code, data_source, vehicle_type, day_type, start_time),
    FOREIGN KEY (carpark_code, data_source) REFERENCES carparks(carpark_code, data_source) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE carpark_short_term_rates;
-- +goose StatementEnd
