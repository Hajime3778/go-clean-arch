package user_test

import (
	"context"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Hajime3778/go-clean-arch/domain"
	infrastructure "github.com/Hajime3778/go-clean-arch/infrastructure/database"
	"github.com/Hajime3778/go-clean-arch/interface/database"
	mockSqlDriver "github.com/Hajime3778/go-clean-arch/interface/database/mock"
	userRepository "github.com/Hajime3778/go-clean-arch/interface/database/user"
	"github.com/stretchr/testify/assert"
)

func TestGetByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock error: '%s'", err)
	}

	sqlDriver := new(infrastructure.SqlDriver)
	sqlDriver.Conn = db

	mockUser := domain.User{
		Name:     "test name",
		Email:    "test email",
		Password: "test password",
		Salt:     "test salt",
	}

	repo := userRepository.NewUserRepository(sqlDriver)
	query := "SELECT * FROM users WHERE id = ?"

	t.Run("正常系 存在するIDで1件取得", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name", "email", "password", "salt", "updated_at", "created_at"}).
			AddRow(mockUser.ID, mockUser.Name, mockUser.Email, mockUser.Password, mockUser.Salt, mockUser.UpdatedAt, mockUser.CreatedAt)
		mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(mockUser.ID).WillReturnRows(rows)

		got, err := repo.GetByID(context.TODO(), mockUser.ID)
		assert.NoError(t, err)
		assert.Equal(t, mockUser, got)
	})

	t.Run("準正常系 存在しないIDで検索してエラーとなること", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name", "email", "password", "salt", "updated_at", "created_at"})
		mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(int64(2)).WillReturnRows(rows)

		got, err := repo.GetByID(context.TODO(), int64(2))
		assert.Equal(t, domain.ErrRecordNotFound, err)
		assert.Equal(t, domain.User{}, got)
	})

	t.Run("異常系 クエリ実行で失敗した場合エラーが返却されること", func(t *testing.T) {
		mockErr := errors.New("query failed error")
		mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(int64(2)).WillReturnError(mockErr)

		got, err := repo.GetByID(context.TODO(), int64(2))
		assert.Equal(t, mockErr, err)
		assert.Equal(t, domain.User{}, got)
	})

	t.Run("異常系 Scan実行で失敗した場合エラーが返却されること", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"foo"}).AddRow("bar")
		mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(int64(2)).WillReturnRows(rows)

		got, err := repo.GetByID(context.TODO(), int64(2))
		assert.NotNil(t, err)
		assert.Equal(t, domain.User{}, got)
	})

	t.Run("異常系 Rows.Close実行で失敗した場合にログが出力されること", func(t *testing.T) {
		mockErr := errors.New("rows close error")
		rows := sqlmock.NewRows([]string{"foo"}).AddRow("bar").CloseError(mockErr)
		mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(int64(2)).WillReturnRows(rows)

		got, err := repo.GetByID(context.TODO(), int64(2))
		assert.NotNil(t, err)
		assert.Equal(t, domain.User{}, got)
	})
}

func TestGetByEmail(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock error: '%s'", err)
	}

	sqlDriver := new(infrastructure.SqlDriver)
	sqlDriver.Conn = db

	mockUser := domain.User{
		Name:     "test name",
		Email:    "test email",
		Password: "test password",
		Salt:     "test salt",
	}

	repo := userRepository.NewUserRepository(sqlDriver)
	query := "SELECT * FROM users WHERE email = ?"

	t.Run("正常系 存在するEmailで1件取得", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name", "email", "password", "salt", "updated_at", "created_at"}).
			AddRow(mockUser.ID, mockUser.Name, mockUser.Email, mockUser.Password, mockUser.Salt, mockUser.UpdatedAt, mockUser.CreatedAt)
		mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(mockUser.Email).WillReturnRows(rows)

		got, err := repo.GetByEmail(context.TODO(), mockUser.Email)
		assert.NoError(t, err)
		assert.Equal(t, mockUser, got)
	})

	t.Run("準正常系 存在しないEmailで検索してエラーとなること", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name", "email", "password", "salt", "updated_at", "created_at"})
		mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(mockUser.Email).WillReturnRows(rows)

		got, err := repo.GetByEmail(context.TODO(), mockUser.Email)
		assert.Equal(t, domain.ErrRecordNotFound, err)
		assert.Equal(t, domain.User{}, got)
	})

	t.Run("異常系 クエリ実行で失敗した場合エラーが返却されること", func(t *testing.T) {
		mockErr := errors.New("query failed error")
		mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(mockUser.Email).WillReturnError(mockErr)

		got, err := repo.GetByEmail(context.TODO(), mockUser.Email)
		assert.Equal(t, mockErr, err)
		assert.Equal(t, domain.User{}, got)
	})

	t.Run("異常系 Scan実行で失敗した場合エラーが返却されること", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"foo"}).AddRow("bar")
		mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(mockUser.Email).WillReturnRows(rows)

		got, err := repo.GetByEmail(context.TODO(), mockUser.Email)
		assert.NotNil(t, err)
		assert.Equal(t, domain.User{}, got)
	})

	t.Run("異常系 Rows.Close実行で失敗した場合にログが出力されること", func(t *testing.T) {
		mockErr := errors.New("rows close error")
		rows := sqlmock.NewRows([]string{"foo"}).AddRow("bar").CloseError(mockErr)
		mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(mockUser.Email).WillReturnRows(rows)

		got, err := repo.GetByEmail(context.TODO(), mockUser.Email)
		assert.NotNil(t, err)
		assert.Equal(t, domain.User{}, got)
	})
}

func TestCreate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock error: '%s'", err)
	}

	t.Run("正常系 1件追加", func(t *testing.T) {
		sqlDriver := new(infrastructure.SqlDriver)
		sqlDriver.Conn = db

		repo := userRepository.NewUserRepository(sqlDriver)
		query := "INSERT INTO users(name,email,password,salt) VALUES(?,?,?,?)"

		mockUser := domain.User{
			Name:     "test user",
			Email:    "test@example.com",
			Password: "test passsword",
			Salt:     "test salt",
		}
		prep := mock.ExpectPrepare(regexp.QuoteMeta(query))
		prep.ExpectExec().
			WithArgs(mockUser.Name, mockUser.Email, mockUser.Password, mockUser.Salt).
			WillReturnResult(sqlmock.NewResult(12, 1))

		id, err := repo.Create(context.TODO(), mockUser)
		assert.NoError(t, err)
		assert.NotEqual(t, int64(0), id)
	})

	t.Run("異常系 クエリ実行で失敗した場合エラーが返却されること", func(t *testing.T) {
		sqlDriver := new(infrastructure.SqlDriver)
		sqlDriver.Conn = db

		repo := userRepository.NewUserRepository(sqlDriver)
		query := "INSERT INTO users(name,email,password,salt) VALUES(?,?,?,?)"

		mockUser := domain.User{
			Name:     "test user",
			Email:    "test@example.com",
			Password: "test passsword",
			Salt:     "test salt",
		}
		mockErr := errors.New("query failed error")
		prep := mock.ExpectPrepare(regexp.QuoteMeta(query))
		prep.ExpectExec().
			WithArgs(mockUser.Name, mockUser.Email, mockUser.Password, mockUser.Salt).
			WillReturnError(mockErr)

		id, err := repo.Create(context.TODO(), mockUser)
		assert.Equal(t, mockErr, err)
		assert.Equal(t, int64(0), id)
	})

	t.Run("異常系 追加後IDで失敗した場合エラーが返却されること", func(t *testing.T) {
		mockErr := errors.New("test error")
		mockResult := &mockSqlDriver.MockResult{
			MockLastInsertId: func() (int64, error) {
				return 0, mockErr
			},
		}
		mockDriver := &mockSqlDriver.MockSqlDriver{
			MockExecuteContext: func(context.Context, string, ...interface{}) (database.Result, error) {
				return mockResult, nil
			},
		}
		mockDriver.Conn = db

		repo := userRepository.NewUserRepository(mockDriver)
		query := "INSERT INTO users(name,email,password,salt) VALUES(?,?,?,?)"

		mockUser := domain.User{
			Name:     "test user",
			Email:    "test@example.com",
			Password: "test passsword",
			Salt:     "test salt",
		}

		prep := mock.ExpectPrepare(regexp.QuoteMeta(query))
		prep.ExpectExec().
			WithArgs(mockUser.Name, mockUser.Email, mockUser.Password, mockUser.Salt).
			WillReturnResult(sqlmock.NewResult(12, 1))

		id, err := repo.Create(context.TODO(), mockUser)
		assert.Equal(t, mockErr, err)
		assert.Equal(t, int64(0), id)
	})
}
