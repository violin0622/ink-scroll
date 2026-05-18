package retry

import (
	"context"
	"errors"
	"maps"
	"sync"
	"time"
)

// Function signature of retryable function
type RetryableFunc func() error

// Function signature of retryable function with data
type RetryableWithFunc[T any] func() (T, error)

// Stats reports cumulative retry execution metrics for a retrier.
type Stats struct {
	Retries      uint64
	Successes    uint64
	Failures     uint64
	TotalElapsed time.Duration
}

// Runs returns the number of completed retry runs.
func (s Stats) Runs() uint64 {
	return s.Successes + s.Failures
}

// Attempts returns the number of logical attempts across all runs.
func (s Stats) Attempts() uint64 {
	return s.Runs() + s.Retries
}

// Retrier is for retry operations that return only an error.
type Retrier struct {
	*retrierCore
}

// RetrierWithData is for retry operations that return data and an error.
type RetrierWithData[T any] struct {
	*retrierCore
}

// New creates a new Retrier with the given options.
// The returned Retrier can be safely reused across multiple retry operations.
func New(opts ...Option) *Retrier {
	return &Retrier{retrierCore: newRetrieerCore(opts...)}
}

// NewWithData creates a new RetrierWithData[T] with the given options.
// The returned retrier can be safely reused across multiple retry operations.
func NewWithData[T any](opts ...Option) *RetrierWithData[T] {
	return &RetrierWithData[T]{retrierCore: newRetrieerCore(opts...)}
}

// Do executes the retryable function using this Retrier's configuration.
func (r *Retrier) Do(fn RetryableFunc) error {
	retryableWith := func() (any, error) { return nil, fn() }
	_, err := doWith(r.retrierCore, retryableWith)
	return err
}

// Stats returns a cumulative snapshot of this Retrier's executions.
func (r *Retrier) Stats() Stats {
	return r.retrierCore.statsSnapshot()
}

// Do executes the retryable function using this RetrierWithData's configuration.
func (r *RetrierWithData[T]) Do(retryableFunc RetryableWithFunc[T]) (T, error) {
	return doWith(r.retrierCore, retryableFunc)
}

// Stats returns a cumulative snapshot of this RetrierWithData's executions.
func (r *RetrierWithData[T]) Stats() Stats {
	return r.retrierCore.statsSnapshot()
}

func doWith[T any](r *retrierCore, retryableFunc RetryableWithFunc[T]) (result T, err error) {
	var n uint
	var attempts uint64
	success := false
	start := time.Now()

	defer func() {
		r.recordStats(attempts, success, time.Since(start))
	}()

	if err := context.Cause(r.context); err != nil {
		return result, err
	}

	// Setting r.attempts to 0 means we'll retry until we succeed
	var lastErr error
	if r.attempts == 0 {
		for {
			attempts++
			t, retryErr := retryableFunc()
			if retryErr == nil {
				success = true
				return t, nil
			}

			if !r.retryIf(retryErr) {
				return result, unpackTransientError(retryErr)
			}

			lastErr = unpackTransientError(retryErr)

			r.onRetry(n, retryErr)
			n++
			select {
			case <-r.timer.After(r.computeDelay(n, retryErr)):
			case <-r.context.Done():
				if r.wrapContextErrorWithLastError {
					return result, Error{context.Cause(r.context), lastErr}
				}
				return result, context.Cause(r.context)
			}
		}
	}

	errorLog := Error{}

	attemptsForErrorCopy := make(map[error]uint, len(r.attemptsForError))
	maps.Copy(attemptsForErrorCopy, r.attemptsForError)

shouldRetry:
	for {
		attempts++
		t, retryErr := retryableFunc()
		if retryErr == nil {
			success = true
			return t, nil
		}

		errorLog = append(errorLog, unpackTransientError(retryErr))

		if !r.retryIf(retryErr) {
			break
		}

		r.onRetry(n, retryErr)

		for errToCheck, attemptsForThisError := range attemptsForErrorCopy {
			if errors.Is(retryErr, errToCheck) {
				attemptsForThisError--
				attemptsForErrorCopy[errToCheck] = attemptsForThisError
				if attemptsForThisError <= 0 {
					break shouldRetry
				}
			}
		}

		// if this is last attempt - don't wait
		if n == r.attempts-1 {
			break shouldRetry
		}
		n++
		select {
		case <-r.timer.After(r.computeDelay(n, retryErr)):
		case <-r.context.Done():
			if r.lastErrorOnly {
				return result, context.Cause(r.context)
			}

			return result, append(errorLog, context.Cause(r.context))
		}
	}

	if r.lastErrorOnly {
		return result, errorLog[len(errorLog)-1]
	}
	return result, errorLog
}

func (r *retrierCore) computeDelay(n uint, err error) time.Duration {
	delayTime := r.delayType(n, err, r)
	if r.maxDelay > 0 && delayTime > r.maxDelay {
		delayTime = r.maxDelay
	}
	return delayTime
}

// retrierCore holds the core configuration and business logic for retry operations.
// this is then used by Retrier and RetrierWithData public APIs
type retrierCore struct {
	attempts                      uint
	attemptsForError              map[error]uint
	delay                         time.Duration
	maxDelay                      time.Duration
	maxJitter                     time.Duration
	onRetry                       OnRetryFunc
	retryIf                       RetryIfFunc
	delayType                     DelayTypeFunc
	lastErrorOnly                 bool
	context                       context.Context
	timer                         Timer
	wrapContextErrorWithLastError bool

	statsMu sync.Mutex
	stats   Stats

	maxBackOffN uint // pre-computed for BackOffDelay, immutable after New()
}

func (r *retrierCore) recordStats(attempts uint64, success bool, elapsed time.Duration) {
	r.statsMu.Lock()
	defer r.statsMu.Unlock()

	if success {
		r.stats.Successes++
	} else {
		r.stats.Failures++
	}
	if attempts > 0 {
		r.stats.Retries += attempts - 1
	}
	r.stats.TotalElapsed += elapsed
}

func (r *retrierCore) statsSnapshot() Stats {
	r.statsMu.Lock()
	defer r.statsMu.Unlock()

	return r.stats
}

// Delay implements DelayContext
func (r *retrierCore) Delay() time.Duration {
	return r.delay
}

// MaxJitter implements DelayContext
func (r *retrierCore) MaxJitter() time.Duration {
	return r.maxJitter
}

// MaxBackOffN implements DelayContext
func (r *retrierCore) MaxBackOffN() uint {
	return r.maxBackOffN
}

// MaxDelay implements DelayContext
func (r *retrierCore) MaxDelay() time.Duration {
	return r.maxDelay
}
