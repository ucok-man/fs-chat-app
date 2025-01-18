package dto

import "github.com/ucok-man/fs-chat-app-backend/internal/validator"

type ReqSignupDto struct {
	FullName string `json:"fullname"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (dto *ReqSignupDto) Validate() map[string]string {
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
