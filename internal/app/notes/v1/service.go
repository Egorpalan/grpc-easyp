package notes_v1

import (
	"google.golang.org/grpc"

	pb "github.com/Egorpalan/grpc-easyp/pkg/api/notes/v1"

	"github.com/Egorpalan/grpc-easyp/internal/service/notes"
)

type Implementation struct {
	service notes.Service
	pb.UnimplementedNoteAPIServer
}

func New(service notes.Service) *Implementation {
	return &Implementation{
		service: service,
	}
}

func (i *Implementation) RegisterServer(server *grpc.Server) {
	pb.RegisterNoteAPIServer(server, i)
}
