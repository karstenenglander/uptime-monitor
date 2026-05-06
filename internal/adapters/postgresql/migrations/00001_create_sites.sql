-- +goose Up
CREATE TABLE IF NOT EXISTS sites (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR NOT NULL,
    url VARCHAR NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    polled_at TIMESTAMPZ
);

-- +goose Down
DROP TABLE IF EXISTS sites;
