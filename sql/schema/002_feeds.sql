-- +goose Up
create table feeds (
    id varchar(36) not null primary key,
    name varchar(255) not null,
    url text not null unique,
    user_id varchar(256) not null references users(id) on delete cascade,
    created_at timestamp not null,
    updated_at timestamp not null default now()
);


-- +goose Down
drop table feeds;
