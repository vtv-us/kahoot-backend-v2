-- name: SaveHistory :one
INSERT INTO "answer_history"
    ("id", "slide_id", "raw_question", "raw_answer", "num_chosen")
VALUES
    ($1, $2, $3, $4, $5)
RETURNING *;

