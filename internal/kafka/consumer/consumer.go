package consumer

import (
	"context"
	"fmt"
	"time"

	"main/internal/common/interfaces"
	"main/internal/logger"
	"main/pkg/messages/proto"

	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"
)

type KafkaConsumerConfig struct {
	Brokers []string `yaml:"brokers"`
	Topic   string   `yaml:"topic"`
	GroupID string   `yaml:"group_id"`
}

type TaskHandler interface {
	HandleTask(ctx context.Context, task *messages.WrapperResponse) error
}

type Consumer struct {
	reader  *kafka.Reader
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
		reader:  r,
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
			logger := logger.GetLoggerFromContext(ctx)

			m, err := c.reader.ReadMessage(ctx)
			if err != nil {
				time.Sleep(1 * time.Second)
				continue
			}
			logger.Info("exec: Consumer\nmessage received", interfaces.LogField{
				Key:   "topic",
				Value: m.Topic,
			}, interfaces.LogField{
				Key:   "partition",
				Value: m.Partition,
			}, interfaces.LogField{
				Key:   "offset",
				Value: m.Offset,
			}, interfaces.LogField{
				Key:   "key",
				Value: string(m.Key),
			})
			var task messages.WrapperResponse
			if err := proto.Unmarshal(m.Value, &task); err != nil {

				logger.Error("exec: Consumer\nfailed to unmarshal message", interfaces.LogField{
					Key:   "error",
					Value: err.Error(),
				})
				continue
			}

			if err := c.handler.HandleTask(ctx, &task); err != nil {
				logger.Error("exec: Consumer\ntask handle error", interfaces.LogField{
					Key:   "error",
					Value: err.Error(),
				})
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
