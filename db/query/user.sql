
-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;

-- name: CreateUser :one
INSERT INTO users (
  first_name,
  last_name,
  email,
  encrypted_password,
  company_id
) VALUES (
  $1, $2, $3, $4, $5
)
RETURNING *;
