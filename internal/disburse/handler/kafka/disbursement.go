package kafkahandler

import (
	"context"
	"encoding/json"

	"github.com/durianpay/dpay-common/logger"
	"github.com/layarda-durianpay/go-skeleton/internal/disburse/app"
	"github.com/layarda-durianpay/go-skeleton/internal/disburse/app/command"
	"github.com/layarda-durianpay/go-skeleton/pkg/common/errors"
	commonkafka "github.com/layarda-durianpay/go-skeleton/pkg/common/kafka"
	schemakafka "github.com/layarda-durianpay/go-skeleton/pkg/common/schema"
	"github.com/segmentio/kafka-go"
)

// can use interface if neede

type DisbursementKafkaReader struct {
	app *app.Application
}

func NewDisbursementKafkaReader(apps *app.Application) *DisbursementKafkaReader {
	return &DisbursementKafkaReader{
		app: apps,
	}
}

func (r DisbursementKafkaReader) DisburseProcessor(ctx context.Context, message kafka.Message) error {
	var body commonkafka.ResponseMessage[schemakafka.DisburseKafkaRequest]

	err := json.Unmarshal(message.Value, &body)
	if err != nil {
		// TODO: this can moved to not logging on here
		logger.Errorw(
			ctx, "error unmarshalling kafka message",
			"error", err.Error(),
			"request", string(message.Value),
			"headers", message.Headers,
		)

		return errors.NewDpayError(
			err,
			"error unmarshalling kafka message",
			errors.DpayInternalError,
		)
	}

	err = r.app.Commands.Disburse.Handle(ctx, &command.DisburseParam{
		Amount: body.Data.Amount,
	})
	if err != nil {
		return errors.WrapDpayErrTrace(err)
	}

	return nil
}
