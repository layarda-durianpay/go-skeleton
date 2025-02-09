package kafka

import (
	"context"

	"github.com/durianpay/dpay-common/dckafka"
	"github.com/durianpay/dpay-common/logger"
	"github.com/segmentio/kafka-go"
)

type (
	BeforeFunc func(ctx context.Context) context.Context
	AfterFunc  func(ctx context.Context, msg kafka.Message, err error)
)

func Read(
	reader *kafka.Reader,
	opts ...Option,
) dckafka.Reader {
	cr := newConfig(opts...)

	return func(ctx context.Context) (msg kafka.Message, err error) {
		ctx = cr.beforeFunc(ctx)
		defer cr.afterFunc(ctx, msg, err)

		msg, err = reader.ReadMessage(ctx)
		if err != nil {
			logger.Errorw(
				ctx,
				"error reading from kafka",
				"error", err.Error(),
			)
			return
		}

		logger.Infow(ctx, "successfully received message from kafka")
		return
	}
}
