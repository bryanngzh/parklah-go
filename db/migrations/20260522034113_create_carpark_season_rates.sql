-- +goose Up
-- +goose StatementBegin
CREATE TABLE carpark_season_rates (
    id           SERIAL PRIMARY KEY,
    carpark_code VARCHAR(10)  NOT NULL,
    data_source  VARCHAR(10)  NOT NULL,
    vehicle_type CHAR(1)      NOT NULL CHECK (vehicle_type IN ('C', 'M', 'H')),
    ticket_type  VARCHAR(20)  NOT NULL CHECK (ticket_type IN ('Commercial', 'Residential')),
    parking_hrs  TEXT,
    monthly_rate NUMERIC(8,2) NOT NULL,
    updated_at   TIMESTAMP    DEFAULT now(),
    UNIQUE (carpark_code, data_source, vehicle_type, ticket_type),
    FOREIGN KEY (carpark_code, data_source) REFERENCES carparks(carpark_code, data_source) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE carpark_season_rates;
-- +goose StatementEnd
