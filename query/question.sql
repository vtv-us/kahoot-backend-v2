-- name: CreateQuestion :one
INSERT INTO "question" (
    id,
    slide_id,
    index,
    raw_question,
    meta,
    long_description,
    type,
    created_at,
    updated_at
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    now(),
    now()
)
RETURNING *;

-- name: GetQuestion :one
SELECT * FROM "question" WHERE id = $1;

-- name: GetQuestionsBySlide :many
SELECT * FROM "question" WHERE slide_id = $1
ORDER BY index ASC;

-- name: GetQuestionBySlideAndIndex :one
SELECT * FROM "question" WHERE slide_id = $1 AND index = $2;

-- name: GetOwnerOfQuestion :one
SELECT s.owner FROM "question" q
JOIN "slide" s ON q.slide_id = s.id
WHERE q.id = $1;

-- name: UpdateQuestion :one
UPDATE "question" SET
    raw_question = $2,
    meta = $3,
    long_description = $4,
    index = $5,
    type = $6,
    updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteQuestion :exec
DELETE FROM "question" WHERE id = $1;

-- name: DeleteQuestionsBySlide :exec
DELETE FROM "question" WHERE slide_id = $1;

-- name: CheckQuestionPermission :one
SELECT EXISTS (
    SELECT 1
    FROM "question" q
    JOIN "slide" s ON q.slide_id = s.id
    WHERE q.id = $1
    AND (
        s.owner = $2
        OR EXISTS (
            SELECT 1
            FROM "collab"
            WHERE user_id = $2
            AND slide_id = s.id
        )
    )
) AS is_permitted;