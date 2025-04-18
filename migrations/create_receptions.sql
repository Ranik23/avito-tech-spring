-- receptions.sql
--liquibase formatted sql


--changeset anton:create-receptions
CREATE TABLE reception (
    id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    date_time TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    pvz_id INTEGER NOT NULL REFERENCES pvz(id) ON DELETE CASCADE,
    status VARCHAR(20) NOT NULL CHECK (status IN ('open', 'closed'))
);
