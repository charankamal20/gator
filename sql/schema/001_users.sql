-- +goose Up
CREATE TABLE users (
    id  varchar(256),
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    name VARCHAR(256)
);  

-- +goose Down
DROP TABLE users;
