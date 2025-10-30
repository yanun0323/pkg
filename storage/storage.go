package storage

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/gob"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"github.com/yanun0323/errors"
)

type storage[T any] struct {
	path string
	db   *sql.DB
	tx   *sql.Tx
}

func (l *storage[T]) driver() db {
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

	return &storage[T]{
		path: path,
		db:   db,
	}, nil
}

// Delete deletes the storage file completely
func Delete(path string) error {
	return os.Remove(path)
}

func (l *storage[T]) Exists(ctx context.Context, key string) (bool, error) {
	var count int
	err := l.driver().QueryRowContext(ctx, "SELECT COUNT(*) FROM storage WHERE key = ?", key).Scan(&count)
	if err != nil {
		return false, wrapError("exists, err: %+v", err)
	}

	return count != 0, nil
}

func (l *storage[T]) Set(ctx context.Context, key string, value T) error {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(value); err != nil {
		return wrapError("encode value, err: %+v", err)
	}

	_, err := l.driver().ExecContext(ctx, "INSERT INTO storage (key, value) VALUES (?, ?)", key, buf.Bytes())
	return wrapError("set value, err: %+v", err)
}

func (l *storage[T]) Get(ctx context.Context, key string) (T, error) {
	var (
		value    T
		blobData []byte
		err      error
	)

	err = l.driver().QueryRowContext(ctx, "SELECT value FROM storage WHERE key = ?", key).Scan(&blobData)
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

func (l *storage[T]) Find(ctx context.Context, keys ...string) ([]T, error) {
	var (
		rows *sql.Rows
		err  error
	)

	if len(keys) == 0 {
		rows, err = l.driver().QueryContext(ctx, "SELECT value FROM storage")
		if err != nil {
			return nil, wrapError("find values, err: %+v", err)
		}
	} else {
		args := make([]any, 0, len(keys))
		for _, key := range keys {
			args = append(args, key)
		}

		sql := "SELECT value FROM storage WHERE key IN (" + strings.Repeat("?, ", len(keys)-1) + "?)"

		rows, err = l.driver().QueryContext(ctx, sql, args...)
		if err != nil {
			return nil, wrapError("find values, err: %+v", err)
		}
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

func (l *storage[T]) Delete(ctx context.Context, key string) error {
	_, err := l.driver().ExecContext(ctx, "DELETE FROM storage WHERE key = ?", key)
	return wrapError("delete value, err: %+v", err)
}

func (l *storage[T]) Clear(ctx context.Context) error {
	_, err := l.driver().ExecContext(ctx, "DELETE FROM storage")
	return wrapError("clear storage, err: %+v", err)
}

func (l *storage[T]) Close() error {
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

func (l *storage[T]) Atomic(ctx context.Context, fn func(Local[T]) error) error {
	if l.tx != nil {
		return fn(l)
	}

	tx, err := l.db.BeginTx(ctx, nil)
	if err != nil {
		return wrapError("begin transaction, err: %+v", err)
	}
	defer tryRollback(tx)

	if err := fn(&storage[T]{
		path: l.path,
		db:   l.db,
		tx:   tx,
	}); err != nil {
		return wrapError("atomic operation, err: %+v", err)
	}

	return tryCommit(tx)
}
