package auth_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/Hajime3778/go-clean-arch/domain"
	"github.com/Hajime3778/go-clean-arch/infrastructure/env"
	"github.com/Hajime3778/go-clean-arch/interface/database/user/mock"
	usecase "github.com/Hajime3778/go-clean-arch/usecase/auth"
	jwt "github.com/form3tech-oss/jwt-go"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestMain(m *testing.M) {
	env.NewEnv().LoadEnvFile("../../.env")
	exitVal := m.Run()
	os.Exit(exitVal)
}

func TestSignUp(t *testing.T) {
	t.Run("正常系 新規登録", func(t *testing.T) {
		createdUserID := int64(1)
		mockUserRepo := &mock.MockUserRepo{
			MockCreate: func(ctx context.Context, user domain.User) (int64, error) {
				return createdUserID, nil
			},
		}
		mockUser := domain.User{
			Name:      "test user",
			Email:     generateRandomEmail(),
			Password:  "test password",
			CreatedAt: time.Time{},
			UpdatedAt: time.Time{},
		}
		authUsecase := usecase.NewAuthUsecase(mockUserRepo)
		tokenString, err := authUsecase.SignUp(context.TODO(), mockUser)
		token, _ := jwt.ParseWithClaims(tokenString, &usecase.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("SECRET_KEY")), nil
		})

		assert.NoError(t, err)
		claims := token.Claims.(*usecase.CustomClaims)
		assert.Equal(t, createdUserID, claims.UserID)
		assert.Equal(t, mockUser.Name, claims.UserName)
	})

	t.Run("異常系 Repository実行時にエラーが発生した場合、エラーとなること", func(t *testing.T) {
		mockUserRepo := &mock.MockUserRepo{
			MockCreate: func(ctx context.Context, task domain.User) (int64, error) {
				return 0, domain.ErrInternalServerError
			},
		}
		authUsecase := usecase.NewAuthUsecase(mockUserRepo)
		token, err := authUsecase.SignUp(context.TODO(), domain.User{})

		assert.Equal(t, domain.ErrInternalServerError, err)
		assert.Empty(t, token)
	})
}

func TestSignIn(t *testing.T) {
	t.Run("正常系 サインイン", func(t *testing.T) {
		password := "test password"
		salt := "salt"
		hashed, _ := bcrypt.GenerateFromPassword([]byte(password+salt), bcrypt.DefaultCost)
		mockUser := domain.User{
			ID:        1,
			Name:      "test user",
			Email:     generateRandomEmail(),
			Password:  string(hashed),
			Salt:      salt,
			CreatedAt: time.Time{},
			UpdatedAt: time.Time{},
		}
		mockUserRepo := &mock.MockUserRepo{
			MockGetByEmail: func(ctx context.Context, email string) (domain.User, error) {
				return mockUser, nil
			},
		}

		authUsecase := usecase.NewAuthUsecase(mockUserRepo)
		tokenString, err := authUsecase.SignIn(context.TODO(), mockUser.Email, password)
		assert.NoError(t, err)
		assert.NotEmpty(t, tokenString)
		token, _ := jwt.ParseWithClaims(tokenString, &usecase.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("SECRET_KEY")), nil
		})
		claims := token.Claims.(*usecase.CustomClaims)
		assert.Equal(t, mockUser.ID, claims.UserID)
		assert.Equal(t, mockUser.Name, claims.UserName)
	})
	t.Run("準正常系 パスワードが間違っている場合、ErrMismatchedPasswordエラーとなること", func(t *testing.T) {
		password := "test password"
		salt := "salt"
		hashed, _ := bcrypt.GenerateFromPassword([]byte(password+salt), bcrypt.DefaultCost)
		mockUser := domain.User{
			ID:        1,
			Name:      "test user",
			Email:     generateRandomEmail(),
			Password:  string(hashed),
			Salt:      salt,
			CreatedAt: time.Time{},
			UpdatedAt: time.Time{},
		}
		mockUserRepo := &mock.MockUserRepo{
			MockGetByEmail: func(ctx context.Context, email string) (domain.User, error) {
				return mockUser, nil
			},
		}
		authUsecase := usecase.NewAuthUsecase(mockUserRepo)
		tokenString, err := authUsecase.SignIn(context.TODO(), mockUser.Email, "foo bar")
		assert.Equal(t, domain.ErrMismatchedPassword, err)
		assert.Empty(t, tokenString)
	})
	t.Run("準正常系 Repositoryで取得したパスワードがハッシュ文字列でない場合エラーとなること", func(t *testing.T) {
		mockUser := domain.User{
			ID:        1,
			Name:      "test user",
			Email:     generateRandomEmail(),
			Password:  "test password",
			Salt:      "salt",
			CreatedAt: time.Time{},
			UpdatedAt: time.Time{},
		}
		mockUserRepo := &mock.MockUserRepo{
			MockGetByEmail: func(ctx context.Context, email string) (domain.User, error) {
				return mockUser, nil
			},
		}
		authUsecase := usecase.NewAuthUsecase(mockUserRepo)
		tokenString, err := authUsecase.SignIn(context.TODO(), mockUser.Email, mockUser.Password)
		assert.NotEmpty(t, err)
		assert.Empty(t, tokenString)
	})
	t.Run("異常系 Repository実行時にエラーが発生した場合、エラーとなること", func(t *testing.T) {
		mockUserRepo := &mock.MockUserRepo{
			MockGetByEmail: func(ctx context.Context, email string) (domain.User, error) {
				return domain.User{}, domain.ErrInternalServerError
			},
		}
		authUsecase := usecase.NewAuthUsecase(mockUserRepo)
		token, err := authUsecase.SignIn(context.TODO(), "", "")

		assert.Equal(t, domain.ErrInternalServerError, err)
		assert.Empty(t, token)
	})
}

func generateRandomEmail() string {
	return fmt.Sprintf("%d@example.com", time.Now().UnixNano())
}
