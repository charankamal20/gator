-- name: CreateFeed :one
insert into feeds (
    id,
    name,
    url,
    user_id,
    created_at,
    updated_at
) values (
    $1,
    $2,
    $3,
    $4,
    now(),
    now()
) returning *;


-- name: GetAllFeeds :many
select f.name, f.url, u.name
from feeds f
join users u
on f.user_id = u.id
order by f.created_at desc;

-- name: GetFeedByUrl :one
select f.id, f.name, f.url, f.user_id, f.created_at, f.updated_at
from feeds f
where f.url = $1;

-- name: MarkFeedFetched :exec
update feeds
set
updated_at = now(),
last_fetched_at = now()
where id = $1;

-- name: GetNextFeedtoFetch :one
select *
from feeds
order by last_fetched_at ASC NULLS FIRST
limit 1;
