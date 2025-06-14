// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: feeds.sql

package database

import (
	"context"
	"database/sql"
	"time"
)

const createFeed = `-- name: CreateFeed :one
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
) returning id, name, url, user_id, created_at, updated_at, last_fetched_at
`

type CreateFeedParams struct {
	ID     string
	Name   string
	Url    string
	UserID string
}

func (q *Queries) CreateFeed(ctx context.Context, arg CreateFeedParams) (Feed, error) {
	row := q.db.QueryRowContext(ctx, createFeed,
		arg.ID,
		arg.Name,
		arg.Url,
		arg.UserID,
	)
	var i Feed
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Url,
		&i.UserID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.LastFetchedAt,
	)
	return i, err
}

const getAllFeeds = `-- name: GetAllFeeds :many
select f.name, f.url, u.name
from feeds f
join users u
on f.user_id = u.id
order by f.created_at desc
`

type GetAllFeedsRow struct {
	Name   string
	Url    string
	Name_2 sql.NullString
}

func (q *Queries) GetAllFeeds(ctx context.Context) ([]GetAllFeedsRow, error) {
	rows, err := q.db.QueryContext(ctx, getAllFeeds)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetAllFeedsRow
	for rows.Next() {
		var i GetAllFeedsRow
		if err := rows.Scan(&i.Name, &i.Url, &i.Name_2); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getFeedByUrl = `-- name: GetFeedByUrl :one
select f.id, f.name, f.url, f.user_id, f.created_at, f.updated_at
from feeds f
where f.url = $1
`

type GetFeedByUrlRow struct {
	ID        string
	Name      string
	Url       string
	UserID    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (q *Queries) GetFeedByUrl(ctx context.Context, url string) (GetFeedByUrlRow, error) {
	row := q.db.QueryRowContext(ctx, getFeedByUrl, url)
	var i GetFeedByUrlRow
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Url,
		&i.UserID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getNextFeedtoFetch = `-- name: GetNextFeedtoFetch :one
select id, name, url, user_id, created_at, updated_at, last_fetched_at
from feeds
order by last_fetched_at ASC NULLS FIRST
limit 1
`

func (q *Queries) GetNextFeedtoFetch(ctx context.Context) (Feed, error) {
	row := q.db.QueryRowContext(ctx, getNextFeedtoFetch)
	var i Feed
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Url,
		&i.UserID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.LastFetchedAt,
	)
	return i, err
}

const markFeedFetched = `-- name: MarkFeedFetched :exec
update feeds
set
updated_at = now(),
last_fetched_at = now()
where id = $1
`

func (q *Queries) MarkFeedFetched(ctx context.Context, id string) error {
	_, err := q.db.ExecContext(ctx, markFeedFetched, id)
	return err
}
