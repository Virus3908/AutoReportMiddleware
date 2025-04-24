package producer

import (
	"context"
	"fmt"

	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"
)

type KafkaProducerConfig struct {
	Brokers []string `yaml:"brokers"`
	Topic string `yaml:"topic"`
}

type Producer struct {
	writer *kafka.Writer
	topic string
}

func NewProducer(cfg KafkaProducerConfig) (*Producer, error) {
	if err := checkKafkaConnection(cfg.Brokers); err != nil {
		return nil, fmt.Errorf("error connect to kafka: %s", err)
	}
	writer := &kafka.Writer{
		Addr:     kafka.TCP(cfg.Brokers...),
		Balancer: &kafka.LeastBytes{},
	}

	topics := cfg.Topic

	return &Producer{
		writer: writer,
		topic: topics,
	}, nil
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

func (p *Producer) SendMessage(ctx context.Context, key string, message proto.Message) error {
	value, err := proto.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	return p.writer.WriteMessages(ctx, kafka.Message{
		Topic: p.topic,
		Key:   []byte(key),
		Value: value,
	})
}
