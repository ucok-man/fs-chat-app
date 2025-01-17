package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
	ErrDuplicateEmail = errors.New("duplicate email")
)

type Models struct {
	User UserModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		User: UserModel{DB: db},
	}
}
