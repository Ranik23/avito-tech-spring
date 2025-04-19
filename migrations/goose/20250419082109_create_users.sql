-- +goose Up
-- +goose StatementBegin
CREATE TABLE reception (
    id SERIAL PRIMARY KEY,
    date_time TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    pvz_id INTEGER NOT NULL REFERENCES pvz(id) ON DELETE CASCADE,
    status VARCHAR(20) NOT NULL CHECK (status IN ('open', 'closed'))
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS reception;
-- +goose StatementEnd
