-- name: ListFeedsWithUsers :many

SELECT 
feeds.name AS feed_name,
feeds.url AS feeds_url,
users.name AS user_name
FROM feeds
JOIN users ON feeds.user_id = users.id;