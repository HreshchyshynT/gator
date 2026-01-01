-- name: CreateFeedFollow :one
WITH inserted_feed_follows AS (
  INSERT INTO feed_follows (
    id,
    created_at,
    updated_at,
    user_id,
    feed_id
    ) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
  ) RETURNING *
)
SELECT inserted_feed_follows.*, feeds.name as feed_name, users.name as user_name 
FROM inserted_feed_follows 
INNER JOIN users ON inserted_feed_follows.user_id = users.id 
INNER JOIN feeds ON inserted_feed_follows.feed_id = feeds.id 
;

-- name: GetFeedFollowsForUser :many
select feed_follows.*, users.name as user_name, feeds.name as feed_name
from feed_follows 
INNER JOIN users on users.id = feed_follows.user_id
INNER JOIN feeds on feeds.id = feed_follows.feed_id
where feed_follows.user_id = @user_id;
