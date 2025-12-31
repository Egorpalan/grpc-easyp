package chat_v1

import (
	"google.golang.org/grpc"

	pb "github.com/Egorpalan/grpc-easyp/pkg/api/chat/v1"

	"github.com/Egorpalan/grpc-easyp/internal/service/chat"
)

type Implementation struct {
	service chat.Service
	pb.UnimplementedChatAPIServer
}

func New(service chat.Service) *Implementation {
	return &Implementation{
		service: service,
	}
}

func (i *Implementation) RegisterServer(server *grpc.Server) {
	pb.RegisterChatAPIServer(server, i)
}
