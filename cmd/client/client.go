package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	pb "github.com/Egorpalan/grpc-easyp/pkg/api/chat/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func main() {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewChatAPIClient(conn)

	ctx := context.Background()
	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer my-secret-token")

	stream, err := client.Chat(ctx)
	if err != nil {
		log.Fatalf("failed to create stream: %v", err)
	}

	go func() {
		for {
			msg, err := stream.Recv()
			if err == io.EOF {
				return
			}
			if err != nil {
				log.Printf("recv error: %v", err)
				return
			}

			if msg.GetMessage() != nil {
				m := msg.GetMessage()
				log.Printf("Received: [%s] %s", m.CorrelationId, m.Text)
			} else if msg.GetError() != nil {
				log.Printf("Received error: %v", msg.GetError())
			}
		}
	}()

	messages := []string{
		"Привет, сервер!",
		"Как дела?",
		"Тестирую двунаправленный стриминг",
	}

	for i, text := range messages {
		msg := &pb.ChatMessage{
			Content: &pb.ChatMessage_Message{
				Message: &pb.Message{
					CorrelationId: fmt.Sprintf("msg-%d", i+1),
					Text:          text,
					Timestamp:     time.Now().Unix(),
				},
			},
		}

		if err := stream.Send(msg); err != nil {
			log.Fatalf("send error: %v", err)
		}

		time.Sleep(1 * time.Second)
	}

	time.Sleep(10 * time.Second)

	stream.CloseSend()
}
