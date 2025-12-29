package events_v1

import (
	"github.com/Egorpalan/grpc-easyp/internal/model/event"
	pb "github.com/Egorpalan/grpc-easyp/pkg/api/events/v1"
)

func convertEventToProto(e event.Event) *pb.EventResponse {
	switch e.Type {
	case event.MessageTypeNoteCreated:
		return &pb.EventResponse{
			Event: &pb.EventResponse_NoteCreated{
				NoteCreated: &pb.NoteCreated{
					Id: e.NoteID,
				},
			},
		}
	default:
		return nil
	}
}
