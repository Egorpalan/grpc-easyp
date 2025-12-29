package main

import (
	"context"
	"log"
	"os"

	chatv1 "github.com/Egorpalan/grpc-easyp/internal/app/chat/v1"
	eventsv1 "github.com/Egorpalan/grpc-easyp/internal/app/events/v1"
	metricsv1 "github.com/Egorpalan/grpc-easyp/internal/app/metrics/v1"
	notesv1 "github.com/Egorpalan/grpc-easyp/internal/app/notes/v1"
	"github.com/Egorpalan/grpc-easyp/internal/config"
	"github.com/Egorpalan/grpc-easyp/internal/lib/app"
	"github.com/Egorpalan/grpc-easyp/internal/logger"
	"github.com/Egorpalan/grpc-easyp/internal/repository/postgresql"
	"github.com/Egorpalan/grpc-easyp/internal/service/chat"
	"github.com/Egorpalan/grpc-easyp/internal/service/events"
	"github.com/Egorpalan/grpc-easyp/internal/service/metrics"
	"github.com/Egorpalan/grpc-easyp/internal/service/notes"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to initialize config: %v", err)
	}

	pLogger := logger.NewLogger(os.Stdout)

	pgConnection, err := postgresql.NewPGConnection(ctx, cfg)
	if err != nil {
		log.Fatalf("failed to establish databi pg connection: %v", err)
	}
	defer pgConnection.Close()

	NotePgDB := postgresql.NewRepository(pgConnection)

	eventsService := events.NewService()
	metricsService := metrics.NewService(pLogger)
	notesService := notes.NewService(pLogger, NotePgDB, eventsService)
	chatService := chat.NewService(pLogger)

	s, err := app.NewServer(app.DefaultServerOptions(cfg, pLogger))
	if err != nil {
		log.Fatalf("Failed to initialize server: %v", err)
	}

	notesv1.New(notesService).RegisterServer(s)
	eventsv1.New(eventsService).RegisterServer(s)
	metricsv1.New(metricsService).RegisterServer(s)
	chatv1.New(chatService).RegisterServer(s)

	if err = app.Run(ctx, s, cfg, pLogger); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
