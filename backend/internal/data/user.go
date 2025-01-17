package data

import (
	"context"
	"database/sql"
	"time"
)

type User struct {
	ID         string    `json:"id"`
	FullName   string    `json:"fullname"`
	Email      string    `json:"email"`
	Password   password  `json:"-"`
	ProfilePic string    `json:"profile_pic"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type UserModel struct {
	DB *sql.DB
}

func (m UserModel) Insert(user *User) error {
	query := `
		INSERT INTO users (fullname, email, password_hash)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at
		`

	args := []any{user.FullName, user.Email, user.Password.hash}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		switch {
		case err.Error() == `ERROR: duplicate key value violates unique constraint "users_email_key" (SQLSTATE 23505)`:
			return ErrDuplicateEmail
		default:
			return err
		}
	}

	return nil
}
