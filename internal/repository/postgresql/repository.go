package postgresql

import (
	"context"
	"errors"
	"fmt"

	"github.com/Egorpalan/grpc-easyp/internal/config"
	"github.com/Egorpalan/grpc-easyp/internal/lib/postgres"
	"github.com/Egorpalan/grpc-easyp/internal/repository/postgresql/notes"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	NewNotesQuery(ctx context.Context) notes.Querier
	RunInTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}

type repository struct {
	txManager *postgres.TransactionManager
}

func NewPGConnection(ctx context.Context, cfg *config.Config) (*postgres.Connection, error) {
	dbConfig, err := pgxpool.ParseConfig(fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.PostgresConfig.PostgresUser,
		cfg.PostgresConfig.PostgresPassword,
		cfg.PostgresConfig.PostgresHost,
		cfg.PostgresConfig.PostgresPort,
		cfg.PostgresConfig.PostgresDB,
		cfg.PostgresConfig.PostgresSslMode))
	if err != nil {
		return nil, errors.New("failed to parse pg dsn config string")
	}

	dbConfig.MaxConns = cfg.PostgresConfig.PostgresMaxConn

	pool, err := pgxpool.NewWithConfig(ctx, dbConfig)
	if err != nil {
		return nil, err
	}

	err = pool.Ping(ctx)
	if err != nil {
		return nil, err
	}

	return postgres.NewConnection(pool), nil
}

func NewRepository(connection *postgres.Connection) Repository {
	txManager := postgres.NewTransactionManager(connection)
	return &repository{
		txManager: txManager,
	}
}

func (r *repository) NewNotesQuery(ctx context.Context) notes.Querier {
	queryEngine := r.txManager.GetQueryEngine(ctx)
	return notes.NewQuery(ctx, queryEngine)
}

func (r *repository) RunInTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return r.txManager.RunReadCommitted(ctx, fn)
}
