package metrics_v1

import (
	"github.com/Egorpalan/grpc-easyp/internal/model/metrics"
	pb "github.com/Egorpalan/grpc-easyp/pkg/api/metrics/v1"
)

func convertMetricRequestToService(req *pb.UploadMetricsRequest) *metrics.Metric {
	if req == nil {
		return nil
	}

	return &metrics.Metric{
		Value: req.Value,
	}
}

func convertSummaryToProto(summary *metrics.UploadMetricsResponse) *pb.UploadMetricsResponse {
	if summary == nil {
		return nil
	}

	return &pb.UploadMetricsResponse{
		Sum:     summary.Sum,
		Average: summary.Average,
		Count:   summary.Count,
	}
}
