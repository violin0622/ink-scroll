package backoff

import (
	"context"
	"time"

	"github.com/cenkalti/backoff/v5"
)

func RetryWith[A any](
	ctx context.Context,
	fn func() (A, error),
	opts ...backoff.RetryOption,
) (A, error) {
	return backoff.Retry[A](ctx, fn, opts...)
}

func novalue(fn func() error) backoff.Operation[struct{}] {
	return func() (struct{}, error) { return struct{}{}, fn() }
}

func Retry(
	ctx context.Context,
	fn func() error,
	opts ...backoff.RetryOption,
) error {
	_, err := backoff.Retry(
		ctx,
		novalue(fn),
		opts...)
	return err
}

const (
	DefaultInitialInterval     = 500 * time.Millisecond
	DefaultRandomizationFactor = 0.5
	DefaultMultiplier          = 1.5
	DefaultMaxInterval         = 60 * time.Second
)
const DefaultMaxElapsedTime = 15 * time.Minute
const Stop time.Duration = -1

func Permanent(err error) error { return backoff.Permanent(err) }
