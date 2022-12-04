// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.14.0
// source: question.sql

package repositories

import (
	"context"
)

const createQuestion = `-- name: CreateQuestion :one
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
) RETURNING id, slide_id, raw_question, answer_a, answer_b, answer_c, answer_d, correct_answer, created_at, updated_at
`

type CreateQuestionParams struct {
	ID            string `json:"id"`
	SlideID       string `json:"slide_id"`
	RawQuestion   string `json:"raw_question"`
	AnswerA       string `json:"answer_a"`
	AnswerB       string `json:"answer_b"`
	AnswerC       string `json:"answer_c"`
	AnswerD       string `json:"answer_d"`
	CorrectAnswer string `json:"correct_answer"`
}

func (q *Queries) CreateQuestion(ctx context.Context, arg CreateQuestionParams) (Question, error) {
	row := q.db.QueryRowContext(ctx, createQuestion,
		arg.ID,
		arg.SlideID,
		arg.RawQuestion,
		arg.AnswerA,
		arg.AnswerB,
		arg.AnswerC,
		arg.AnswerD,
		arg.CorrectAnswer,
	)
	var i Question
	err := row.Scan(
		&i.ID,
		&i.SlideID,
		&i.RawQuestion,
		&i.AnswerA,
		&i.AnswerB,
		&i.AnswerC,
		&i.AnswerD,
		&i.CorrectAnswer,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteQuestion = `-- name: DeleteQuestion :exec
DELETE FROM "question" WHERE id = $1
`

func (q *Queries) DeleteQuestion(ctx context.Context, id string) error {
	_, err := q.db.ExecContext(ctx, deleteQuestion, id)
	return err
}

const getOwnerOfQuestion = `-- name: GetOwnerOfQuestion :one
SELECT s.owner FROM "question" q
JOIN "slide" s ON q.slide_id = s.id
WHERE q.id = $1
`

func (q *Queries) GetOwnerOfQuestion(ctx context.Context, id string) (string, error) {
	row := q.db.QueryRowContext(ctx, getOwnerOfQuestion, id)
	var owner string
	err := row.Scan(&owner)
	return owner, err
}

const getQuestion = `-- name: GetQuestion :one
SELECT id, slide_id, raw_question, answer_a, answer_b, answer_c, answer_d, correct_answer, created_at, updated_at FROM "question" WHERE id = $1
`

func (q *Queries) GetQuestion(ctx context.Context, id string) (Question, error) {
	row := q.db.QueryRowContext(ctx, getQuestion, id)
	var i Question
	err := row.Scan(
		&i.ID,
		&i.SlideID,
		&i.RawQuestion,
		&i.AnswerA,
		&i.AnswerB,
		&i.AnswerC,
		&i.AnswerD,
		&i.CorrectAnswer,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getQuestionsBySlide = `-- name: GetQuestionsBySlide :many
SELECT id, slide_id, raw_question, answer_a, answer_b, answer_c, answer_d, correct_answer, created_at, updated_at FROM "question" WHERE slide_id = $1
`

func (q *Queries) GetQuestionsBySlide(ctx context.Context, slideID string) ([]Question, error) {
	rows, err := q.db.QueryContext(ctx, getQuestionsBySlide, slideID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Question{}
	for rows.Next() {
		var i Question
		if err := rows.Scan(
			&i.ID,
			&i.SlideID,
			&i.RawQuestion,
			&i.AnswerA,
			&i.AnswerB,
			&i.AnswerC,
			&i.AnswerD,
			&i.CorrectAnswer,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateQuestion = `-- name: UpdateQuestion :one
UPDATE "question" SET
    raw_question = $2,
    answer_a = $3,
    answer_b = $4,
    answer_c = $5,
    answer_d = $6,
    correct_answer = $7,
    updated_at = now()
WHERE id = $1
RETURNING id, slide_id, raw_question, answer_a, answer_b, answer_c, answer_d, correct_answer, created_at, updated_at
`

type UpdateQuestionParams struct {
	ID            string `json:"id"`
	RawQuestion   string `json:"raw_question"`
	AnswerA       string `json:"answer_a"`
	AnswerB       string `json:"answer_b"`
	AnswerC       string `json:"answer_c"`
	AnswerD       string `json:"answer_d"`
	CorrectAnswer string `json:"correct_answer"`
}

func (q *Queries) UpdateQuestion(ctx context.Context, arg UpdateQuestionParams) (Question, error) {
	row := q.db.QueryRowContext(ctx, updateQuestion,
		arg.ID,
		arg.RawQuestion,
		arg.AnswerA,
		arg.AnswerB,
		arg.AnswerC,
		arg.AnswerD,
		arg.CorrectAnswer,
	)
	var i Question
	err := row.Scan(
		&i.ID,
		&i.SlideID,
		&i.RawQuestion,
		&i.AnswerA,
		&i.AnswerB,
		&i.AnswerC,
		&i.AnswerD,
		&i.CorrectAnswer,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
