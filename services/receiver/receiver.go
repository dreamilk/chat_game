package receiver

import (
	"context"
	"os"
	"os/signal"

	"github.com/IBM/sarama"
	"go.uber.org/zap"

	"chat_game/log"
)

type Receiver struct {
	consumer sarama.Consumer
}

func NewReceiver(addrs []string) (*Receiver, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	consumer, err := sarama.NewConsumer(addrs, config)
	if err != nil {
		return nil, err
	}

	return &Receiver{consumer: consumer}, nil
}

func (r *Receiver) Consume(ctx context.Context, channel string, handler func(ctx context.Context, message []byte)) error {
	consumer, err := r.consumer.ConsumePartition(channel, 0, sarama.OffsetNewest)
	if err != nil {
		return err
	}
	defer consumer.Close()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	done := make(chan struct{})

	log.Info(ctx, "Consumer started")
	go func() {
		for {
			select {
			case msg := <-consumer.Messages():
				log.Info(ctx, "Received message", zap.String("topic", msg.Topic), zap.String("key", string(msg.Key)), zap.String("value", string(msg.Value)))
				handler(ctx, msg.Value)
			case err := <-consumer.Errors():
				log.Error(ctx, "Error consuming message", zap.Error(err))
			case <-signals:
				done <- struct{}{}
			}
		}
	}()

	log.Info(ctx, "Consumer stopped")

	<-done
	return nil
}

func (r *Receiver) Close() error {
	return r.consumer.Close()
}
