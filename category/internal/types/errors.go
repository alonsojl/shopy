package types

import "category/pkg/errorx"

const (
	CodeInvalidArgument errorx.ErrorCode = iota
	CodePrecondition
	CodeNotFound
)

var (
	ErrBodyRequest = errorx.NewErrorf(CodeInvalidArgument, "invalid body request")
	ErrParams      = errorx.NewErrorf(CodeInvalidArgument, "invalid params")
	ErrNotFound    = errorx.NewErrorf(CodeNotFound, "item not found")
)
