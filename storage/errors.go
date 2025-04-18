package storage

import (
	"database/sql"

	"github.com/pkg/errors"
)

var (
	// ErrTypeMismatch is returned when storage type doesn't match value type
	ErrTypeMismatch = errors.New("storage type mismatch")

	// ErrDBClosed is returned when database is closed
	ErrDBClosed = errors.New("database is closed")

	// ErrNotFound is returned when key is not found
	ErrNotFound = errors.New("key not found")
)

func wrapError(format string, err error) error {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, ErrTypeMismatch):
		return ErrTypeMismatch
	case errors.Is(err, ErrDBClosed):
		return ErrDBClosed
	case errors.Is(err, ErrNotFound):
		return ErrNotFound
	case errors.Is(err, sql.ErrConnDone):
		return ErrDBClosed
	case errors.Is(err, sql.ErrTxDone):
		return nil
	case errors.Is(err, sql.ErrNoRows):
		return ErrNotFound
	default:
		return errors.Errorf(format, err)
	}
}
