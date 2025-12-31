package metrics_v1

import (
	"io"

	"github.com/Egorpalan/grpc-easyp/internal/model/metrics"
	pb "github.com/Egorpalan/grpc-easyp/pkg/api/metrics/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (i *Implementation) UploadMetrics(stream pb.MetricsAPI_UploadMetricsServer) error {
	ctx := stream.Context()
	var metricList []*metrics.Metric

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return status.Error(codes.Internal, "failed to receive metric")
		}

		metric := convertMetricRequestToService(req)
		if metric != nil {
			metricList = append(metricList, metric)
		}
	}

	if len(metricList) == 0 {
		return status.Error(codes.InvalidArgument, "no metrics received")
	}

	summary, err := i.service.UploadMetrics(ctx, metricList)
	if err != nil {
		return status.Error(codes.Internal, "failed to process metrics")
	}

	response := convertSummaryToProto(summary)
	if response == nil {
		return status.Error(codes.Internal, "failed to convert summary")
	}

	return stream.SendAndClose(response)
}
