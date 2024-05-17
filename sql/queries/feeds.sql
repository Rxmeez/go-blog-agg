-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetAllFeeds :many
SELECT id, created_at, updated_at, name, url, user_id
FROM feeds;

-- name: FeedIdExists :one
SELECT id
FROM feeds
WHERE id = $1;