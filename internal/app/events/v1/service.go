package events_v1

import (
	"google.golang.org/grpc"

	pb "github.com/Egorpalan/grpc-easyp/pkg/api/events/v1"

	"github.com/Egorpalan/grpc-easyp/internal/service/events"
)

type Implementation struct {
	service events.Service
	pb.UnimplementedEventsAPIServer
}

func New(eventBus events.Service) *Implementation {
	return &Implementation{
		service: eventBus,
	}
}

func (i *Implementation) RegisterServer(server *grpc.Server) {
	pb.RegisterEventsAPIServer(server, i)
}
