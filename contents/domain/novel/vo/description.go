package vo

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Description string

func NewDescription(d string) (Description, error) {
	if len(d) < 20 || len(d) > 500 {
		return ``, status.Error(codes.InvalidArgument, `invalid description length, must be 20-500 words`)
	}

	return Description(d), nil
}
