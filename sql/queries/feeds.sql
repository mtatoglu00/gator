-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES (
  $1,
  $2,
  $3,
  $4,
  $5,
  $6
  )
  RETURNING *;

-- name: GetFeeds :many
SELECT * from feeds;

-- name: CreateFeedFollow :one
WITH inserted_feed_follow AS (
INSERT INTO feed_follows (created_at, updated_at, user_id, feed_id)
VALUES (
  $1,
  $2,
  $3,
  $4
  )
  RETURNING *
)
SELECT
  inserted_feed_follow.*,
  feeds.name AS feed_name,
  users.name AS user_name
FROM inserted_feed_follow
INNER JOIN users ON  inserted_feed_follow.user_id = users.id
INNER JOIN feeds ON inserted_feed_follow.feed_id = feeds.id;

-- name: GetFeedByURL :one
SELECT * FROM feeds 
WHERE url = $1;

-- name: GetFeedFollowsForUser :many
SELECT users.name AS User_Name, feeds.name AS Feed_Name, feeds.url AS url FROM users
INNER JOIN feed_follows ON users.id = feed_follows.user_id
INNER JOIN feeds ON feeds.id = feed_follows.feed_id
WHERE users.name = $1;
