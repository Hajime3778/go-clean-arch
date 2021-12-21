package database

import "context"

type SqlDriver interface {
	QueryContext(context.Context, string, ...interface{}) (Rows, error)
	ExecuteContext(context.Context, string, ...interface{}) (Result, error)
	ErrNoRows() error
}

type Rows interface {
	Scan(...interface{}) error
	Next() bool
	Close() error
}

type Result interface {
	LastInsertId() (int64, error)
	RowsAffected() (int64, error)
}
