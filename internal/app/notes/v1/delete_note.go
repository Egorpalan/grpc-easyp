package notes_v1

import (
	"context"

	"github.com/Egorpalan/grpc-easyp/internal/model/exception"
	pb "github.com/Egorpalan/grpc-easyp/pkg/api/notes/v1"
)

func (i *Implementation) DeleteNote(ctx context.Context, req *pb.DeleteNoteRequest) (*pb.DeleteNoteResponse, error) {
	err := i.service.DeleteNote(ctx, req.Id)
	if err != nil {
		return nil, exception.WrapError(err)
	}

	return &pb.DeleteNoteResponse{}, nil
}
