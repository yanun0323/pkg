package storage

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

func openConnAndCheckType[T any](path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, wrapError("create sqlite db, err: %+v", err)
	}

	db.Exec(_schemaStorageType)
	db.Exec(_schemaStorage)

	if err := checkStorageType[T](db); err != nil {
		return nil, wrapError("check storage type, err: %+v", err)
	}

	return db, nil
}

func checkStorageType[T any](db *sql.DB) error {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM storage_type").Scan(&count)
	if err != nil {
		return wrapError("check storage type, err: %+v", err)
	}

	typeName := fmt.Sprintf("%T", new(T))

	if count != 0 {
		var name string
		err := db.QueryRow("SELECT name FROM storage_type").Scan(&name)
		if err != nil {
			return wrapError("check storage type, err: %+v", err)
		}

		if strings.EqualFold(name, typeName) {
			return nil
		} else {
			return ErrTypeMismatch
		}
	}

	if _, err := db.Exec("INSERT INTO storage_type (name) VALUES (?)", typeName); err != nil {
		return wrapError("create storage type, err: %+v", err)
	}

	return nil
}

func tryCommit(tx *sql.Tx) error {
	if tx == nil {
		return nil
	}

	if err := tx.Commit(); err != nil {
		if errors.Is(err, sql.ErrTxDone) {
			return nil
		}

		return wrapError("commit transaction, err: %+v", err)
	}

	return nil
}

func tryRollback(tx *sql.Tx) {
	if tx == nil {
		return
	}

	_ = tx.Rollback()
}
