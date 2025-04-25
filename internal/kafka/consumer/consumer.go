package consumer

import (
	"context"
	"log"
	"time"
    "fmt"

	"main/pkg/messages/proto"

	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"
)

type KafkaConsumerConfig struct {
	Brokers []string `yaml:"brokers"`
	Topic string `yaml:"topic"`
	GroupID string `yaml:"group_id"`
}

type TaskHandler interface {
	HandleTask(ctx context.Context, task *messages.WrapperResponse) error
}

type Consumer struct {
	reader *kafka.Reader
	handler TaskHandler
}

func NewConsumer(cfg KafkaConsumerConfig, taskHandler TaskHandler) (*Consumer, error) {
    if err := checkKafkaConnection(cfg.Brokers); err != nil {
        return nil, fmt.Errorf("kafka consumer connection error: %w", err)
    }

	r := kafka.NewReader(kafka.ReaderConfig{
        Brokers:     cfg.Brokers,
        Topic:       cfg.Topic,
        GroupID:     cfg.GroupID,
        StartOffset: kafka.FirstOffset,
    })

	return &Consumer{
		reader: r,
		handler: taskHandler,
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

func (c *Consumer) Start(ctx context.Context) {
	go func() {
        for {
            m, err := c.reader.ReadMessage(ctx)
            if err != nil {
                time.Sleep(1 * time.Second)
                continue
            }

            var task messages.WrapperResponse
            if err := proto.Unmarshal(m.Value, &task); err != nil {
                continue
            }

            if err := c.handler.HandleTask(ctx, &task); err != nil {
                log.Printf("Task handle failed: %v", err)
            }
        }
    }()
}

func (c *Consumer) Close() error {
	if err := c.reader.Close(); err != nil {
		return fmt.Errorf("failed to close Kafka consumer: %w", err)
	}
	return nil
}