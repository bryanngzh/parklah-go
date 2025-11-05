-- +goose Up
-- +goose StatementBegin
CREATE TABLE ura_carpark_season_rates (
    id SERIAL PRIMARY KEY,
    pp_code VARCHAR(5) NOT NULL REFERENCES ura_carparks(pp_code) ON DELETE CASCADE,
    vehicle_type VARCHAR(2),
    monthly_rate NUMERIC(6,2),
    parking_hrs TEXT,
    ticket_type TEXT,
    updated_at TIMESTAMP DEFAULT now(),
    UNIQUE (pp_code, vehicle_type, ticket_type)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE ura_carpark_season_rates;
-- +goose StatementEnd
