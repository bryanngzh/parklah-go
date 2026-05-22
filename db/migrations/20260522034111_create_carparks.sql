-- +goose Up
-- +goose StatementBegin
CREATE TABLE carparks (
    id             SERIAL PRIMARY KEY,
    carpark_code   VARCHAR(10)  NOT NULL,
    carpark_name   TEXT         NOT NULL,
    data_source    VARCHAR(10)  NOT NULL CHECK (data_source IN ('ura', 'hdb')),
    carpark_type   VARCHAR(30),
    parking_system VARCHAR(20)  CHECK (parking_system IN ('electronic', 'coupon')),
    lat            FLOAT,
    lon            FLOAT,
    total_lots     INT,
    created_at     TIMESTAMP    DEFAULT now(),
    updated_at     TIMESTAMP    DEFAULT now(),
    UNIQUE (carpark_code, data_source)
);
CREATE INDEX idx_carparks_code   ON carparks(carpark_code);
CREATE INDEX idx_carparks_source ON carparks(data_source);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_carparks_source;
DROP INDEX IF EXISTS idx_carparks_code;
DROP TABLE carparks;
-- +goose StatementEnd
