package metrics_v1

import (
	"google.golang.org/grpc"

	pb "github.com/Egorpalan/grpc-easyp/pkg/api/metrics/v1"

	"github.com/Egorpalan/grpc-easyp/internal/service/metrics"
)

type Implementation struct {
	service metrics.Service
	pb.UnimplementedMetricsAPIServer
}

func New(service metrics.Service) *Implementation {
	return &Implementation{
		service: service,
	}
}

func (i *Implementation) RegisterServer(server *grpc.Server) {
	pb.RegisterMetricsAPIServer(server, i)
}
