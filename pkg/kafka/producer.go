package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ranggakrisnaa/go-fiber-starter/config"
	"github.com/segmentio/kafka-go"
)

type KafkaProducer interface {
	ProduceUserCreatedEvent(ctx context.Context, event UserCreatedEvent) error
	ProduceUserUpdatedEvent(ctx context.Context, event UserUpdatedEvent) error
	ProduceUserDeletedEvent(ctx context.Context, event UserDeletedEvent) error
	ProduceAuthLoginEvent(ctx context.Context, event AuthLoginEvent) error
	ProduceAuthLogoutEvent(ctx context.Context, event AuthLogoutEvent) error
	Close() error
}

type kafkaProducer struct {
	userWriter *kafka.Writer
	authWriter *kafka.Writer
}

type UserCreatedEvent struct {
	UserID    string `json:"user_id"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	CreatedAt string `json:"created_at"`
}

type UserUpdatedEvent struct {
	UserID    string `json:"user_id"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	UpdatedAt string `json:"updated_at"`
}

type UserDeletedEvent struct {
	UserID    string `json:"user_id"`
	Email     string `json:"email"`
	DeletedAt string `json:"deleted_at"`
}

type AuthLoginEvent struct {
	UserID    string `json:"user_id"`
	Email     string `json:"email"`
	IP        string `json:"ip"`
	UserAgent string `json:"user_agent"`
	LoginAt   string `json:"login_at"`
}

type AuthLogoutEvent struct {
	UserID    string `json:"user_id"`
	Email     string `json:"email"`
	IP        string `json:"ip"`
	UserAgent string `json:"user_agent"`
	LogoutAt  string `json:"logout_at"`
	SessionID string `json:"session_id"`
}

func NewKafkaProducer() (KafkaProducer, error) {
	userWriter := config.GetKafkaWriter("user-events")
	authWriter := config.GetKafkaWriter("auth-events")

	return &kafkaProducer{
		userWriter: userWriter,
		authWriter: authWriter,
	}, nil
}

func (p *kafkaProducer) ProduceUserCreatedEvent(ctx context.Context, event UserCreatedEvent) error {
	return p.produceEvent(ctx, p.userWriter, "user-created", event)
}

func (p *kafkaProducer) ProduceUserUpdatedEvent(ctx context.Context, event UserUpdatedEvent) error {
	return p.produceEvent(ctx, p.userWriter, "user-updated", event)
}

func (p *kafkaProducer) ProduceUserDeletedEvent(ctx context.Context, event UserDeletedEvent) error {
	return p.produceEvent(ctx, p.userWriter, "user-deleted", event)
}

func (p *kafkaProducer) ProduceAuthLoginEvent(ctx context.Context, event AuthLoginEvent) error {
	return p.produceEvent(ctx, p.authWriter, "auth-login", event)
}

func (p *kafkaProducer) ProduceAuthLogoutEvent(ctx context.Context, event AuthLogoutEvent) error {
	return p.produceEvent(ctx, p.authWriter, "auth-logout", event)
}

func (p *kafkaProducer) produceEvent(ctx context.Context, writer *kafka.Writer, eventType string, event interface{}) error {
	key := []byte(eventType)

	value, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	headers := []kafka.Header{
		{
			Key:   "event-type",
			Value: []byte(eventType),
		},
		{
			Key:   "timestamp",
			Value: []byte(time.Now().Format(time.RFC3339)),
		},
	}

	message := kafka.Message{
		Key:     key,
		Value:   value,
		Headers: headers,
		Time:    time.Now(),
	}

	return writer.WriteMessages(ctx, message)
}

func (p *kafkaProducer) Close() error {
	if p.userWriter != nil {
		if err := p.userWriter.Close(); err != nil {
			return fmt.Errorf("failed to close user writer: %w", err)
		}
	}

	if p.authWriter != nil {
		if err := p.authWriter.Close(); err != nil {
			return fmt.Errorf("failed to close auth writer: %w", err)
		}
	}

	return nil
}
