-- name: CreateUser :one
INSERT INTO "user" (
  user_id,
  email,
  name,
  password
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;

-- name: GetUser :one
SELECT * FROM "user"
WHERE user_id = $1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM "user"
WHERE email = $1 LIMIT 1;

-- name: ListUser :many
SELECT * FROM "user"
ORDER BY user_id
LIMIT $1
OFFSET $2;

-- name: UpdatePassword :one
UPDATE "user"
SET password = $2
WHERE email = $1
RETURNING *;

-- name: Verify :one
UPDATE "user"
SET verified = true
WHERE email = $1
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM "user"
WHERE email = $1;