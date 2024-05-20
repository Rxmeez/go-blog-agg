-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetAllFeeds :many
SELECT id, created_at, updated_at, name, url, user_id, last_fetched_at
FROM feeds;

-- name: FeedIdExists :one
SELECT id
FROM feeds
WHERE id = $1;

-- name: GetNextFeedsToFetch :many
SELECT id, url
FROM feeds
ORDER BY last_fetched_at ASC NULLS FIRST
Limit $1;

-- name: MarkFeedFetched :exec
UPDATE feeds
SET last_fetched_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP
WHERE id = $1;
