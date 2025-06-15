-- name: DeleteFeedFollow :exec
DELETE FROM feed_follows ff WHERE ff.user_id = $1 AND ff.feed_id =(
    SELECT feeds.id FROM feeds WHERE feeds.url = $2
); 