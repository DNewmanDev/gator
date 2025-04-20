-- name: CreateFeedFollow :many
WITH inserted_feed_follow AS(

    INSERT INTO feed_follows()VALUES()RETURNING *
)
SELECT 
inserted_feed_follow.*, feeds.name AS feed_name, users_name AS user_name
FROM inserted_feed_follow
INNER JOIN
INNER JOIN