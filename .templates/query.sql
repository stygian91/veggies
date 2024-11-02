-- name: GetAllUsers :many
select * from users order by id asc;

-- name: GetUsers :many
select * from users where id in (sqlc.slice(user_ids)) order by id asc;

-- name: GetUser :one
select * from users where id = ? limit 1;

-- name: GetUserItems :many
select * from user_items where user_id in (sqlc.slice(user_ids)) order by user_id asc;

-- name: GetUserCategories :many
select categories_users.*, sqlc.embed(categories) from categories_users
inner join categories on categories.id = categories_users.category_id
where categories_users.user_id in (sqlc.slice(user_ids));
