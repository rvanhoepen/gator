-- name: CreatePost :one
INSERT INTO posts (
  id,
  title,
  url,
  description,
  published_at,
  feed_id
) VALUES (
  $1,
  $2,
  $3,
  $4,
  $5,
  $6
)
ON CONFLICT (url) DO NOTHING
RETURNING id, title, url, description, published_at, feed_id;

-- name: GetPostsForUser :many
SELECT id, title, url, description, published_at, feed_id
FROM posts
WHERE feed_id in (
  SELECT feed_id
  FROM feed_follows
  WHERE user_id = $1
)
ORDER BY published_at DESC
LIMIT $2;
