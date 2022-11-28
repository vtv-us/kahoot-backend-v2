// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.14.0
// source: user.sql

package repositories

import (
	"context"
	"database/sql"
)

const createUser = `-- name: CreateUser :one
INSERT INTO "user" (
  user_id,
  email,
  name,
  password,
  verified,
  verified_code,
  google_id,
  facebook_id
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8
)
RETURNING user_id, email, name, password, verified, verified_code, created_at, google_id, facebook_id, avatar_url
`

type CreateUserParams struct {
	UserID       string         `json:"user_id"`
	Email        string         `json:"email"`
	Name         string         `json:"name"`
	Password     string         `json:"password"`
	Verified     bool           `json:"verified"`
	VerifiedCode string         `json:"verified_code"`
	GoogleID     sql.NullString `json:"google_id"`
	FacebookID   sql.NullString `json:"facebook_id"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser,
		arg.UserID,
		arg.Email,
		arg.Name,
		arg.Password,
		arg.Verified,
		arg.VerifiedCode,
		arg.GoogleID,
		arg.FacebookID,
	)
	var i User
	err := row.Scan(
		&i.UserID,
		&i.Email,
		&i.Name,
		&i.Password,
		&i.Verified,
		&i.VerifiedCode,
		&i.CreatedAt,
		&i.GoogleID,
		&i.FacebookID,
		&i.AvatarUrl,
	)
	return i, err
}

const deleteUser = `-- name: DeleteUser :exec
DELETE FROM "user"
WHERE email = $1
`

func (q *Queries) DeleteUser(ctx context.Context, email string) error {
	_, err := q.db.ExecContext(ctx, deleteUser, email)
	return err
}

const getUser = `-- name: GetUser :one
SELECT user_id, email, name, password, verified, verified_code, created_at, google_id, facebook_id, avatar_url FROM "user"
WHERE user_id = $1 LIMIT 1
`

func (q *Queries) GetUser(ctx context.Context, userID string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUser, userID)
	var i User
	err := row.Scan(
		&i.UserID,
		&i.Email,
		&i.Name,
		&i.Password,
		&i.Verified,
		&i.VerifiedCode,
		&i.CreatedAt,
		&i.GoogleID,
		&i.FacebookID,
		&i.AvatarUrl,
	)
	return i, err
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT user_id, email, name, password, verified, verified_code, created_at, google_id, facebook_id, avatar_url FROM "user"
WHERE email = $1 LIMIT 1
`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByEmail, email)
	var i User
	err := row.Scan(
		&i.UserID,
		&i.Email,
		&i.Name,
		&i.Password,
		&i.Verified,
		&i.VerifiedCode,
		&i.CreatedAt,
		&i.GoogleID,
		&i.FacebookID,
		&i.AvatarUrl,
	)
	return i, err
}

const listUser = `-- name: ListUser :many
SELECT user_id, email, name, password, verified, verified_code, created_at, google_id, facebook_id, avatar_url FROM "user"
ORDER BY user_id
LIMIT $1
OFFSET $2
`

type ListUserParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListUser(ctx context.Context, arg ListUserParams) ([]User, error) {
	rows, err := q.db.QueryContext(ctx, listUser, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []User{}
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.UserID,
			&i.Email,
			&i.Name,
			&i.Password,
			&i.Verified,
			&i.VerifiedCode,
			&i.CreatedAt,
			&i.GoogleID,
			&i.FacebookID,
			&i.AvatarUrl,
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

const updateAvatarUrl = `-- name: UpdateAvatarUrl :one
UPDATE "user"
SET avatar_url = $2
WHERE user_id = $1
RETURNING user_id, email, name, password, verified, verified_code, created_at, google_id, facebook_id, avatar_url
`

type UpdateAvatarUrlParams struct {
	UserID    string         `json:"user_id"`
	AvatarUrl sql.NullString `json:"avatar_url"`
}

func (q *Queries) UpdateAvatarUrl(ctx context.Context, arg UpdateAvatarUrlParams) (User, error) {
	row := q.db.QueryRowContext(ctx, updateAvatarUrl, arg.UserID, arg.AvatarUrl)
	var i User
	err := row.Scan(
		&i.UserID,
		&i.Email,
		&i.Name,
		&i.Password,
		&i.Verified,
		&i.VerifiedCode,
		&i.CreatedAt,
		&i.GoogleID,
		&i.FacebookID,
		&i.AvatarUrl,
	)
	return i, err
}

const updatePassword = `-- name: UpdatePassword :one
UPDATE "user"
SET password = $2
WHERE email = $1
RETURNING user_id, email, name, password, verified, verified_code, created_at, google_id, facebook_id, avatar_url
`

type UpdatePasswordParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (q *Queries) UpdatePassword(ctx context.Context, arg UpdatePasswordParams) (User, error) {
	row := q.db.QueryRowContext(ctx, updatePassword, arg.Email, arg.Password)
	var i User
	err := row.Scan(
		&i.UserID,
		&i.Email,
		&i.Name,
		&i.Password,
		&i.Verified,
		&i.VerifiedCode,
		&i.CreatedAt,
		&i.GoogleID,
		&i.FacebookID,
		&i.AvatarUrl,
	)
	return i, err
}

const updateProfile = `-- name: UpdateProfile :one
UPDATE "user"
SET name = $2
WHERE user_id = $1
RETURNING user_id, email, name, password, verified, verified_code, created_at, google_id, facebook_id, avatar_url
`

type UpdateProfileParams struct {
	UserID string `json:"user_id"`
	Name   string `json:"name"`
}

func (q *Queries) UpdateProfile(ctx context.Context, arg UpdateProfileParams) (User, error) {
	row := q.db.QueryRowContext(ctx, updateProfile, arg.UserID, arg.Name)
	var i User
	err := row.Scan(
		&i.UserID,
		&i.Email,
		&i.Name,
		&i.Password,
		&i.Verified,
		&i.VerifiedCode,
		&i.CreatedAt,
		&i.GoogleID,
		&i.FacebookID,
		&i.AvatarUrl,
	)
	return i, err
}

const updateSocialID = `-- name: UpdateSocialID :one
UPDATE "user"
SET google_id = $2, facebook_id = $3
WHERE email = $1
RETURNING user_id, email, name, password, verified, verified_code, created_at, google_id, facebook_id, avatar_url
`

type UpdateSocialIDParams struct {
	Email      string         `json:"email"`
	GoogleID   sql.NullString `json:"google_id"`
	FacebookID sql.NullString `json:"facebook_id"`
}

func (q *Queries) UpdateSocialID(ctx context.Context, arg UpdateSocialIDParams) (User, error) {
	row := q.db.QueryRowContext(ctx, updateSocialID, arg.Email, arg.GoogleID, arg.FacebookID)
	var i User
	err := row.Scan(
		&i.UserID,
		&i.Email,
		&i.Name,
		&i.Password,
		&i.Verified,
		&i.VerifiedCode,
		&i.CreatedAt,
		&i.GoogleID,
		&i.FacebookID,
		&i.AvatarUrl,
	)
	return i, err
}

const verify = `-- name: Verify :one
UPDATE "user"
SET verified = true
WHERE email = $1
RETURNING user_id, email, name, password, verified, verified_code, created_at, google_id, facebook_id, avatar_url
`

func (q *Queries) Verify(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRowContext(ctx, verify, email)
	var i User
	err := row.Scan(
		&i.UserID,
		&i.Email,
		&i.Name,
		&i.Password,
		&i.Verified,
		&i.VerifiedCode,
		&i.CreatedAt,
		&i.GoogleID,
		&i.FacebookID,
		&i.AvatarUrl,
	)
	return i, err
}
