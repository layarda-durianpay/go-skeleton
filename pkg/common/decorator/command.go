package decorator

import "context"

type CommandHandler[C any] interface {
	Handle(ctx context.Context, cmd C) error
}

func ApplyCommandDecorators[H any](
	handler CommandHandler[H],
) CommandHandler[H] {
	return commandOTelDecorator[H]{
		base: commandLoggingDecorator[H]{
			base: commandErrorDecorator[H]{
				base: handler,
			},
		},
	}
}
