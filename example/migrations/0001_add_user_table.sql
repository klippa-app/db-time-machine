-- +migrate Up
CREATE TABLE users (
	id INT PRIMARY KEY
);

-- +migrate Down
DROP TABLE users;

