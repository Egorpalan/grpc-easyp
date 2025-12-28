package app

import (
	"context"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/Egorpalan/grpc-easyp/internal/config"
	"github.com/Egorpalan/grpc-easyp/internal/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
)

// ServerOptions опции для создания gRPC сервера
type ServerOptions struct {
	Config            *config.Config
	Logger            *slog.Logger
	UnaryInterceptors []grpc.UnaryServerInterceptor
	EnableReflection  bool
	EnableValidation  bool
}

// NewServer создает новый gRPC сервер с настроенными опциями
func NewServer(opts ServerOptions) (*grpc.Server, error) {
	cfg := opts.Config

	kaep := keepalive.EnforcementPolicy{
		MinTime:             cfg.ServerConfig.KeepAliveTime,
		PermitWithoutStream: true,
	}

	kasp := keepalive.ServerParameters{
		MaxConnectionIdle:     cfg.ServerConfig.MaxConnectionIdle,
		MaxConnectionAge:      cfg.ServerConfig.MaxConnectionAge,
		MaxConnectionAgeGrace: cfg.ServerConfig.MaxConnectionAgeGrace,
		Time:                  cfg.ServerConfig.KeepAliveTime,
		Timeout:               cfg.ServerConfig.KeepAliveTimeout,
	}

	interceptors := make([]grpc.UnaryServerInterceptor, 0)
	interceptors = append(interceptors, opts.UnaryInterceptors...)

	if opts.EnableValidation {
		interceptors = append(interceptors, middleware.ValidationInterceptor(opts.Logger))
	}

	serverOpts := []grpc.ServerOption{
		grpc.MaxConcurrentStreams(cfg.ServerConfig.MaxConcurrentStreams),
		grpc.KeepaliveEnforcementPolicy(kaep),
		grpc.KeepaliveParams(kasp),
	}

	if len(interceptors) > 0 {
		serverOpts = append(serverOpts, grpc.ChainUnaryInterceptor(interceptors...))
	}

	s := grpc.NewServer(serverOpts...)

	if opts.EnableReflection {
		reflection.Register(s)
	}

	return s, nil
}

// DefaultServerOptions создает опции по умолчанию
func DefaultServerOptions(cfg *config.Config, logger *slog.Logger) ServerOptions {
	return ServerOptions{
		Config: cfg,
		Logger: logger,
		UnaryInterceptors: []grpc.UnaryServerInterceptor{
			middleware.AuthInterceptor,
			middleware.LoggerInterceptor(logger),
		},
		EnableReflection: true,
		EnableValidation: true,
	}
}

// Run запускает gRPC сервер с graceful shutdown
func Run(ctx context.Context, s *grpc.Server, cfg *config.Config, logger *slog.Logger) error {
	addr := net.JoinHostPort(cfg.ServerConfig.GRPCHost, cfg.ServerConfig.GRPCPort)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	errChan := make(chan error, 1)
	go func() {
		logger.Info("gRPC server started",
			"address", addr,
		)
		if err = s.Serve(lis); err != nil {
			errChan <- err
		}
	}()

	select {
	case err = <-errChan:
		return err
	case sig := <-quit:
		logger.InfoContext(ctx, "Shutting down server...", "signal", sig.String())
		s.GracefulStop()
		return nil
	case <-ctx.Done():
		logger.InfoContext(ctx, "Shutting down server...", "reason", "context cancelled")
		s.GracefulStop()
		return ctx.Err()
	}
}
