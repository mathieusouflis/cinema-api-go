-- name: GetUserById :one
SELECT * FROM users WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: GetUserByUsername :one
SELECT * FROM users WHERE username = $1;

-- name: CreateUser :one
INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;

-- name: GetUserByGoogleID :one
SELECT * FROM users WHERE google_id = $1;

-- name: UpdateGoogleID :exec
UPDATE users SET google_id = $2 WHERE id = $1;

-- name: CreateUserWithOAuth :one
INSERT INTO users (username, email, password, google_id) VALUES ($1, $2, '', $3) RETURNING *;
