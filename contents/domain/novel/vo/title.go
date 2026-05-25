package vo

import (
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrTitleTooLong = status.Error(
		codes.InvalidArgument,
		`title exceed max length 256`,
	)
	ErrTitleInvalidChar = status.Error(
		codes.InvalidArgument,
		`title contains invalid charactor @#$%^&*()_+-=\/`,
	)
	ErrTitleEmpty = status.Error(codes.InvalidArgument, `empty title`)
)

type Title string

func NewTitle(s string) (Title, error) {
	if len(s) == 0 {
		return ``, ErrTitleEmpty
	}

	if len(s) >= 256 {
		return ``, ErrTitleTooLong
	}

	if strings.Contains(s, `@#$%^&*()_+-=\/`) {
		return ``, ErrTitleInvalidChar
	}

	return Title(s), nil
}
