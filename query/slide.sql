-- name: CreateSlide :one
INSERT INTO "slide" (
    id,
    owner,
    title,
    content
) VALUES (
    $1,
    $2,
    $3,
    $4
) RETURNING *;

-- name: GetSlide :one
SELECT * FROM "slide" WHERE id = $1;

-- name: GetSlidesByOwner :many
SELECT * FROM "slide" WHERE owner = $1;

-- name: UpdateSlide :one
UPDATE "slide" SET
    title = $2,
    content = $3,
    updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteSlide :exec
DELETE FROM "slide" WHERE id = $1;

-- name: CheckSlidePermission :one
SELECT EXISTS (
    SELECT 1
    FROM "slide"
    WHERE id = $1
    AND (
        owner = $2
        OR EXISTS (
            SELECT 1
            FROM "collab"
            WHERE user_id = $2
            AND slide_id = $1
        )
    )
) AS is_permitted;