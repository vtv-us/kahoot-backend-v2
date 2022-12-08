-- name: CreateQuestion :one
INSERT INTO "question" (
    id,
    slide_id,
    index,
    raw_question,
    meta,
    long_description,
    created_at,
    updated_at
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    now(),
    now()
)
RETURNING *;

-- name: GetQuestion :one
SELECT * FROM "question" WHERE id = $1;

-- name: GetQuestionsBySlide :many
SELECT * FROM "question" WHERE slide_id = $1
ORDER BY index ASC;

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
    updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteQuestion :exec
DELETE FROM "question" WHERE id = $1;

-- name: DeleteQuestionsBySlide :exec
DELETE FROM "question" WHERE slide_id = $1;