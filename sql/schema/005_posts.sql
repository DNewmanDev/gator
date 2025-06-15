-- +goose Up
CREATE TABLE posts(
id  UUID UNIQUE PRIMARY KEY DEFAULT gen_random_uuid(),
created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
title VARCHAR(255) NOT NULL,
url varchar(255) NOT NULL UNIQUE,
description varchar(255),
published_at TIMESTAMP,
feed_id UUID NOT NULL REFERENCES feeds(id) ON DELETE CASCADE

);

-- +goose Down
DROP TABLE posts;


