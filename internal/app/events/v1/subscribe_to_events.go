package events_v1

import (
	pb "github.com/Egorpalan/grpc-easyp/pkg/api/events/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (i *Implementation) SubscribeToEvents(req *pb.SubscribeToEventsRequest, stream pb.EventsAPI_SubscribeToEventsServer) error {
	ctx := stream.Context()

	healthCheck := &pb.EventResponse{
		Event: &pb.EventResponse_HealthCheck{
			HealthCheck: &pb.HealthCheck{
				Message: "Connected to events stream",
			},
		},
	}

	if err := stream.Send(healthCheck); err != nil {
		return status.Error(codes.Internal, "failed to send health check")
	}

	eventCh := i.service.Subscribe(ctx)

	for {
		select {
		case <-ctx.Done():
			return nil

		case e, ok := <-eventCh:
			if !ok {
				return nil
			}

			eventResponse := convertEventToProto(e)
			if eventResponse == nil {
				continue
			}

			if err := stream.Send(eventResponse); err != nil {
				return status.Error(codes.Internal, "failed to send event")
			}
		}
	}
}
