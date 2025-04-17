-- products.sql
--liquibase formatted sql


--changeset anton:create-products
CREATE TABLE product (
    id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    date_time TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    type VARCHAR(50) NOT NULL,
    reception_id INTEGER NOT NULL REFERENCES reception(id) ON DELETE CASCADE
);