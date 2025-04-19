-- +goose Up
-- +goose StatementBegin
CREATE TABLE pvz (
    id SERIAL PRIMARY KEY,
    registration_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    city VARCHAR(50) NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS pvz;
-- +goose StatementEnd
