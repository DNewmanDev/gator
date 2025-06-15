-- name: GetFeedFollowsForUser :many

SELECT ff.*, f.name AS feed_name, u.name AS user_name FROM feed_follows ff
INNER JOIN feeds f ON ff.feed_id=f.id
INNER JOIN users u ON ff.user_id=u.id
WHERE ff.user_id = $1;
