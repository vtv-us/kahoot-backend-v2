-- name: AddCollab :exec
INSERT INTO "collab" (
  user_id,
  slide_id
) VALUES (
  $1, $2
);

-- name: RemoveCollab :exec
DELETE FROM "collab"
WHERE user_id = $1
AND slide_id = $2;

-- name: ListCollab :many
SELECT g.*
FROM "collab" c
JOIN "slide" g on g.id = c.slide_id
WHERE c.user_id = $1;

-- name: ListCollabBySlide :many
SELECT u.*
FROM "collab" c
JOIN "user" u using (user_id)
WHERE c.slide_id = $1;

-- name: CheckIsCollab :one
SELECT EXISTS (
  SELECT 1
  FROM "collab"
  WHERE user_id = $1
  AND slide_id = $2
) AS is_collab;