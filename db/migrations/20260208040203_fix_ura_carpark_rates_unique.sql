-- +goose Up
-- +goose StatementBegin
ALTER TABLE ura_carpark_rates DROP CONSTRAINT ura_carpark_rates_pp_code_vehicle_type_start_time_end_time__key;
ALTER TABLE ura_carpark_rates ADD UNIQUE (pp_code, vehicle_type, start_time, end_time);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE ura_carpark_rates DROP CONSTRAINT ura_carpark_rates_pp_code_vehicle_type_start_time_end_time_key;
ALTER TABLE ura_carpark_rates ADD UNIQUE (pp_code, vehicle_type, start_time, end_time, weekday_min, satday_min, sunph_min);
-- +goose StatementEnd