package sender

import (
	"context"

	"github.com/IBM/sarama"
)

type WsHub struct {
	producer sarama.SyncProducer
}

func NewWsHub(addrs []string) (*WsHub, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	producer, err := sarama.NewSyncProducer(addrs, config)
	if err != nil {
		return nil, err
	}

	return &WsHub{producer: producer}, nil
}

func (h *WsHub) Send(ctx context.Context, channel string, message []byte) error {
	_, _, err := h.producer.SendMessage(&sarama.ProducerMessage{
		Topic: channel,
		Value: sarama.StringEncoder(message),
	})
	return err
}

func (h *WsHub) Close() error {
	return h.producer.Close()
}
