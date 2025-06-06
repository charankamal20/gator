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
