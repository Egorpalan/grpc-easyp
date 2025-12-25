package notes_v1

import (
	"github.com/Egorpalan/grpc-easyp/internal/model/notes"
	pb "github.com/Egorpalan/grpc-easyp/pkg/api/notes/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func convertNoteToProto(note *notes.Note) *pb.Note {
	if note == nil {
		return nil
	}

	var id string
	if note.ID != nil {
		id = *note.ID
	}

	return &pb.Note{
		Id:          id,
		Title:       note.Title,
		Description: note.Description,
		CreatedAt:   timestamppb.New(note.CreatedAt),
		UpdatedAt:   timestamppb.New(note.UpdatedAt),
	}
}
