package retry

import (
	"errors"
	"fmt"
	"strings"
)

// Error type represents list of errors in retry
type Error []error

// Error method return string representation of Error
// It is an implementation of error interface
func (e Error) Error() string {
	logWithNumber := make([]string, len(e))
	for i, l := range e {
		if l != nil {
			logWithNumber[i] = fmt.Sprintf("#%d: %s", i+1, l.Error())
		}
	}

	return fmt.Sprintf(
		"All attempts fail:\n%s",
		strings.Join(logWithNumber, "\n"),
	)
}

func (e Error) Is(target error) bool {
	for _, v := range e {
		if errors.Is(v, target) {
			return true
		}
	}
	return false
}

func (e Error) As(target any) bool {
	for _, v := range e {
		if errors.As(v, target) {
			return true
		}
	}
	return false
}

// Unwrap returns the list of errors that this Error is wrapping.
//
// This method implements the Unwrap() []error interface introduced in Go 1.20
// for multi-error unwrapping. This allows errors.Is and errors.As to traverse
// all wrapped errors, not just the last one.
//
// IMPORTANT: errors.Unwrap(err) will return nil because the standard library's
// errors.Unwrap function only calls Unwrap() error, not Unwrap() []error.
// This is the same behavior as errors.Join in Go 1.20.
//
// Example - Use errors.Is to check for specific errors:
//
//	err := retry.New(retry.Attempts(3)).Do(func() error {
//		return os.ErrNotExist
//	})
//	if errors.Is(err, os.ErrNotExist) {
//		// Handle not exist error
//	}
//
// Example - Use errors.As to extract error details:
//
//	var pathErr *fs.PathError
//	if errors.As(err, &pathErr) {
//		fmt.Println("Failed at path:", pathErr.Path)
//	}
//
// Example - Get the last error directly (for migration):
//
//	if retryErr, ok := err.(retry.Error); ok {
//		lastErr := retryErr.LastError()
//	}
//
// See also: LastError() for direct access to the last error.
func (e Error) Unwrap() []error {
	return e
}

// WrappedErrors returns the list of errors that this Error is wrapping.
// It is an implementation of the `errwrap.Wrapper` interface
// in package [errwrap](https://github.com/hashicorp/errwrap) so that
// `retry.Error` can be used with that library.
func (e Error) WrappedErrors() []error {
	return e
}

// LastError returns the last error in the error list.
//
// This is a convenience method for users migrating from retry-go v4.x where
// errors.Unwrap(err) returned the last error. In v5.0.0, errors.Unwrap(err)
// returns nil due to the switch to Unwrap() []error for Go 1.20 compatibility.
//
// Migration example:
//
//	// v4.x code:
//	lastErr := errors.Unwrap(retryErr)
//
//	// v5.0.0 code (option 1 - recommended):
//	if errors.Is(retryErr, specificError) { ... }
//
//	// v5.0.0 code (option 2 - if you need the last error):
//	lastErr := retryErr.(retry.Error).LastError()
//
// Note: Using errors.Is or errors.As is preferred as they check ALL wrapped
// errors, not just the last one.
func (e Error) LastError() error {
	if len(e) == 0 {
		return nil
	}
	return e[len(e)-1]
}

type transientError struct{ error }

func Transient(e error) error {
	if e == nil {
		return nil
	}

	return transientError{e}
}

func (e transientError) Error() string {
	if e.error == nil {
		return "<nil>"
	}
	return fmt.Sprintf(`transient error: %s`, e.error.Error())
}

func (e transientError) Unwrap() error { return e.error }

func (e transientError) Is(err error) bool {
	_, ok := err.(transientError)
	return ok
}

func IsTransient(e error) bool {
	return errors.Is(e, transientError{})
}

func unpackTransientError(e error) error {
	if te, ok := e.(transientError); ok {
		return te.error
	}
	return e
}
