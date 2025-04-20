-- name: CreateFeedFollow :many
WITH inserted_feed_follow AS(

    INSERT INTO feed_follows (user_id, feed_id)
    VALUES($1, $2) RETURNING *
)
SELECT 
ff.*, f.name AS feed_name, u.name AS user_name
FROM inserted_feed_follow ff
INNER JOIN feeds f on ff.feed_id = f.id 
INNER JOIN users u on ff.user_id = u.id;



-- come look at this one later