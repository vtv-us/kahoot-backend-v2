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

-- name: GetAnswerByQuestionAndIndex :one
SELECT * FROM "answer" WHERE question_id = $1 AND index = $2;

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

-- name: CheckAnswerPermission :one
-- Check if the user has permission to access the answer or collaborator
-- of the slide that the answer belongs to.
SELECT EXISTS (
    SELECT 1
    FROM "answer" a
    JOIN "question" q ON a.question_id = q.id
    JOIN "slide" s ON q.slide_id = s.id
    WHERE a.id = $1
    AND (
        s.owner = $2
        OR EXISTS (
            SELECT 1
            FROM "collab"
            WHERE user_id = $2
            AND slide_id = s.id
        )
    )
) AS has_permission;