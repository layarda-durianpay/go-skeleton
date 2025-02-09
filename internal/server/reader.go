package server

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	commoncfg "github.com/durianpay/dpay-common/config"
	"github.com/durianpay/dpay-common/dckafka"
	"github.com/durianpay/dpay-common/logger"
	kafkahandler "github.com/layarda-durianpay/go-skeleton/internal/disburse/handler/kafka"
	commonkafka "github.com/layarda-durianpay/go-skeleton/pkg/common/kafka"
	"github.com/samber/lo"
	"github.com/segmentio/kafka-go"
	// "github.com/segmentio/kafka-go"
)

type readers struct {
	DisbursementReaders *kafka.Reader
}

func NewReader(ctx context.Context) (readers, error) {
	disbursementKafkaReader, err := initKafkaReader(disbursementCfg.GetDisbursementKafkaTopic())
	if err != nil {
		logger.Errorw(ctx, "error initializing disbursement kafka reader", "error", err.Error())
		return lo.Empty[readers](), err
	}

	return readers{
		DisbursementReaders: disbursementKafkaReader,
	}, nil
}

func StartReaders() error {
	ctx, cancel := context.WithCancel(context.Background())
	readers, err := NewReader(ctx)
	if err != nil {
		cancel()
		return err
	}

	handler := kafkahandler.NewDisbursementKafkaReader(&appObj)

	serveHelperServer(ctx)

	var wg sync.WaitGroup

	runDisbursementReader(
		ctx,
		readers.DisbursementReaders,
		handler,
		&wg,
	)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	defer close(c)

	oscall := <-c
	logger.Infof(ctx, "exiting %v \n", oscall)
	cancel()
	wg.Wait()

	return appObjCleanup()
}

func runDisbursementReader(
	ctx context.Context,
	disbursementKafkaReader *kafka.Reader,
	handler *kafkahandler.DisbursementKafkaReader,
	wg *sync.WaitGroup,
) error {

	wg.Add(1)

	go func() {
		defer wg.Done()

		// TODO: improve InitConsumer with fork Handler so it can handle multiple type at once
		consumer := dckafka.InitConsumer(
			dckafka.ConsumerEntity("disbursement"),
			dckafka.InitWorker(3),
			commonkafka.Read(
				disbursementKafkaReader,
				commonkafka.WithBeforeFunc(func(ctx context.Context) context.Context {
					return context.WithValue(ctx, "iseng-aja", "iseng-aja")
				}),
				commonkafka.WithAfterFunc(func(ctx context.Context, msg kafka.Message, err error) {
					logger.Debugw(ctx, "after func log")
				}),
			),
			handler.DisburseProcessor,
			disbursementKafkaReader.Close,
		)

		consumer.ConsumeMessage(ctx)
	}()

	return nil
}

func initKafkaReader(topic string) (reader *kafka.Reader, err error) {
	readerCfg := kafka.ReaderConfig{
		Brokers:         commoncfg.KafkaBrokerURLs(),
		GroupID:         commoncfg.KafkaClientID(),
		Topic:           topic,
		MinBytes:        10e3, // 10KB
		MaxBytes:        10e6, // 10MB
		MaxWait:         1 * time.Second,
		ReadLagInterval: -1,
	}
	reader = kafka.NewReader(readerCfg)
	return
}
