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

-- name: GetPostsForUser :many
select posts.* from posts
inner join feed_follows 
on feed_follows.feed_id = posts.feed_id and feed_follows.user_id = @user_id
order by published_at DESC 
limit @lim;
