package log

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
)

// New interceptorLogger adapts logr logger to interceptor logger.
// This code is simple enough to be copied and not imported.
func New(l logr.Logger) logging.Logger {
	return logging.LoggerFunc(
		func(_ context.Context, lvl logging.Level, msg string, fields ...any) {
			l := l.WithValues(fields...)
			switch lvl {
			case logging.LevelDebug:
				l.V(debugVerbosity).Info(msg)
			case logging.LevelInfo:
				l.V(infoVerbosity).Info(msg)
			case logging.LevelWarn:
				l.V(warnVerbosity).Info(msg)
			case logging.LevelError:
				l.V(errorVerbosity).Info(msg)
			default:
				panic(fmt.Sprintf("unknown level %v", lvl))
			}
		},
	)
}

const (
	debugVerbosity = 4
	infoVerbosity  = 2
	warnVerbosity  = 1
	errorVerbosity = 0
)
