-- name: CreateFeed :one 
insert into feeds (id, created_at, updated_at, name, url, user_id) 
values (
  $1,
  $2,
  $3,
  $4,
  $5,
  $6
) RETURNING *;

-- name: GetAllFeeds :many
select feeds.*, users.name as userName from feeds 
inner join users on feeds.user_id = users.id;

-- name: GetFeedByUrl :one
select * from feeds
where feeds.url = @url;

-- name: MarkFeedFetched :one
UPDATE feeds 
SET updated_at = NOW(), last_fetched_at = NOW()
WHERE id = @feed_id RETURNING feeds.*;

-- name: GetNextFeedToFetch :one
SELECT * FROM feeds
ORDER BY last_fetched_at ASC NULLS FIRST;
