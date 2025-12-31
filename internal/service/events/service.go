package events

import (
	"context"
	"sync"

	"github.com/Egorpalan/grpc-easyp/internal/model/event"
	"github.com/Egorpalan/grpc-easyp/internal/service/notes"
)

type Service interface {
	notes.EventPublisher
	Subscribe(ctx context.Context) <-chan event.Event
}

type service struct {
	mu          sync.RWMutex
	subscribers map[chan event.Event]struct{}
}

func NewService() Service {
	return &service{
		subscribers: make(map[chan event.Event]struct{}),
	}
}

func (s *service) Publish(ctx context.Context, e event.Event) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for ch := range s.subscribers {
		select {
		case ch <- e:
		case <-ctx.Done():
			return
		default:
		}
	}
}

func (s *service) Subscribe(ctx context.Context) <-chan event.Event {
	ch := make(chan event.Event, 10)

	s.mu.Lock()
	s.subscribers[ch] = struct{}{}
	s.mu.Unlock()

	go func() {
		<-ctx.Done()
		s.mu.Lock()
		delete(s.subscribers, ch)
		close(ch)
		s.mu.Unlock()
	}()

	return ch
}
