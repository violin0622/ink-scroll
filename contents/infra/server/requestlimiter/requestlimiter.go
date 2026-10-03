// Package requestlimter provides Interceptor for request rate limiter
// and dependent interfaces.
package requestlimter

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/protoadapt"

	"github.com/violin0622/ink-scroll/contents/infra/server/stats"
)

type RequestLimiter interface {
	Request(RequestInfo) Result
}

type RequestLimiterStats interface {
	Stats() Stats
}

type Stats struct {
	Total, Denied uint64
	Limit, Tokens float64
	Burst         int
}

type RequestInfo struct {
	FullMethod string
	Peer       *peer.Peer
	MD         metadata.MD
	Custom     map[string]any
	PS         *stats.PayloadSize
}

type Result struct {
	Deny    bool
	Details []protoadapt.MessageV1
}

func UnaryServerInterceptor(rl RequestLimiter) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (
		resp any, err error,
	) {
		ri := fromContext(ctx)
		ri.FullMethod = info.FullMethod

		if r := rl.Request(ri); r.Deny {
			s, _ := status.New(
				codes.ResourceExhausted,
				`rate limit exceeded`).
				WithDetails(r.Details...)

			return nil, s.Err()
		}

		return handler(ctx, req)
	}
}

func StreamServerInterceptor(rl RequestLimiter) grpc.StreamServerInterceptor {
	return func(
		srv any,
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		ri := fromContext(ss.Context())
		ri.FullMethod = info.FullMethod

		if r := rl.Request(ri); r.Deny {
			s, _ := status.New(
				codes.ResourceExhausted,
				`rate limit exceeded`).
				WithDetails(r.Details...)

			return s.Err()
		}

		return handler(srv, ss)
	}
}

func fromContext(ctx context.Context) (ri RequestInfo) {
	p, ok := peer.FromContext(ctx)
	if ok {
		ri.Peer = p
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		ri.MD = md
	}

	ps, ok := stats.FromContext(ctx)
	if ok {
		ri.PS = ps
	}

	return ri
}
