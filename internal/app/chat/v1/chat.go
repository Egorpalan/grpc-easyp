package chat_v1

import (
	"io"
	"log/slog"
	"time"

	pb "github.com/Egorpalan/grpc-easyp/pkg/api/chat/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (i *Implementation) Chat(stream pb.ChatAPI_ChatServer) error {
	ctx := stream.Context()
	logger := slog.Default()

	sendCh := make(chan *pb.ChatMessage, 10)
	errCh := make(chan error, 1)

	go func() {
		for {
			req, err := stream.Recv()
			if err == io.EOF {
				close(sendCh)
				return
			}

			if err != nil {
				errCh <- err
				return
			}

			if msg := req.GetMessage(); msg != nil {
				domainMsg := convertMessageToService(msg)
				if domainMsg == nil {
					continue
				}

				ack, processErr := i.service.ProcessMessage(ctx, domainMsg)
				if processErr != nil {
					errorMsg := convertErrorToChatMessage(processErr)
					select {
					case sendCh <- errorMsg:
					case <-ctx.Done():
						return
					}
					continue
				}

				ackProto := convertMessageToChatMessage(ack)
				select {
				case sendCh <- ackProto:
				case <-ctx.Done():
					return
				}
			}
		}
	}()

	notificationTicker := time.NewTicker(5 * time.Second)
	defer notificationTicker.Stop()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-notificationTicker.C:
				notification := i.service.GenerateNotification(ctx)
				notificationProto := convertMessageToChatMessage(notification)

				select {
				case sendCh <- notificationProto:
				case <-ctx.Done():
					return
				}
			}
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case err := <-errCh:
			if err != nil {
				logger.ErrorContext(ctx, "error receiving message", "error", err)
				return status.Error(codes.Internal, "failed to receive message")
			}

		case msg, ok := <-sendCh:
			if !ok {
				return nil
			}

			if err := stream.Send(msg); err != nil {
				logger.ErrorContext(ctx, "error sending message", "error", err)
				return status.Error(codes.Internal, "failed to send message")
			}
		}
	}
}
