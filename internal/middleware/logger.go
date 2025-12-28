package middleware

import (
	"context"
	"log/slog"
	"time"

	"google.golang.org/grpc"
)

func LoggerInterceptor(logger *slog.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		start := time.Now()
		logger.Info("request started", "method", info.FullMethod)
		resp, err := handler(ctx, req)

		duration := time.Since(start)

		if err != nil {
			logger.Error("request failed",
				"method", info.FullMethod,
				"duration", duration,
				"error", err,
			)
		} else {
			logger.Info("request completed",
				"method", info.FullMethod,
				"duration", duration,
			)
		}

		return resp, err
	}
}
