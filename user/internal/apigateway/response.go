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

type UserAuthorized struct {
	BaseResponse
	Token string `json:"token"`
}

type UserAdded struct {
	BaseResponse
	User *models.User `json:"user"`
}

type UserDeleted struct {
	BaseResponse
	User string `json:"user"`
}
