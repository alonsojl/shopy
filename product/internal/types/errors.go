package types

import "product/pkg/errorx"

const (
	CodeBadRequest errorx.Code = iota
	CodeNotFound
)

var (
	ErrRequest  = errorx.NewErrorf(CodeBadRequest, "invalid body request")
	ErrParams   = errorx.NewErrorf(CodeBadRequest, "invalid params")
	ErrNotFound = errorx.NewErrorf(CodeNotFound, "item not found")
)
