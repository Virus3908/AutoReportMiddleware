package kafka

import (
	"context"
	"fmt"
	"main/internal/models"

	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"
)

type KafkaConfig struct {
	Brokers []string `yaml:"brokers"`
}

type Producer struct {
	writer *kafka.Writer
	topics map[models.TaskType]string
}

func NewProducer(cfg KafkaConfig) (*Producer, error) {
	if err := checkKafkaConnection(cfg.Brokers); err != nil {
		return nil, fmt.Errorf("error connect to kafka: %s", err)
	}
	writer := &kafka.Writer{
		Addr:     kafka.TCP(cfg.Brokers...),
		Balancer: &kafka.LeastBytes{},
	}

	topics := map[models.TaskType]string{
		models.ConvertTask:    "convert",
		models.DiarizeTask:    "diarize",
		models.TranscribeTask: "transcribe",
	}

	return &Producer{
		writer: writer,
		topics: topics,
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

func (p *Producer) SendMessage(ctx context.Context, taskType models.TaskType, key string, message proto.Message) error {
	topic, ok := p.topics[taskType]
	if !ok {
		return fmt.Errorf("topic not found for task type: %d", taskType)
	}

	value, err := proto.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	return p.writer.WriteMessages(ctx, kafka.Message{
		Topic: topic,
		Key:   []byte(key),
		Value: value,
	})
}
