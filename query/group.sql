-- name: CreateGroup :one
INSERT INTO "group" (
  group_id,
  group_name,
  created_by,
  description
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;

-- name: GetGroup :one
SELECT *
FROM "group"
WHERE group_id = $1;

-- name: ListGroupOwned :many
SELECT *
FROM "group"
JOIN "user_group" ug using (group_id)
WHERE ug.user_id = $1
AND ug.role = 'owner'
AND ug.status = 'joined'
ORDER BY ug.group_id;

-- name: DeleteGroup :exec
DELETE FROM "group"
WHERE group_id = $1;
