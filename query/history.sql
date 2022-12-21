-- name: UpsertAnswerHistory :one
INSERT INTO "answer_history" (
    "username",
    "slide_id",
    "question_id",
    "answer_id"
) VALUES (
    $1, $2, $3, $4
) ON CONFLICT ON CONSTRAINT "answer_history_pkey" DO UPDATE SET
    "answer_id" = $4,
    "updated_at" = now()
RETURNING *;

-- name: GetAnswerHistory :one
SELECT *
FROM "answer_history"
WHERE username = $1 AND slide_id = $2 AND question_id = $3;

-- name: ListAnswerHistoryBySlideID :many
SELECT *
FROM "answer_history"
WHERE slide_id = $1
ORDER BY updated_at DESC;

-- name: ListAnswerHistoryByQuestionID :many
SELECT *
FROM "answer_history"
WHERE question_id = $1
ORDER BY updated_at DESC;

-- name: ListAnswerHistoryByAnswerID :many
SELECT *
FROM "answer_history"
WHERE answer_id = $1
ORDER BY updated_at DESC;

-- name: CountAnswerByQuestionID :many
SELECT slide_id, question_id, answer_id, count(*) as count
FROM "answer_history"
WHERE question_id = $1
GROUP BY slide_id, question_id, answer_id
ORDER BY count DESC;
