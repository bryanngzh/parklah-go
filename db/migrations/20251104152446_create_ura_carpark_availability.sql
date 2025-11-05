-- +goose Up
-- +goose StatementBegin
CREATE TABLE ura_carpark_availability (
    id SERIAL PRIMARY KEY,
    pp_code VARCHAR(5) NOT NULL REFERENCES ura_carparks(pp_code) ON DELETE CASCADE,
    vehicle_type VARCHAR(2),
    lots_available INT,
    updated_at TIMESTAMP DEFAULT now(),
    UNIQUE (pp_code, vehicle_type)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE ura_carpark_availability;
-- +goose StatementEnd
