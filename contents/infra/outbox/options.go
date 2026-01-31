package outbox

import (
	"time"

	"github.com/go-logr/logr"
)

// Option is a function that configures a Relay.
type Option func(*Relay)

// WithMaxFetch sets the maximum number of events to fetch per tick.
func WithMaxFetch(n int) Option {
	return func(r *Relay) {
		if n > 0 {
			r.maxFetch = n
		}
	}
}

// WithMaxPub sets the maximum number of events to publish per batch.
func WithMaxPub(n int) Option {
	return func(r *Relay) {
		if n > 0 {
			r.maxPub = n
		}
	}
}

// WithTickTimeout sets the timeout for each tick operation.
func WithTickTimeout(d time.Duration) Option {
	return func(r *Relay) {
		if d > 0 {
			r.tickTimeout = d
		}
	}
}

// WithInterval sets the polling interval.
func WithInterval(d time.Duration) Option {
	return func(r *Relay) {
		if d > 0 {
			r.interval = d
		}
	}
}

func WithNotify(n <-chan struct{}) Option {
	return func(r *Relay) {
		r.notify = n
	}
}

func WithLogr(l logr.Logger) Option {
	return func(r *Relay) { r.log = l }
}
