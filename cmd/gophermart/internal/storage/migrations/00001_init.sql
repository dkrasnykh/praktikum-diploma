-- +goose Up

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS users
(
    id            SERIAL       NOT NULL UNIQUE,
    login         VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    created_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS rewards
(
    id           SERIAL PRIMARY KEY,
    user_id      SERIAL,
    order_number VARCHAR NOT NULL UNIQUE,
    status       VARCHAR,
    accrual      BIGINT,
    withdraw     BIGINT,
    updated_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE users;
DROP TABLE rewards;
DROP EXTENSION IF EXISTS "uuid-ossp";