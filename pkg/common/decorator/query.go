package decorator

import (
	"context"
)

type QueryHandler[Q any, R any] interface {
	Handle(ctx context.Context, q Q) (R, error)
}

func ApplyQueryDecorators[H any, R any](
	handler QueryHandler[H, R],
) QueryHandler[H, R] {
	return queryOTelDecorator[H, R]{
		base: queryLoggingDecorator[H, R]{
			base: queryErrorDecorator[H, R]{
				base: handler,
			},
		},
	}
}
