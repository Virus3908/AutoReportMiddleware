package kafka

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/segmentio/kafka-go"
)

type KafkaConfig struct {
	Brokers []string `yaml:"brokers"`
	Topic   string   `yaml:"topic"`
}

type KafkaMessage struct {
	Data        string `json:"data"`
	CallbackURL string `json:"callback_url"`
}

type Producer struct {
	writer      *kafka.Writer
	callbackURL string
}

func NewProducer(cfg KafkaConfig, callbackURL string) (*Producer, error) {
	if err := checkKafkaConnection(cfg.Brokers); err != nil {
		return nil, fmt.Errorf("error connect to kafka: %s", err)
	}
	writer := &kafka.Writer{
		Addr:     kafka.TCP(cfg.Brokers...),
		Topic:    cfg.Topic,
		Balancer: &kafka.LeastBytes{},
	}

	return &Producer{
		writer:      writer,
		callbackURL: callbackURL,
	}, nil
}

func (p *Producer) SendMessage(ctx context.Context, key string, data string) error {
	message := KafkaMessage{
		Data:        data,
		CallbackURL: p.callbackURL,
	}

	value, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal Kafka message: %w", err)
	}

	return p.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(key),
		Value: value,
	})
}

func checkKafkaConnection(brokers []string) error {
	var lastErr error
	for _, broker := range brokers {
		conn, err := kafka.Dial("tcp", broker)
		if err == nil {
			_ = conn.Close()
			return nil
		}
		lastErr = err
	}
	return fmt.Errorf("no Kafka brokers available: %v, last error: %w", brokers, lastErr)
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
