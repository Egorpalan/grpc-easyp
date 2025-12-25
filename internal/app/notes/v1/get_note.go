package notes_v1

import (
	"context"

	"github.com/Egorpalan/grpc-easyp/internal/model/exception"
	pb "github.com/Egorpalan/grpc-easyp/pkg/api/notes/v1"
)

func (i *Implementation) GetNote(ctx context.Context, req *pb.GetNoteRequest) (*pb.GetNoteResponse, error) {
	result, err := i.service.GetNote(ctx, req.Id)
	if err != nil {
		return nil, exception.WrapError(err)
	}

	return &pb.GetNoteResponse{
		Note: convertNoteToProto(result),
	}, nil
}
