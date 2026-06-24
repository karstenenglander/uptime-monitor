-- +goose Up
ALTER TABLE sites
ADD COLUMN last_status_code INT,
ADD COLUMN latency BIGINT;

-- +goose Down
ALTER TABLE sites
DROP COLUMN last_status_code,
DROP COLUMN latency;
