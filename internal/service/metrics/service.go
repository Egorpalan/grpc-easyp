package metrics

import (
	"context"
	"log/slog"

	"github.com/Egorpalan/grpc-easyp/internal/model/metrics"
)

type Service interface {
	UploadMetrics(ctx context.Context, metricList []*metrics.Metric) (*metrics.UploadMetricsResponse, error)
}

type service struct {
	logger *slog.Logger
}

func NewService(logger *slog.Logger) Service {
	return &service{
		logger: logger,
	}
}

func (s *service) UploadMetrics(ctx context.Context, metricList []*metrics.Metric) (*metrics.UploadMetricsResponse, error) {
	if len(metricList) == 0 {
		return &metrics.UploadMetricsResponse{
			Sum:     0,
			Average: 0,
			Count:   0,
		}, nil
	}
	var sum float64
	count := int64(len(metricList))

	for _, m := range metricList {
		sum += m.Value
	}

	average := sum / float64(count)

	result := &metrics.UploadMetricsResponse{
		Sum:     sum,
		Average: average,
		Count:   count,
	}

	s.logger.InfoContext(ctx,
		"processed metrics",
		"count", count,
		"sum", sum,
		"average", average,
	)

	return result, nil
}
