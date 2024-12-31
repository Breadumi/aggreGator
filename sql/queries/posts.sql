-- name: CreatePost :one
INSERT INTO posts (created_at, updated_at, title, url, description, published_at, feed_id)
VALUES($1, $2, $3, $4, $5, $6, $7)
ON CONFLICT (url) DO NOTHING
RETURNING *;

-- name: GetPostsByUser :many
SELECT * FROM posts
WHERE feed_id IN (
    SELECT feed_id FROM feed_follows
    WHERE feed_follows.user_id = $1
)
LIMIT $2;