package mock

import (
	"context"
	"database/sql"

	"github.com/Hajime3778/go-clean-arch/interface/database"
)

type MockSqlDriver struct {
	database.SqlDriver
	Conn               *sql.DB
	MockQueryContext   func(context.Context, string, ...interface{}) (database.Rows, error)
	MockExecuteContext func(context.Context, string, ...interface{}) (database.Result, error)
	MockErrNoRows      func() error
}

func (m *MockSqlDriver) QueryContext(ctx context.Context, query string, args ...interface{}) (database.Rows, error) {
	return m.MockQueryContext(ctx, query, args)
}

func (m *MockSqlDriver) ExecuteContext(ctx context.Context, query string, args ...interface{}) (database.Result, error) {
	return m.MockExecuteContext(ctx, query, args)
}

func (m *MockSqlDriver) ErrNoRows() error {
	return m.MockErrNoRows()
}

type MockRows struct {
	database.Rows
	MockScan  func(...interface{}) error
	MockNext  func() bool
	MockClose func() error
}

func (m *MockRows) Scan(args ...interface{}) error {
	return m.MockScan(args)
}

func (m *MockRows) Next() bool {
	return m.MockNext()
}

func (m *MockRows) Close() error {
	return m.MockClose()
}

type MockResult struct {
	database.Result
	MockLastInsertId func() (int64, error)
	MockRowsAffected func() (int64, error)
}

func (m *MockResult) LastInsertId() (int64, error) {
	return m.MockLastInsertId()
}

func (m *MockResult) RowsAffected() (int64, error) {
	return m.MockLastInsertId()
}
