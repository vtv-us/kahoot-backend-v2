-- name: SaveChat :one
INSERT INTO "chat_msg" (
    id,
    slide_id,
    username,
    content,
    created_at
) VALUES (
    $1,
    $2,
    $3,
    $4,
    now()
)
RETURNING *;

-- name: GetChatBySlide :many
SELECT * FROM "chat_msg" WHERE slide_id = $1
ORDER BY created_at ASC;