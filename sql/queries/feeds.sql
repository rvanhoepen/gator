-- name: CreateFeed :one
WITH inserted_feed AS (
  INSERT INTO feeds (id, name, url, user_id, created_at, updated_at)
  VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
  )
  RETURNING *
),
inserted_feed_follow AS (
  INSERT INTO feed_follows (user_id, feed_id)
  SELECT
    inserted_feed.user_id,
    inserted_feed.id
  FROM inserted_feed
)
SELECT *
FROM inserted_feed;

-- name: GetFeeds :many
SELECT f.id, f.name, f.url, u.name as created_by, f.created_at, f.updated_at
FROM feeds f
JOIN users u ON f.user_id = u.id
ORDER BY f.name;

-- name: GetFeedByUrl :one
SELECT f.id, f.name, f.url, u.name as created_by, f.created_at, f.updated_at
FROM feeds f
JOIN users u ON f.user_id = u.id
WHERE f.url = $1;

