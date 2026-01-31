// Package outbox provides a struct that read and relay events to event bus.
package outbox

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"sync"
	"sync/atomic"
	"time"

	"github.com/go-logr/logr"
	"github.com/oklog/ulid/v2"

	novelv1 "play-ddd/proto/gen/go/contents/novel/v1"
	"play-ddd/utils/xerr"
	"play-ddd/utils/xslice"
)

// ErrRelayStopped is the error returned when the relay is stopped.
var ErrRelayStopped = errors.New("relay stopped")

// Relay relays events from the event store to the event bus.
type Relay struct {
	bus  EventBus
	repo EventRepo

	// Configuration
	maxFetch    int
	maxPub      int
	interval    time.Duration
	tickTimeout time.Duration
	instance    ulid.ULID

	// Runtime
	ticker *time.Ticker
	notify <-chan struct{}
	log    *logr.Logger

	mu       sync.Mutex
	running  bool
	ctx      context.Context
	cancelFn context.CancelCauseFunc
	wg       sync.WaitGroup

	// Health check
	lastTickAt atomic.Value // time.Time
	// lastTickErr atomic.Value // error
}

// NewRelay creates a new Relay with the given dependencies and options.
func NewRelay(
	bus EventBus,
	repo EventRepo,
	// log logr.Logger,
	// notify <-chan struct{},
	opts ...Option,
) *Relay {
	r := &Relay{
		bus:         bus,
		repo:        repo,
		maxFetch:    100,
		maxPub:      100,
		interval:    time.Second,
		tickTimeout: 10 * time.Second,
		notify:      make(<-chan struct{}),
		instance:    ulid.Make(),
	}

	for _, opt := range opts {
		opt(r)
	}

	r.ticker = time.NewTicker(r.interval)
	r.ticker.Stop() // Stop until Start() is called

	return r
}

// Start starts the relay. It is idempotent and concurrent-safe.
func (r *Relay) Start() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.running {
		return nil
	}

	r.ctx, r.cancelFn = context.WithCancelCause(context.Background())
	r.ticker.Reset(r.interval)
	r.running = true

	r.wg.Go(r.run)

	return nil
}

// Stop stops the relay gracefully. It is idempotent and concurrent-safe.
// It waits for the current batch to complete before returning.
func (r *Relay) Stop() error {
	r.mu.Lock()
	if !r.running {
		r.mu.Unlock()
		return nil
	}
	r.running = false
	r.mu.Unlock()

	r.ticker.Stop()
	r.cancelFn(ErrRelayStopped)
	r.wg.Wait()

	return nil
}

// Healthy checks if the relay is healthy (last tick was within expected time).
func (r *Relay) Healthy() bool {
	r.mu.Lock()
	running := r.running
	r.mu.Unlock()

	if !running {
		return true // not started is considered healthy
	}

	t, ok := r.lastTickAt.Load().(time.Time)
	if !ok {
		return true // just started, no tick yet
	}

	// Last tick should be within 2x interval
	return time.Since(t) < 2*r.interval
}

// Ready checks if the relay is ready (running and healthy).
func (r *Relay) Ready() bool {
	r.mu.Lock()
	running := r.running
	r.mu.Unlock()

	return running && r.Healthy()
}

// LastError returns the last tick error, if any.
// func (r *Relay) LastError() error {
// 	if err, ok := r.lastTickErr.Load().(error); ok {
// 		return err
// 	}
// 	return nil
// }

// run is the main loop that listens for triggers and processes events.
func (r *Relay) run() {
NEXT:
	select {
	case <-r.ctx.Done():
		r.log.Info("Relay stopped", "cause", context.Cause(r.ctx))
		return
	case _, ok := <-r.notify:
		if !ok {
			// Channel closed, set to nil to prevent busy loop
			r.notify = nil
			goto NEXT
		}
		r.tick()
	case <-r.ticker.C:
		r.tick()
	}
	goto NEXT
}

// tick performs a single tick and updates health check state.
func (r *Relay) tick() {
	r.lastTickAt.Store(time.Now())

	ctx, cancel := context.WithTimeoutCause(
		r.ctx,
		r.tickTimeout,
		fmt.Errorf("max duration per tick is %s", r.tickTimeout))
	defer cancel()

	err := r.repo.Process(r.processEvents(ctx))
	if err != nil {
		r.log.Error(err, "Relay tick failed.")
		return
	}

	r.log.V(1).Info("Relay ticked.")
}

func (r *Relay) processEvents(
	ctx context.Context) (
	context.Context, func(ef EventsBatch) error,
) {
	return ctx, func(eb EventsBatch) error {
	NEXT:
		events, err := eb.Poll(Arg{Limit: r.maxFetch})
		if err != nil {
			return err
		}

		if len(events) == 0 {
			return nil
		}

		err = r.pub(ctx, eb, events)

		switch {
		case err != nil:
			return err
		case len(events) < r.maxFetch:
			goto NEXT
		default:
			return nil
		}
	}
}

// pub publishes events in batches, maintaining order within aggregates.
// If an event fails, subsequent events in the same aggregate are skipped.
func (r *Relay) pub(
	ctx context.Context, eb EventsBatch, events []*novelv1.Event) (
	err error,
) {
	defer xerr.Expect(&err, "publish events")

	// Track failed aggregates to maintain order within aggregate
	failedAggs := make(map[ID]bool)
	results := make([]Result, len(events))
	toSend := make([]*novelv1.Event, 0, r.maxPub)
	toSendIdx := make([]int, 0, r.maxPub)

	for i, processed := 0, 0; processed < len(events); processed += i {
		var e *novelv1.Event
		for i, e = range events[processed:] {
			if len(toSend) >= r.maxPub {
				break
			}

			id, aid := e.GetId().Into(), e.GetAggregateId().Into()
			switch {
			case failedAggs[aid]:
				results[processed+i] = resultBlocked(id)
			default:
				toSend = append(toSend, e)
				toSendIdx = append(toSendIdx, processed+i)
			}
		}

		pubErr := r.bus.Pub(ctx, toSend)
		xslice.ForeachIdx(
			slices.All(toSend),
			func(i int, e *novelv1.Event) {
				if pubErr != nil {
					results[toSendIdx[i]] = resultFailed(e.GetId().Into(), pubErr)
					failedAggs[e.GetId().Into()] = true
				}

				results[toSendIdx[i]] = resultPublished(e.GetId().Into())
			})

		clear(toSend)
		clear(toSendIdx)
		toSend, toSendIdx = toSend[:0], toSendIdx[:0]
	}

	if err = eb.AdvanceCursor(results...); err != nil {
		return err
	}

	return nil
}

func resultBlocked(id ID) Result {
	return Result{
		ID:     id,
		Status: StatusBlocked,
		Reason: "One of previous events in same aggregate failed.",
	}
}

func resultFailed(id ID, err error) Result {
	return Result{
		ID:     id,
		Status: StatusFailed,
		Reason: err.Error(),
	}
}

func resultPublished(id ID) Result {
	return Result{
		ID:     id,
		Status: StatusPublished,
	}
}
