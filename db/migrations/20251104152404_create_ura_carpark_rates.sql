-- +goose Up
-- +goose StatementBegin
CREATE TABLE ura_carpark_rates (
    id SERIAL PRIMARY KEY,
    pp_code VARCHAR(5) NOT NULL REFERENCES ura_carparks(pp_code) ON DELETE CASCADE,
    vehicle_type VARCHAR(2),
    start_time TIME,
    end_time TIME,
    weekday_min INTEGER,
    weekday_rate NUMERIC(6,2),
    satday_min INTEGER,
    satday_rate NUMERIC(6,2),
    sunph_min INTEGER,
    sunph_rate NUMERIC(6,2),
    updated_at TIMESTAMP DEFAULT now(),
    UNIQUE (pp_code, vehicle_type, start_time, end_time, weekday_min, satday_min, sunph_min)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE ura_carpark_rates;
-- +goose StatementEnd
