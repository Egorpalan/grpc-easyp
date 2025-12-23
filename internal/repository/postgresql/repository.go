package postgresql

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Egorpalan/grpc-easyp/internal/config"
	"github.com/Egorpalan/grpc-easyp/internal/repository/postgresql/notes"
)

type Repository interface {
	NewNotesQuery(ctx context.Context) notes.Querier
}

type repository struct {
	pool *pgxpool.Pool
}

func NewPGConnection(ctx context.Context, cfg *config.Config) (*pgxpool.Pool, error) {
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

	return pool, nil
}

func NewRepository(pool *pgxpool.Pool) Repository {
	return &repository{pool: pool}
}

func (r *repository) NewNotesQuery(ctx context.Context) notes.Querier {
	return notes.NewQuery(ctx, r.pool)
}
