package apigateway

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type UserAddRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (u UserAddRequest) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.Email,
			validation.Required,
			is.EmailFormat,
		),
		validation.Field(&u.Password,
			validation.Required,
			is.Alphanumeric,
		),
	)
}

type UserCredentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (u UserCredentials) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.Email,
			validation.Required,
			is.EmailFormat,
		),
		validation.Field(&u.Password,
			validation.Required,
			is.Alphanumeric,
		),
	)
}
