-- name: CreatePost :one
insert into posts
(
    id,
    title,
    description,
    url,
    feed_id,
    published_at,
    created_at,
    updated_at
)
values
(
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    now(),
    now()
)
returning *;

-- name: GetPostsForUser :many
select *
from posts p
join feeds f on p.feed_id = f.id
join users u on f.user_id = u.id
where u.id = $1
order by p.published_at desc
limit $2;
