package kafka

import (
	"context"
	"github.com/segmentio/kafka-go"
	"fmt"
)

type KafkaConfig struct {
	Brokers []string `yaml:"brokers"`
	Topic   string   `yaml:"topic"`
}

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(cfg KafkaConfig) (*Producer, error) {
	if err := checkKafkaConnection(cfg.Brokers); err != nil {
		return nil, fmt.Errorf("error connect to kafka: %s", err)
	}
	writer := &kafka.Writer{
		Addr:     kafka.TCP(cfg.Brokers...),
		Topic:    cfg.Topic,
		Balancer: &kafka.LeastBytes{},
	}

	return &Producer{writer: writer}, nil
}

func (p *Producer) SendMessage(ctx context.Context, key string, value []byte) error {
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