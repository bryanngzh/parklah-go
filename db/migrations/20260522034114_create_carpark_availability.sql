-- +goose Up
-- +goose StatementBegin
CREATE TABLE carpark_availability (
    id             SERIAL PRIMARY KEY,
    carpark_code   VARCHAR(10) NOT NULL,
    data_source    VARCHAR(10) NOT NULL CHECK (data_source IN ('ura', 'hdb')),
    vehicle_type   CHAR(1)     NOT NULL CHECK (vehicle_type IN ('C', 'M', 'H')),
    lots_available INT         NOT NULL,
    total_lots     INT,
    snapshot_time  TIMESTAMP   NOT NULL,
    created_at     TIMESTAMP   DEFAULT now()
);
CREATE INDEX idx_availability_code_time ON carpark_availability(carpark_code, snapshot_time DESC);
CREATE INDEX idx_availability_source    ON carpark_availability(data_source);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_availability_source;
DROP INDEX IF EXISTS idx_availability_code_time;
DROP TABLE carpark_availability;
-- +goose StatementEnd
