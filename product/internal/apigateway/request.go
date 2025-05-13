package apigateway

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type ProductAddRequest struct {
	Name     string          `json:"name"`
	Price    float64         `json:"price"`
	Image    string          `json:"image"`
	QRCode   string          `json:"qrcode"`
	IsTop    bool            `json:"is_top"`
	Category CategoryRequest `json:"category"`
}

func (p ProductAddRequest) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.Name,
			validation.Required,
			validation.Length(1, 100),
		),
		validation.Field(&p.Price,
			validation.Required,
			validation.Min(0.0),
		),
		validation.Field(&p.Image,
			validation.Required,
			is.Base64,
		),
		validation.Field(&p.QRCode,
			validation.When(p.QRCode != "",
				is.Alphanumeric,
			)),
		validation.Field(&p.Category),
	)
}

type ProductPutRequest struct {
	Name     string          `json:"name"`
	Price    float64         `json:"price"`
	Image    string          `json:"image"`
	QRCode   string          `json:"qrcode"`
	IsTop    bool            `json:"is_top"`
	Category CategoryRequest `json:"category"`
}

func (p ProductPutRequest) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.Name,
			validation.Required,
			validation.Length(1, 100),
		),
		validation.Field(&p.Price,
			validation.Required,
			validation.Min(0.0),
		),
		validation.Field(&p.Image,
			validation.When(p.Image != "",
				is.Base64,
			)),
		validation.Field(&p.QRCode,
			validation.When(p.QRCode != "",
				is.Alphanumeric,
			)),
		validation.Field(&p.Category),
	)
}

type CategoryRequest struct {
	Uuid string `json:"uuid"`
	Name string `json:"name"`
}

func (c CategoryRequest) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Uuid,
			validation.Required,
			is.UUID,
		),
		validation.Field(&c.Name,
			validation.Required,
			validation.Length(1, 50),
		),
	)
}
