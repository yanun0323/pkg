package storage

import (
	"context"
	"database/sql"
)

// Local represents a local storage for a single type
type Local[T any] interface {
	// Exists checks if the key exists in the storage
	Exists(ctx context.Context, key string) (bool, error)

	// Set sets the value for the key
	Set(ctx context.Context, key string, value T) error

	// Get retrieves the value associated with the specified key
	//
	// Returns ErrNotFound if the key doesn't exist in the storage
	Get(ctx context.Context, key string) (T, error)

	// Find retrieves the values associated with the specified keys
	//
	// Returns empty slice if any of the keys don't exist in the storage
	Find(ctx context.Context, keys ...string) ([]T, error)

	// Delete deletes the value for the key
	Delete(ctx context.Context, key string) error

	// Clear clears the storage
	Clear(ctx context.Context) error

	// Atomic executes a function within a transaction
	//
	// Note: Nested Atomic calls will use the same transaction
	Atomic(ctx context.Context, fn func(tx Local[T]) error) error

	// Close disconnects from the storage
	//
	// Note: Close will also commit any pending transaction
	Close() error
}

type db interface {
	// Exec(query string, args ...any) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	// Prepare(query string) (*sql.Stmt, error)
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	// Query(query string, args ...any) (*sql.Rows, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	// QueryRow(query string, args ...any) *sql.Row
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}
