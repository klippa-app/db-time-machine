-- +migrate Up
ALTER TABLE users DROP COLUMN email;

-- +migrate Down
ALTER TABLE users ADD COLUMN email VARCHAR(256);

