-- +goose Up
create table posts (
    id varchar(256) not null primary key,
    title varchar(256) not null,
    url varchar(256) not null unique,
    description text not null,
    published_at timestamp not null,
    feed_id varchar(256) not null references feeds(id) on delete cascade,
    created_at timestamp not null default now(),
    updated_at timestamp not null default now()
);


-- +goose Down
drop table posts;
