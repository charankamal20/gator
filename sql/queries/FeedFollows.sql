-- name: CreateFeedFollow :one
with inserted as (
    insert into feed_follows (
        id,
        feed_id,
        user_id,
        created_at
    ) values (
        $1,
        $2,
        $3,
        now()
    )
    returning id, feed_id, user_id, created_at
)
select
    inserted.id,
    inserted.feed_id,
    inserted.user_id,
    inserted.created_at,
    feeds.name as feed_name,
    users.name as user_name
from inserted
join feeds on feeds.id = inserted.feed_id
join users on users.id = inserted.user_id;


-- name: GetFeedFollowsForUser :many
select
    feed_follows.id,
    feed_follows.feed_id,
    feed_follows.user_id,
    feed_follows.created_at,
    feeds.name as feed_name,
    users.name as user_name
from feed_follows
join feeds on feeds.id = feed_follows.feed_id
join users on users.id = feed_follows.user_id
where feed_follows.user_id = $1;
