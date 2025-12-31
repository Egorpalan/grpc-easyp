package chat_v1

import (
	"time"

	"github.com/Egorpalan/grpc-easyp/internal/model/chat"
	pb "github.com/Egorpalan/grpc-easyp/pkg/api/chat/v1"
	"google.golang.org/grpc/status"
)

func convertMessageToService(msg *pb.Message) *chat.Message {
	if msg == nil {
		return nil
	}

	return &chat.Message{
		CorrelationID: msg.CorrelationId,
		Text:          msg.Text,
		Timestamp:     time.Unix(msg.Timestamp, 0),
	}
}

func convertMessageToProto(msg *chat.Message) *pb.Message {
	if msg == nil {
		return nil
	}

	return &pb.Message{
		CorrelationId: msg.CorrelationID,
		Text:          msg.Text,
		Timestamp:     msg.Timestamp.Unix(),
	}
}

func convertMessageToChatMessage(msg *chat.Message) *pb.ChatMessage {
	return &pb.ChatMessage{
		Content: &pb.ChatMessage_Message{
			Message: convertMessageToProto(msg),
		},
	}
}

func convertErrorToChatMessage(err error) *pb.ChatMessage {
	return &pb.ChatMessage{
		Content: &pb.ChatMessage_Error{
			Error: status.Convert(err).Proto(),
		},
	}
}
