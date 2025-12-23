package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	notesv1 "github.com/Egorpalan/grpc-easyp/internal/app/notes/v1"
	"github.com/Egorpalan/grpc-easyp/internal/config"
	"github.com/Egorpalan/grpc-easyp/internal/logger"
	"github.com/Egorpalan/grpc-easyp/internal/repository/postgresql"
	"github.com/Egorpalan/grpc-easyp/internal/service/notes"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to initialize config: %v", err)
	}

	pLogger := logger.NewLogger(os.Stdout)

	dataBiPgConnection, err := postgresql.NewPGConnection(ctx, cfg)
	if err != nil {
		log.Fatalf("failed to establish databi pg connection: %v", err)
	}
	defer dataBiPgConnection.Close()

	dataBiPgDB := postgresql.NewRepository(dataBiPgConnection)
	notesService := notes.NewService(pLogger, dataBiPgDB)

	s := grpc.NewServer()

	notesv1.New(notesService).RegisterServer(s)

	reflection.Register(s)

	addr := net.JoinHostPort(cfg.ServerConfig.GRPCHost, cfg.ServerConfig.GRPCPort)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	go func() {
		pLogger.Info("gRPC server started", "address", addr)
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	pLogger.Info("Shutting down server...")
	s.GracefulStop()
}
