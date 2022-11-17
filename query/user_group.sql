-- name: ListGroupJoined :many
SELECT g.group_id, group_name, ug.role, created_by, g.created_at
FROM "user_group" ug
INNER JOIN "group" g using (group_id)
WHERE user_id = $1
ORDER BY g.group_id;

-- name: ListMemberInGroup :one
SELECT user_id, role
FROM "user_group"
WHERE group_id = $1
ORDER BY user_id;

-- name: AddMemberToGroup :exec
INSERT INTO "user_group" (
  user_id,
  group_id,
  role
) VALUES (
  $1, $2, $3
);

-- name: RemoveMemberFromGroup :exec
DELETE FROM "user_group"
WHERE user_id = $1 AND group_id = $2;

-- name: UpdateMemberRole :exec
UPDATE "user_group"
SET role = $3
WHERE user_id = $1 AND group_id = $2;

-- name: GetRoleInGroup :one
SELECT role
FROM "user_group"
WHERE user_id = $1 AND group_id = $2;