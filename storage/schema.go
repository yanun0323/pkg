package storage

const (
	_schemaStorageType = `
CREATE TABLE IF NOT EXISTS storage_type (
	name TEXT PRIMARY KEY
)
`

	_schemaStorage = `
CREATE TABLE IF NOT EXISTS storage (
	key TEXT PRIMARY KEY,
	value BLOB NOT NULL DEFAULT '',
	created_at INTEGER NOT NULL DEFAULT 0,
	updated_at INTEGER NOT NULL DEFAULT 0
)
`
)
