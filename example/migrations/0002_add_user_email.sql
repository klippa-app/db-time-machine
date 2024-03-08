-- +migrate Up
ALTER TABLE users ADD COLUMN email VARCHAR(256);

-- +migrate Down
ALTER TABLE users DROP COLUMN email;
