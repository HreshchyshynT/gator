-- name: CreatePost :one 
insert into posts (
  id,
  created_at,
  updated_at,
  title,
  url,
  description,
  published_at,
  feed_id
) 
values (
  $1,
  $2,
  $3,
  $4,
  $5,
  $6,
  $7,
  $8
) RETURNING *;

-- name: GetPostsForUserNewest :many
select posts.* from posts
inner join feed_follows 
  on feed_follows.feed_id = posts.feed_id and feed_follows.user_id = @user_id
inner join feeds 
  on feeds.id = posts.feed_id
where 
  (sqlc.narg('feed_name')::text is NULL or feeds.name = sqlc.narg('feed_name'))
  AND 
  (sqlc.narg('since')::timestamp is NULL or posts.published_at > sqlc.narg('since'))
order by posts.published_at DESC, posts.id ASC  
limit @lim;

-- name: GetPostsForUserOldest :many
select posts.* from posts
inner join feed_follows 
  on feed_follows.feed_id = posts.feed_id and feed_follows.user_id = @user_id
inner join feeds 
  on feeds.id = posts.feed_id
where 
  (sqlc.narg('feed_name')::text is NULL or feeds.name = sqlc.narg('feed_name'))
  AND 
  (sqlc.narg('since')::timestamp is NULL or posts.published_at > sqlc.narg('since'))
order by posts.published_at ASC, posts.id ASC
limit @lim;

-- name: GetPostsForUserTitle :many
select posts.* from posts
inner join feed_follows 
  on feed_follows.feed_id = posts.feed_id and feed_follows.user_id = @user_id
inner join feeds 
  on feeds.id = posts.feed_id
where 
  (sqlc.narg('feed_name')::text is NULL or feeds.name = sqlc.narg('feed_name'))
  AND 
  (sqlc.narg('since')::timestamp is NULL or posts.published_at > sqlc.narg('since'))
order by title ASC, posts.id ASC  
limit @lim;

-- name: GetPostsForUserFeed :many
select posts.* from posts
inner join feed_follows 
  on feed_follows.feed_id = posts.feed_id and feed_follows.user_id = @user_id
inner join feeds
  on feeds.id = posts.feed_id
where 
  (sqlc.narg('feed_name')::text is NULL or feeds.name = sqlc.narg('feed_name'))
  AND 
  (sqlc.narg('since')::timestamp is NULL or posts.published_at > sqlc.narg('since'))
order by feeds.name ASC, posts.published_at DESC, posts.id ASC 
limit @lim;

