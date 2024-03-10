-- name: GetUserName :one
SELECT username FROM users WHERE id = $1;

-- name: GetUser :one
SELECT * FROM users WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: CreateUser :one
INSERT INTO users (email, password, username) VALUES ($1, $2, $3) RETURNING *;

-- name: UpdateUser :one
UPDATE users SET email = $1, password = $2, username = $3 WHERE id = $4 RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;

