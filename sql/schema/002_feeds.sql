-- +goose Up

CREATE EXTENSION IF NOT EXISTS "pgcrypto";
CREATE TABLE feeds(
id  UUID UNIQUE PRIMARY KEY DEFAULT gen_random_uuid(),
created_at TIMESTAMP NOT NULL DEFAULT now(),
updated_at TIMESTAMP NOT NULL DEFAULT now(),
name TEXT NOT NULL,
url TEXT UNIQUE NOT NULL,
user_id UUID REFERENCES users(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE feeds;




