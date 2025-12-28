package notes_v1

import (
	"context"
	"errors"
	"fmt"

	"github.com/Egorpalan/grpc-easyp/internal/model/exception"
	pb "github.com/Egorpalan/grpc-easyp/pkg/api/notes/v1"
)

func (i *Implementation) GetNote(ctx context.Context, req *pb.GetNoteRequest) (*pb.GetNoteResponse, error) {
	result, err := i.service.GetNote(ctx, req.Id)
	if err != nil {
		if errors.Is(err, exception.ErrNoteNotFound) {
			reason := fmt.Sprintf("Note with ID %s was searched but not found in DB", req.Id)
			return nil, exception.WrapErrorWithDetails(
				err,
				reason,
				"NOTE_NOT_FOUND",
			)
		}
		return nil, exception.WrapError(err)
	}

	return &pb.GetNoteResponse{
		Note: convertNoteToProto(result),
	}, nil
}
