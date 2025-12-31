package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var _ QueryEngine = (*Transaction)(nil)

type Transaction struct {
	pgx.Tx
}

// QueryRow выполняет запрос, возвращающий одну строку
func (t *Transaction) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	return t.Tx.QueryRow(ctx, sql, args...)
}

// Query выполняет запрос, возвращающий несколько строк
func (t *Transaction) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	return t.Tx.Query(ctx, sql, args...)
}

// Exec выполняет команду (INSERT, UPDATE, DELETE)
func (t *Transaction) Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error) {
	return t.Tx.Exec(ctx, sql, args...)
}

// Commit фиксирует транзакцию
func (t *Transaction) Commit(ctx context.Context) error {
	return t.Tx.Commit(ctx)
}

// Rollback откатывает транзакцию
func (t *Transaction) Rollback(ctx context.Context) error {
	return t.Tx.Rollback(ctx)
}
