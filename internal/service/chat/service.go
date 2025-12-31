package chat

import (
	"context"
	"log/slog"
	"time"

	"github.com/Egorpalan/grpc-easyp/internal/model/chat"
	"github.com/google/uuid"
)

type Service interface {
	ProcessMessage(ctx context.Context, msg *chat.Message) (*chat.Message, error)
	GenerateNotification(ctx context.Context) *chat.Message
}

type service struct {
	logger *slog.Logger
}

func NewService(logger *slog.Logger) Service {
	return &service{
		logger: logger,
	}
}

func (s *service) ProcessMessage(ctx context.Context, msg *chat.Message) (*chat.Message, error) {
	ack := &chat.Message{
		CorrelationID: msg.CorrelationID,
		Text:          "Message received: " + msg.Text,
		Timestamp:     time.Now(),
	}

	s.logger.InfoContext(ctx,
		"message processed",
		"correlation_id", msg.CorrelationID,
		"text", msg.Text,
	)

	return ack, nil
}

func (s *service) GenerateNotification(ctx context.Context) *chat.Message {
	return &chat.Message{
		CorrelationID: uuid.New().String(),
		Text:          "System notification: Server is running",
		Timestamp:     time.Now(),
	}
}
