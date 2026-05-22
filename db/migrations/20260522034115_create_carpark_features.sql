-- +goose Up
-- +goose StatementBegin
CREATE TABLE carpark_features (
    id                   SERIAL PRIMARY KEY,
    carpark_code         VARCHAR(10)  NOT NULL,
    data_source          VARCHAR(10)  NOT NULL,
    short_term_parking   VARCHAR(50),
    free_parking         VARCHAR(100),
    night_parking        BOOLEAN,
    car_park_decks       INT,
    gantry_height        NUMERIC(4,2),
    car_park_basement    BOOLEAN,
    is_central_area      BOOLEAN   DEFAULT false,
    is_peak_hour_carpark BOOLEAN   DEFAULT false,
    updated_at           TIMESTAMP DEFAULT now(),
    UNIQUE (carpark_code, data_source),
    FOREIGN KEY (carpark_code, data_source) REFERENCES carparks(carpark_code, data_source) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE carpark_features;
-- +goose StatementEnd
