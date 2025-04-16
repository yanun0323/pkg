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

	// Get gets the value for the key
	Get(key string) (T, bool, error)

	// Delete deletes the value for the key
	Delete(key string) error

	// Clear clears the storage
	Clear() error

	// Atomic executes a function in a transaction
	Atomic(fn func(tx Local[T]) error) error

	// Close closes the storage connection
	Close() error
}

type localStorage[T any] struct {
	path string
	db   db
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
	err := l.db.QueryRow("SELECT COUNT(*) FROM storage WHERE key = ?", key).Scan(&count)
	if err != nil {
		return false, err
	}

	return count != 0, nil
}

func (l *localStorage[T]) Set(key string, value T) error {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(value); err != nil {
		return errors.Errorf("encode value, err: %+v", err)
	}

	_, err := l.db.Exec("INSERT INTO storage (key, value) VALUES (?, ?)", key, buf.Bytes())
	return err
}

func (l *localStorage[T]) Get(key string) (T, bool, error) {
	var value T
	var blobData []byte
	err := l.db.QueryRow("SELECT value FROM storage WHERE key = ?", key).Scan(&blobData)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return value, false, nil
		}

		return value, false, err
	}

	buf := bytes.NewBuffer(blobData)
	dec := gob.NewDecoder(buf)
	err = dec.Decode(&value)
	if err != nil {
		return value, false, errors.Errorf("decode value, err: %+v", err)
	}

	return value, true, nil
}

func (l *localStorage[T]) Delete(key string) error {
	_, err := l.db.Exec("DELETE FROM storage WHERE key = ?", key)
	return err
}

func (l *localStorage[T]) Clear() error {
	_, err := l.db.Exec("DELETE FROM storage")
	return err
}

func (l *localStorage[T]) Close() error {
	return l.db.Close()
}

func (l *localStorage[T]) Atomic(fn func(Local[T]) error) error {
	tx, err := l.db.Begin()
	if err != nil {
		return errors.Errorf("begin transaction, err: %+v", err)
	}

	if err := fn(l); err != nil {
		return errors.Errorf("atomic operation, err: %+v", err)
	}

	return tx.Commit()
}
