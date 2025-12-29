package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var _ QueryEngine = (*Connection)(nil)

type Connection struct {
	pool *pgxpool.Pool
}

// NewConnection создает новое подключение
func NewConnection(pool *pgxpool.Pool) *Connection {
	return &Connection{pool: pool}
}

// Close закрывает пул соединений
func (c *Connection) Close() {
	c.pool.Close()
}

// QueryRow выполняет запрос, возвращающий одну строку
func (c *Connection) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	return c.pool.QueryRow(ctx, sql, args...)
}

// Query выполняет запрос, возвращающий несколько строк
func (c *Connection) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	return c.pool.Query(ctx, sql, args...)
}

// Exec выполняет команду (INSERT, UPDATE, DELETE)
func (c *Connection) Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error) {
	return c.pool.Exec(ctx, sql, args...)
}

// Begin начинает транзакцию
func (c *Connection) Begin(ctx context.Context) (*Transaction, error) {
	tx, err := c.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("postgres: begin transaction: %w", err)
	}
	return &Transaction{Tx: tx}, nil
}

// BeginTx начинает транзакцию с опциями
func (c *Connection) BeginTx(ctx context.Context, txOpts pgx.TxOptions) (*Transaction, error) {
	tx, err := c.pool.BeginTx(ctx, txOpts)
	if err != nil {
		return nil, fmt.Errorf("postgres: begin transaction: %w", err)
	}
	return &Transaction{Tx: tx}, nil
}
