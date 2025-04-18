package storage

import (
	"bytes"
	"database/sql"
	"encoding/gob"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
)

// Local is a local storage for a single type
type Local[T any] interface {
	// Exists checks if the key exists in the storage
	Exists(key string) (bool, error)

	// Set sets the value for the key
	Set(key string, value T) error

	// Get retrieves the value associated with the specified key
	//
	// Returns ErrNotFound if the key doesn't exist in the storage
	Get(key string) (T, error)

	// Find retrieves the values associated with the specified keys
	//
	// Returns empty slice if any of the keys don't exist in the storage
	Find(keys ...string) ([]T, error)

	// Delete deletes the value for the key
	Delete(key string) error

	// Clear clears the storage
	Clear() error

	// Atomic executes a function within a transaction
	//
	// Note: Nested Atomic calls will use the same transaction
	Atomic(fn func(tx Local[T]) error) error

	// Close disconnects from the storage
	//
	// Note: Close will also commit any pending transaction
	Close() error
}

type localStorage[T any] struct {
	path string
	db   *sql.DB
	tx   *sql.Tx
}

func (l *localStorage[T]) driver() db {
	if l.tx != nil {
		return l.tx
	}
	return l.db
}

// New creates a new local storage
//
// The storage is stored in a sqlite3 file at the given path
func New[T any](path string) (Local[T], error) {
	db, err := openConnAndCheckType[T](path)
	if err != nil {
		return nil, err
	}

	return &localStorage[T]{
		path: path,
		db:   db,
	}, nil
}

// Delete deletes the storage file completely
func Delete(path string) error {
	return os.Remove(path)
}

func (l *localStorage[T]) Exists(key string) (bool, error) {
	var count int
	err := l.driver().QueryRow("SELECT COUNT(*) FROM storage WHERE key = ?", key).Scan(&count)
	if err != nil {
		return false, wrapError("exists, err: %+v", err)
	}

	return count != 0, nil
}

func (l *localStorage[T]) Set(key string, value T) error {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(value); err != nil {
		return wrapError("encode value, err: %+v", err)
	}

	_, err := l.driver().Exec("INSERT INTO storage (key, value) VALUES (?, ?)", key, buf.Bytes())
	return wrapError("set value, err: %+v", err)
}

func (l *localStorage[T]) Get(key string) (T, error) {
	var (
		value    T
		blobData []byte
		err      error
	)

	err = l.driver().QueryRow("SELECT value FROM storage WHERE key = ?", key).Scan(&blobData)
	if err != nil {
		return value, wrapError("get value, err: %+v", err)
	}

	buf := bytes.NewBuffer(blobData)
	dec := gob.NewDecoder(buf)
	err = dec.Decode(&value)
	if err != nil {
		return value, wrapError("decode value, err: %+v", err)
	}

	return value, nil
}

func (l *localStorage[T]) Find(keys ...string) ([]T, error) {
	rows, err := l.driver().Query("SELECT value FROM storage WHERE key IN (?)", keys)
	if err != nil {
		return nil, wrapError("find values, err: %+v", err)
	}

	values := make([]T, 0, len(keys))
	for rows.Next() {
		var (
			value    T
			blobData []byte
		)

		if err := rows.Scan(&blobData); err != nil {
			return nil, errors.Errorf("scan value, err: %+v", err)
		}

		buf := bytes.NewBuffer(blobData)
		dec := gob.NewDecoder(buf)
		err = dec.Decode(&value)
		if err != nil {
			return nil, errors.Errorf("decode value, err: %+v", err)
		}

		values = append(values, value)
	}

	return values, nil
}

func (l *localStorage[T]) Delete(key string) error {
	_, err := l.driver().Exec("DELETE FROM storage WHERE key = ?", key)
	return wrapError("delete value, err: %+v", err)
}

func (l *localStorage[T]) Clear() error {
	_, err := l.driver().Exec("DELETE FROM storage")
	return wrapError("clear storage, err: %+v", err)
}

func (l *localStorage[T]) Close() error {
	if l.tx != nil {
		if err := tryCommit(l.tx); err != nil {
			tryRollback(l.tx)
		}
	}

	if err := l.db.Close(); err != nil {
		if errors.Is(err, sql.ErrConnDone) {
			return nil
		}

		return wrapError("close database, err: %+v", err)
	}

	return nil
}

func (l *localStorage[T]) Atomic(fn func(Local[T]) error) error {
	if l.tx != nil {
		return fn(l)
	}

	tx, err := l.db.Begin()
	if err != nil {
		return wrapError("begin transaction, err: %+v", err)
	}
	defer tryRollback(tx)

	if err := fn(&localStorage[T]{
		path: l.path,
		db:   l.db,
		tx:   tx,
	}); err != nil {
		return wrapError("atomic operation, err: %+v", err)
	}

	return tryCommit(tx)
}
