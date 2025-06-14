package apigateway

import (
	"shopy/internal/models"
	"time"
)

type BaseResponse struct {
	Status   string `json:"status"`
	Code     int    `json:"code"`
	Datetime string `json:"datetime"`
}

func NewBaseResponse(code int) BaseResponse {
	return BaseResponse{
		Status:   "success",
		Code:     code,
		Datetime: time.Now().Format(time.DateTime),
	}
}

type SelectedProducts struct {
	BaseResponse
	Products models.Products `json:"products"`
}

type ProductAdded struct {
	BaseResponse
	Product *models.Product `json:"product"`
}

type ProductDeleted struct {
	BaseResponse
	Product string `json:"product"`
}
