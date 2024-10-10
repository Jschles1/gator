-- name: CreatePost :one
INSERT INTO posts (id, feed_id, created_at, updated_at, published_at, title, url, description)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8
)
RETURNING *;

-- name: GetPostsForUser :many
SELECT
    p.*
FROM posts p
INNER JOIN feed_follows ff ON p.feed_id = ff.feed_id
WHERE ff.user_id = $1
ORDER BY p.published_at DESC
LIMIT $2;