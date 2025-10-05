package config

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/segmentio/kafka-go"
)

type KafkaConfig struct {
	Brokers []string
	GroupID string
	Topics  []string
}

func GetKafkaConfig() *KafkaConfig {
	brokers := os.Getenv("KAFKA_BROKERS")
	if brokers == "" {
		brokers = "localhost:9092"
	}

	groupID := os.Getenv("KAFKA_GROUP_ID")
	if groupID == "" {
		groupID = "go-fiber-clean"
	}

	topics := os.Getenv("KAFKA_TOPICS")
	var topicsList []string
	if topics == "" {
		topicsList = []string{"user-events", "auth-events"}
	} else {
		topicsList = strings.Split(topics, ",")
	}

	return &KafkaConfig{
		Brokers: []string{brokers},
		GroupID: groupID,
		Topics:  topicsList,
	}
}

func GetKafkaWriter(topic string) *kafka.Writer {
	config := GetKafkaConfig()

	batchSize := getEnvInt("KAFKA_BATCH_SIZE", 10)
	batchTimeout := getEnvDuration("KAFKA_BATCH_TIMEOUT", "1s")
	async := getEnvBool("KAFKA_ASYNC", false)

	return &kafka.Writer{
		Addr:         kafka.TCP(config.Brokers...),
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		BatchSize:    batchSize,
		BatchTimeout: batchTimeout,
		Async:        async,
	}
}

func GetKafkaReader(topic string) *kafka.Reader {
	config := GetKafkaConfig()

	minBytes := getEnvInt("KAFKA_MIN_BYTES", 10)
	maxBytes := getEnvInt("KAFKA_MAX_BYTES", 10)
	maxWait := getEnvDuration("KAFKA_MAX_WAIT", "1s")
	readBatchTimeout := getEnvDuration("KAFKA_READ_BATCH_TIMEOUT", "10s")

	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:          config.Brokers,
		GroupID:          config.GroupID,
		Topic:            topic,
		MinBytes:         minBytes * 1024,
		MaxBytes:         maxBytes * 1024 * 1024,
		MaxWait:          maxWait,
		ReadBatchTimeout: readBatchTimeout,
		StartOffset:      kafka.LastOffset,
	})
}

func TestKafkaConnection() error {
	config := GetKafkaConfig()

	conn, err := kafka.Dial("tcp", config.Brokers[0])
	if err != nil {
		return fmt.Errorf("failed to connect to Kafka: %w", err)
	}
	defer conn.Close()

	_, err = conn.Controller()
	if err != nil {
		return fmt.Errorf("failed to get Kafka controller: %w", err)
	}

	return nil
}

func CloseKafkaWriter(writer *kafka.Writer) error {
	if writer == nil {
		return nil
	}
	return writer.Close()
}

func CloseKafkaReader(reader *kafka.Reader) error {
	if reader == nil {
		return nil
	}
	return reader.Close()
}

func ProduceMessage(ctx context.Context, writer *kafka.Writer, key []byte, value []byte) error {
	if writer == nil {
		return fmt.Errorf("kafka writer is not initialized")
	}

	err := writer.WriteMessages(ctx, kafka.Message{
		Key:   key,
		Value: value,
		Time:  time.Now(),
	})

	if err != nil {
		return fmt.Errorf("failed to produce message: %w", err)
	}

	return nil
}

func ConsumeMessages(ctx context.Context, reader *kafka.Reader, handler func(message kafka.Message) error) error {
	if reader == nil {
		return fmt.Errorf("kafka reader is not initialized")
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			m, err := reader.ReadMessage(ctx)
			if err != nil {
				if err == context.Canceled {
					return nil
				}
				return fmt.Errorf("failed to consume message: %w", err)
			}

			if err := handler(m); err != nil {
				return fmt.Errorf("message handler failed: %w", err)
			}
		}
	}
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

func getEnvDuration(key string, defaultValue string) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	if duration, err := time.ParseDuration(defaultValue); err == nil {
		return duration
	}
	return 1 * time.Second
}
