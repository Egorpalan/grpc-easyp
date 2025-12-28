package middleware

import (
	"context"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	authorizationHeader = "authorization"
	bearerPrefix        = "Bearer "
	validToken          = "my-secret-token"
)

func AuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}

	authHeaders := md.Get(authorizationHeader)
	if len(authHeaders) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "authorization token is not provided")
	}

	authHeader := authHeaders[0]
	if !strings.HasPrefix(authHeader, bearerPrefix) {
		return nil, status.Errorf(codes.Unauthenticated, "authorization token is not provided")
	}

	token := strings.TrimPrefix(authHeader, bearerPrefix)
	if token != validToken {
		return nil, status.Errorf(codes.Unauthenticated, "authorization token is not provided")
	}

	return handler(ctx, req)
}
