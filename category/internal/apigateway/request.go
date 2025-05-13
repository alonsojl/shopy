package apigateway

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type CategoryAddRequest struct {
	Name  string `json:"name"`
	Image string `json:"image"`
}

func (c CategoryAddRequest) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Name,
			validation.Required,
			validation.Length(1, 50),
		),
		validation.Field(&c.Image,
			validation.Required,
			is.Base64,
		),
	)
}
