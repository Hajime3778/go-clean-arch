package database

type SqlDriver interface {
	Query(string, ...interface{}) (Rows, error)
	QueryRow(string, ...interface{}) Row
	Execute(string, ...interface{}) (Result, error)
	ErrNoRows() error
}

type Row interface {
	Scan(...interface{}) error
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
