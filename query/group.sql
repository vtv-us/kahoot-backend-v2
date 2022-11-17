-- name: CreateGroup :one
INSERT INTO "group" (
  group_id,
  group_name,
  created_by
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: ListGroupCreatedByUser :many
SELECT * FROM "group"
WHERE created_by = $1
ORDER BY group_id;

-- name: DeleteGroup :exec
DELETE FROM "group"
WHERE group_id = $1;
