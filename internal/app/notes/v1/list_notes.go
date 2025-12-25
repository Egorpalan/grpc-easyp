package notes_v1

import (
	"context"

	"github.com/Egorpalan/grpc-easyp/internal/model/exception"
	pb "github.com/Egorpalan/grpc-easyp/pkg/api/notes/v1"
)

func (i *Implementation) ListNotes(ctx context.Context, req *pb.ListNotesRequest) (*pb.ListNotesResponse, error) {
	result, err := i.service.ListNotes(ctx)
	if err != nil {
		return nil, exception.WrapError(err)
	}

	protoNotes := make([]*pb.Note, 0, len(result))
	for _, note := range result {
		protoNotes = append(protoNotes, convertNoteToProto(note))
	}

	return &pb.ListNotesResponse{
		Items: protoNotes,
	}, nil
}
