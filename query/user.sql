-- name: CreateUser :one
INSERT INTO "user" (
  user_id,
  email,
  name,
  password,
  verified,
  verified_code,
  google_id,
  facebook_id
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8
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

-- name: UpdatePasswordByEmail :one
UPDATE "user"
SET password = $2
WHERE email = $1
RETURNING *;

-- name: UpdatePassword :one
UPDATE "user"
SET password = $2
WHERE user_id = $1
RETURNING *;

-- name: Verify :one
UPDATE "user"
SET verified = true
WHERE email = $1
RETURNING *;

-- name: UpdateSocialID :one
UPDATE "user"
SET google_id = $2, facebook_id = $3
WHERE email = $1
RETURNING *;

-- name: UpdateAvatarUrl :one
UPDATE "user"
SET avatar_url = $2
WHERE user_id = $1
RETURNING *;

-- name: UpdateProfile :one
UPDATE "user"
SET name = $2
WHERE user_id = $1
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM "user"
WHERE email = $1;