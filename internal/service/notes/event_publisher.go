package notes

import (
	"context"

	"github.com/Egorpalan/grpc-easyp/internal/model/event"
)

type EventPublisher interface {
	Publish(ctx context.Context, e event.Event)
}
