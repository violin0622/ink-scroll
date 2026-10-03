package grpc

import (
	"github.com/violin0622/ink-scroll/contents/app"
	contv1 "github.com/violin0622/ink-scroll/proto/gen/go/contents/v1"
)

type Query struct {
	contv1.UnimplementedQueryServiceServer
}

func NewQueryService(
	ah app.QueryHandler,
) Query {
	return Query{}
}
