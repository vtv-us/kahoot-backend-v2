-- name: CreateQuestion :one
INSERT INTO "question" (
    id,
    slide_id,
    raw_question,
    answer_a,
    answer_b,
    answer_c,
    answer_d,
    correct_answer
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8
) RETURNING *;

-- name: GetQuestion :one
SELECT * FROM "question" WHERE id = $1;

-- name: GetQuestionsBySlide :many
SELECT * FROM "question" WHERE slide_id = $1;

-- name: GetOwnerOfQuestion :one
SELECT s.owner FROM "question" q
JOIN "slide" s ON q.slide_id = s.id
WHERE q.id = $1;

-- name: UpdateQuestion :one
UPDATE "question" SET
    raw_question = $2,
    answer_a = $3,
    answer_b = $4,
    answer_c = $5,
    answer_d = $6,
    correct_answer = $7,
    updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteQuestion :exec
DELETE FROM "question" WHERE id = $1;