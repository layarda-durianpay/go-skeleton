package interceptors

import (
	"context"
	"fmt"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"go.uber.org/zap"
)

func interceptorLogger(l *zap.SugaredLogger) logging.Logger {
	return logging.LoggerFunc(func(
		ctx context.Context,
		lvl logging.Level,
		msg string,
		fields ...any,
	) {
		f := make([]any, 0)

		for i := 0; i < len(fields); i += 2 {
			key := fields[i]
			value := fields[i+1]

			f = append(f, key, value)
		}

		switch lvl {
		case logging.LevelDebug:
			l.Debugw(msg, f...)
		case logging.LevelInfo:
			l.Infow(msg, f...)
		case logging.LevelWarn:
			l.Warnw(msg, f...)
		case logging.LevelError:
			l.Errorw(msg, f...)
		default:
			panic(fmt.Sprintf("unknown level %v", lvl))
		}
	})
}
