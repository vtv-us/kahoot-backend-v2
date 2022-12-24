-- name: UpsertUserQuestion :one
INSERT INTO "user_question" (
  "question_id",
  "slide_id",
  "username",
  "content",
  "created_at"
) VALUES (
    $1, $2, $3, $4, now()
) ON CONFLICT (question_id) DO UPDATE SET
    "slide_id" = $2,
    "username" = $3,
    "content" = $4
RETURNING *;

-- name: GetUserQuestion :one
SELECT *
FROM "user_question"
WHERE question_id = $1;

-- name: ListUserQuestion :many
SELECT *
FROM "user_question"
WHERE slide_id = $1
ORDER BY created_at DESC;

-- name: UpvoteUserQuestion :one
UPDATE "user_question"
SET votes = votes + 1
WHERE question_id = $1
RETURNING *;

-- name: ToggleUserQuestionAnswered :one
UPDATE "user_question"
SET answered = NOT answered
WHERE question_id = $1
RETURNING *;
