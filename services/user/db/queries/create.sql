-- name: CreateUser :one
INSERT INTO users (
  auth_id,
  avatar_url,
  bio,
  theme,
  language,
  email_notifications,
  last_login_at
) VALUES (
  $1,
  $2,
  $3,
  $4,
  $5,
  $6,
  $7
)
RETURNING *;
