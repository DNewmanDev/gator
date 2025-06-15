-- name: MarkFeedFetched :one


UPDATE feeds
SET last_fetched_at = now(),
updated_at = now()
WHERE feeds.id = $1
RETURNING *;


