-- pvzs.sql
--liquibase formatted sql


--changeset anton:create-pvzs
CREATE TABLE pvz (
    id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    registration_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    city VARCHAR(50) NOT NULL
);