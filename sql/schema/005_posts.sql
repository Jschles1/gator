-- +goose Up
CREATE TABLE posts (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  feed_id uuid NOT NULL,
  title VARCHAR NOT NULL,
  description VARCHAR,
  url VARCHAR NOT NULL,
  published_at TIMESTAMP NOT NULL,
  FOREIGN KEY (feed_id) REFERENCES feeds(id) ON DELETE CASCADE,
  UNIQUE(url)
);

-- +goose Down
DROP TABLE posts;