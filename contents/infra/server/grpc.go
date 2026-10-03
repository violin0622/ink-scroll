package server

import (
	"buf.build/go/protovalidate"
	"github.com/go-logr/logr"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	pvm "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/protovalidate"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthv1 "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	"github.com/violin0622/ink-scroll/contents/infra/server/log"
	"github.com/violin0622/ink-scroll/contents/infra/server/requestlimiter"
	"github.com/violin0622/ink-scroll/contents/infra/server/stats"
	contentsv1 "github.com/violin0622/ink-scroll/proto/gen/go/contents/v1"
)

func New(
	contCmd contentsv1.CmdServiceServer,
	contQuery contentsv1.QueryServiceServer,
	opts ...Option,
) *grpc.Server {
	var b builder
	b.logr = logr.Discard()
	for _, o := range opts {
		o(&b)
	}

	l := log.New(b.logr)
	b.gRPCOpts = append(
		[]grpc.ServerOption{
			grpc.StatsHandler(otelgrpc.NewServerHandler()),
			grpc.StatsHandler(stats.NewInjectPayloadSize()),
			grpc.ChainUnaryInterceptor(
				requestlimter.UnaryServerInterceptor(b.rl),
				logging.UnaryServerInterceptor(l),
				recovery.UnaryServerInterceptor(),
				pvm.UnaryServerInterceptor(protovalidate.GlobalValidator),
			),
			grpc.ChainStreamInterceptor(
				requestlimter.StreamServerInterceptor(b.rl),
				logging.StreamServerInterceptor(l),
				recovery.StreamServerInterceptor(),
				pvm.StreamServerInterceptor(protovalidate.GlobalValidator),
			),
		},
		b.gRPCOpts...)

	s := grpc.NewServer(b.gRPCOpts...)

	if b.hs != nil {
		healthv1.RegisterHealthServer(s, b.hs)
	}
	reflection.RegisterV1(s)
	contentsv1.RegisterCmdServiceServer(s, contCmd)
	contentsv1.RegisterQueryServiceServer(s, contQuery)

	return s
}

type builder struct {
	gRPCOpts []grpc.ServerOption
	hs       *health.Server
	logr     logr.Logger
	rl       requestlimter.RequestLimiter
}

type Option func(*builder)

func WithRequestLimiter(rl requestlimter.RequestLimiter) Option {
	return func(b *builder) { b.rl = rl }
}

func WithLogr(l logr.Logger) Option {
	return func(b *builder) { b.logr = l }
}

func WithHealth(hs *health.Server) Option {
	return func(b *builder) { b.hs = hs }
}

func WithGRPCOptions(opts ...grpc.ServerOption) Option {
	return func(b *builder) { b.gRPCOpts = append(b.gRPCOpts, opts...) }
}
