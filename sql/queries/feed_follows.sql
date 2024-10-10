-- name: CreateFeedFollow :one
WITH inserted_follow AS (
    INSERT INTO feed_follows (id, user_id, feed_id, created_at, updated_at)
    VALUES ($1, $2, $3, $4, $5)
    RETURNING *
)
SELECT 
    f.*,
    u.name AS user_name,
    fd.name AS feed_name
FROM inserted_follow f
INNER JOIN users u ON f.user_id = u.id
INNER JOIN feeds fd ON f.feed_id = fd.id;

-- name: GetFeedFollowsForUser :many
SELECT
    f.*,
    u.name AS user_name,
    fd.name AS feed_name
FROM feed_follows f
INNER JOIN users u ON f.user_id = u.id
INNER JOIN feeds fd ON f.feed_id = fd.id
WHERE f.user_id = $1;

-- name: Unfollow :exec
DELETE FROM feed_follows
WHERE feed_follows.user_id = $1 AND feed_follows.feed_id = (
    SELECT id FROM feeds WHERE url = $2
);
