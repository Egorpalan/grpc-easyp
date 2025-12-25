package middleware

import (
	"context"
	"log/slog"

	govalidator "buf.build/go/protovalidate"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/protovalidate"
	"google.golang.org/grpc"
)

func ValidationInterceptor(logger *slog.Logger) grpc.UnaryServerInterceptor {
	validator, err := govalidator.New()
	if err != nil {
		logger.Error("failed to initialize validator", "error", err)
		return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
			return handler(ctx, req)
		}
	}

	return protovalidate.UnaryServerInterceptor(validator)
}
