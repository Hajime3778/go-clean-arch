package database

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"os"

	"github.com/Hajime3778/go-clean-arch/interface/database"

	_ "github.com/go-sql-driver/mysql"
)

type SqlDriver struct {
	Conn *sql.DB
}

type Rows struct {
	Rows *sql.Rows
}

// NewSqlConnenction: データベースへ接続します
func NewSqlConnenction() database.SqlDriver {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")

	connStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbName)
	val := url.Values{}
	val.Add("parseTime", "1")
	dsn := fmt.Sprintf("%s?%s", connStr, val.Encode())

	conn, err := sql.Open(`mysql`, dsn)
	if err != nil {
		panic(err.Error())
	}
	return &SqlDriver{conn}
}

// Query: 取得のクエリを実行します
func (driver *SqlDriver) QueryContext(ctx context.Context, query string, args ...interface{}) (database.Rows, error) {
	rows, err := driver.Conn.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	return Rows{rows}, nil
}

// Execute: クエリを実行します
func (driver *SqlDriver) ExecuteContext(ctx context.Context, query string, args ...interface{}) (database.Result, error) {
	res := SqlResult{}
	stmt, err := driver.Conn.PrepareContext(ctx, query)
	if err != nil {
		return res, err
	}
	result, err := stmt.ExecContext(ctx, args...)
	if err != nil {
		return res, err
	}
	res.Result = result
	return res, nil
}

// ErrNoRows: データが見つからなかったときのエラー
func (driver *SqlDriver) ErrNoRows() error {
	return sql.ErrNoRows
}

// Scan: マッピングを行います
func (r Rows) Scan(dest ...interface{}) error {
	return r.Rows.Scan(dest...)
}

func (r Rows) Next() bool {
	return r.Rows.Next()
}

func (r Rows) Close() error {
	return r.Rows.Close()
}

type SqlResult struct {
	Result sql.Result
}

// LastInsertId: 追加された際のIDを返却します
func (r SqlResult) LastInsertId() (int64, error) {
	return r.Result.LastInsertId()
}

// RowsAffected: 影響のあった行数を返却します
func (r SqlResult) RowsAffected() (int64, error) {
	return r.Result.RowsAffected()
}
