package notes_v1

import (
	"context"

	"github.com/Egorpalan/grpc-easyp/internal/model/exception"
	pb "github.com/Egorpalan/grpc-easyp/pkg/api/notes/v1"
)

func (i *Implementation) CreateOrUpdateNote(ctx context.Context, req *pb.CreateOrUpdateNoteRequest) (*pb.CreateOrUpdateNoteResponse, error) {
	var id *string
	if req.Id != nil {
		id = req.Id
	}

	result, err := i.service.CreateOrUpdateNote(ctx, id, req.Title, req.Description)
	if err != nil {
		return nil, exception.WrapError(err)
	}

	return &pb.CreateOrUpdateNoteResponse{
		Note: convertNoteToProto(result),
	}, nil
}
