package data

import (
	"database/sql"
	"time"
)

type User struct {
	ID         int64     `json:"id"`
	FullName   string    `json:"full_name"`
	Email      string    `json:"email"`
	Password   password  `json:"-"`
	ProfilePic string    `json:"profile_pic"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type UserModel struct {
	DB *sql.DB
}
