package dto

import "github.com/ucok-man/fs-chat-app-backend/internal/validator"

type ReqSigninDto struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (dto *ReqSigninDto) Validate() map[string]string {
	v := validator.New()

	validateEmail(v, dto.Email)
	validatePasswordPlaintext(v, dto.Password)

	if !v.Valid() {
		return v.Errors
	}
	return nil
}
