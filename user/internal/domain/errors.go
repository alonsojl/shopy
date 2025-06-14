package domain

import "shopy/pkg/errorx"

const (
	CodeBadRequest errorx.Code = iota
	CodeNotFound
	CodeUnauthorized
)

var (
	ErrBodyRequest  = errorx.NewErrorf(CodeBadRequest, "invalid body request")
	ErrParams       = errorx.NewErrorf(CodeBadRequest, "invalid params")
	ErrNotFound     = errorx.NewErrorf(CodeNotFound, "item not found")
	ErrUnauthorized = errorx.NewErrorf(CodeUnauthorized, "invalid credentials")
)
