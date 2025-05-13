package apigateway

import (
	"category/internal/models"
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

type SelectedCategories struct {
	BaseResponse
	Categories models.Categories `json:"categories"`
}

type CategoryAdded struct {
	BaseResponse
	Category *models.Category `json:"category"`
}

type CategoryDeleted struct {
	BaseResponse
	Category string `json:"category"`
}
