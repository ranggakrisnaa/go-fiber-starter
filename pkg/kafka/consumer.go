package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/ranggakrisnaa/go-fiber-starter/config"
	"github.com/segmentio/kafka-go"
)

type KafkaConsumer interface {
	StartUserEventConsumer(ctx context.Context, handler UserEventHandler) error
	StartAuthEventConsumer(ctx context.Context, handler AuthEventHandler) error
	Close() error
}

type kafkaConsumer struct {
	userReader *kafka.Reader
	authReader *kafka.Reader
}

type UserEventHandler interface {
	HandleUserCreated(event UserCreatedEvent) error
	HandleUserUpdated(event UserUpdatedEvent) error
	HandleUserDeleted(event UserDeletedEvent) error
}

type AuthEventHandler interface {
	HandleAuthLogin(event AuthLoginEvent) error
	HandleAuthLogout(event AuthLogoutEvent) error
}

func NewKafkaConsumer() (KafkaConsumer, error) {
	userReader := config.GetKafkaReader("user-events")
	authReader := config.GetKafkaReader("auth-events")

	return &kafkaConsumer{
		userReader: userReader,
		authReader: authReader,
	}, nil
}

func (c *kafkaConsumer) StartUserEventConsumer(ctx context.Context, handler UserEventHandler) error {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				msg, err := c.userReader.ReadMessage(ctx)
				if err != nil {
					if err == context.Canceled {
						return
					}
					log.Printf("Error reading user event: %v", err)
					continue
				}

				if err := c.handleUserMessage(msg, handler); err != nil {
					log.Printf("Error handling user message: %v", err)
				}
			}
		}
	}()

	return nil
}

func (c *kafkaConsumer) StartAuthEventConsumer(ctx context.Context, handler AuthEventHandler) error {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				msg, err := c.authReader.ReadMessage(ctx)
				if err != nil {
					if err == context.Canceled {
						return
					}
					log.Printf("Error reading auth event: %v", err)
					continue
				}

				if err := c.handleAuthMessage(msg, handler); err != nil {
					log.Printf("Error handling auth message: %v", err)
				}
			}
		}
	}()

	return nil
}

func (c *kafkaConsumer) handleUserMessage(msg kafka.Message, handler UserEventHandler) error {
	eventType := c.getEventType(msg)

	switch eventType {
	case "user-created":
		var event UserCreatedEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			return fmt.Errorf("failed to unmarshal user created event: %w", err)
		}
		return handler.HandleUserCreated(event)
	case "user-updated":
		var event UserUpdatedEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			return fmt.Errorf("failed to unmarshal user updated event: %w", err)
		}
		return handler.HandleUserUpdated(event)
	case "user-deleted":
		var event UserDeletedEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			return fmt.Errorf("failed to unmarshal user deleted event: %w", err)
		}
		return handler.HandleUserDeleted(event)
	default:
		log.Printf("Unknown user event type: %s", eventType)
		return nil
	}
}

func (c *kafkaConsumer) handleAuthMessage(msg kafka.Message, handler AuthEventHandler) error {
	eventType := c.getEventType(msg)

	switch eventType {
	case "auth-login":
		var event AuthLoginEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			return fmt.Errorf("failed to unmarshal auth login event: %w", err)
		}
		return handler.HandleAuthLogin(event)
	case "auth-logout":
		var event AuthLogoutEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			return fmt.Errorf("failed to unmarshal auth logout event: %w", err)
		}
		return handler.HandleAuthLogout(event)
	default:
		log.Printf("Unknown auth event type: %s", eventType)
		return nil
	}
}

func (c *kafkaConsumer) getEventType(msg kafka.Message) string {
	for _, header := range msg.Headers {
		if header.Key == "event-type" {
			return string(header.Value)
		}
	}
	return "unknown"
}

func (c *kafkaConsumer) Close() error {
	if c.userReader != nil {
		if err := c.userReader.Close(); err != nil {
			return fmt.Errorf("failed to close user reader: %w", err)
		}
	}

	if c.authReader != nil {
		if err := c.authReader.Close(); err != nil {
			return fmt.Errorf("failed to close auth reader: %w", err)
		}
	}

	return nil
}
