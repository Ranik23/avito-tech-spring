-- +goose Up
-- +goose StatementBegin
CREATE TABLE product (
    id SERIAL PRIMARY KEY,
    date_time TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    type VARCHAR(50) NOT NULL,
    reception_id INTEGER NOT NULL REFERENCES reception(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS product;
-- +goose StatementEnd
