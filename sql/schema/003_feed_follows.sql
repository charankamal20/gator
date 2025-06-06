-- +goose Up
create table feed_follows (
    id varchar(36) not null primary key,
    user_id varchar(256) not null references users(id) on delete cascade,
    feed_id varchar(256) not null references feeds(id) on delete cascade,
    created_at timestamp not null,
    updated_at timestamp not null default now(),
    unique (user_id, feed_id)
);


-- +goose Down
drop table feed_follows;
