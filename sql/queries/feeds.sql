-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetAllFeeds :many
SELECT * FROM feeds ORDER BY created_at DESC;

-- name: GetNextFeedsToFetch :many
SELECT * FROM feeds
WHERE last_fetched_at IS NULL OR last_fetched_at < NOW() - INTERVAL '1 hour'
ORDER BY COALESCE(last_fetched_at, '1970-01-01'::timestamp) ASC
LIMIT $1;

-- name: MarkFeedFetched :one
UPDATE feeds
SET last_fetched_at = NOW(),
    updated_at = NOW()
WHERE id = $1
RETURNING *;
