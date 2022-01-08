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
