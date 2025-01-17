package dto

import "github.com/ucok-man/fs-chat-app-backend/internal/validator"

type UserReqSignup struct {
	ID       int64  `json:"id"`
	FullName string `json:"fullname"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func validateEmail(v *validator.Validator, email string) {
	v.Check(email != "", "email", "must be provided")
	v.Check(validator.Matches(email, validator.EmailRX), "email", "must be valid email address")
}

func validatePasswordPlaintext(v *validator.Validator, password string) {
	v.Check(password != "", "password", "must be provided")
	v.Check(len(password) >= 8, "password", "must be at least 8 bytes long")
	v.Check(len(password) <= 72, "password", "must not be more than 72 bytes long")
}

func (dto *UserReqSignup) Validate() map[string]string {
	v := validator.New()

	v.Check(dto.FullName != "", "full_name", "must be provided")
	v.Check(len(dto.FullName) <= 500, "full_name", "must not be more than 500 bytes long")
	validateEmail(v, dto.Email)
	validatePasswordPlaintext(v, dto.Password)

	if !v.Valid() {
		return v.Errors
	}
	return nil
}
