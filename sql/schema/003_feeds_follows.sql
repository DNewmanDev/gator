-- +goose Up
CREATE TABLE feed_follows(
id  UUID UNIQUE PRIMARY KEY DEFAULT gen_random_uuid(),
created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
feed_id UUID NOT NULL REFERENCES feeds(id) ON DELETE CASCADE
);

ALTER TABLE feed_follows ADD CONSTRAINT unique_user_feed UNIQUE(user_id, feed_id);
-- +goose Down
DROP TABLE feed_follows;




