-- name: CreateAnswer :one
INSERT INTO "answer" (
    id,
    question_id,
    index,
    raw_answer,
    created_at,
    updated_at
) VALUES (
    $1,
    $2,
    $3,
    $4,
    now(),
    now()
)
RETURNING *;

-- name: GetAnswer :one
SELECT * FROM "answer" WHERE id = $1;

-- name: GetAnswersByQuestion :many
SELECT * FROM "answer" WHERE question_id = $1
ORDER BY index ASC;

-- name: UpdateAnswer :one
UPDATE "answer" SET
    index = $2,
    raw_answer = $3,
    updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteAnswer :exec
DELETE FROM "answer" WHERE id = $1;

-- name: DeleteAnswersByQuestion :exec
DELETE FROM "answer" WHERE question_id = $1;

-- name: DeleteAnswersBySlide :exec
DELETE FROM "answer" WHERE question_id IN (
    SELECT id FROM "question" WHERE slide_id = $1
);