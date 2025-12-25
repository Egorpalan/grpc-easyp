package main

import (
	"context"
	"log"
	"os"

	notesv1 "github.com/Egorpalan/grpc-easyp/internal/app/notes/v1"
	"github.com/Egorpalan/grpc-easyp/internal/config"
	"github.com/Egorpalan/grpc-easyp/internal/lib/app"
	"github.com/Egorpalan/grpc-easyp/internal/logger"
	"github.com/Egorpalan/grpc-easyp/internal/repository/postgresql"
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

	dataBiPgConnection, err := postgresql.NewPGConnection(ctx, cfg)
	if err != nil {
		log.Fatalf("failed to establish databi pg connection: %v", err)
	}
	defer dataBiPgConnection.Close()

	dataBiPgDB := postgresql.NewRepository(dataBiPgConnection)
	notesService := notes.NewService(pLogger, dataBiPgDB)

	s, err := app.NewServer(app.DefaultServerOptions(cfg, pLogger))
	if err != nil {
		log.Fatalf("Failed to initialize server: %v", err)
	}

	notesv1.New(notesService).RegisterServer(s)

	if err = app.Run(ctx, s, cfg, pLogger); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
