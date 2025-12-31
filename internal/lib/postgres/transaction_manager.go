package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type key string

const (
	txKey key = "tx"
)

// TransactionManager - менеджер транзакций
type TransactionManager struct {
	connection *Connection
}

// NewTransactionManager создает новый менеджер транзакций
func NewTransactionManager(connection *Connection) *TransactionManager {
	return &TransactionManager{
		connection: connection,
	}
}

// GetQueryEngine возвращает QueryEngine (транзакцию или connection)
func (m *TransactionManager) GetQueryEngine(ctx context.Context) QueryEngine {
	if tx, ok := ctx.Value(txKey).(*Transaction); ok {
		return tx
	}

	return m.connection
}

// runTransaction выполняет функцию в транзакции
func (m *TransactionManager) runTransaction(ctx context.Context, txOpts pgx.TxOptions, fn func(ctx context.Context) error) (err error) {
	if _, ok := ctx.Value(txKey).(*Transaction); ok {
		return fn(ctx)
	}

	tx, err := m.connection.BeginTx(ctx, txOpts)
	if err != nil {
		return fmt.Errorf("can't begin transaction: %w", err)
	}

	ctx = context.WithValue(ctx, txKey, tx)

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic recovered: %v", r)
		}

		if err == nil {
			if commitErr := tx.Commit(ctx); commitErr != nil {
				err = fmt.Errorf("commit failed: %w", commitErr)
			}
		}

		if err != nil {
			if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
				err = fmt.Errorf("rollback failed: %w", rollbackErr)
			}
		}
	}()

	err = fn(ctx)
	return err
}

// RunReadCommitted выполняет функцию в транзакции с уровнем изоляции ReadCommitted
func (m *TransactionManager) RunReadCommitted(ctx context.Context, fn func(ctx context.Context) error) error {
	return m.runTransaction(ctx, pgx.TxOptions{
		IsoLevel:   pgx.ReadCommitted,
		AccessMode: pgx.ReadWrite,
	}, fn)
}

// RunRepeatableRead выполняет функцию в транзакции с уровнем изоляции RepeatableRead
func (m *TransactionManager) RunRepeatableRead(ctx context.Context, fn func(ctx context.Context) error) error {
	return m.runTransaction(ctx, pgx.TxOptions{
		IsoLevel:   pgx.RepeatableRead,
		AccessMode: pgx.ReadWrite,
	}, fn)
}

// RunSerializable выполняет функцию в транзакции с уровнем изоляции Serializable
func (m *TransactionManager) RunSerializable(ctx context.Context, fn func(ctx context.Context) error) error {
	return m.runTransaction(ctx, pgx.TxOptions{
		IsoLevel:   pgx.Serializable,
		AccessMode: pgx.ReadWrite,
	}, fn)
}

func (m *TransactionManager) RunInTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return m.RunReadCommitted(ctx, fn)
}
