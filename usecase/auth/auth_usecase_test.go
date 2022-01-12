package auth_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Hajime3778/go-clean-arch/domain"
	"github.com/Hajime3778/go-clean-arch/infrastructure/database"
	"github.com/Hajime3778/go-clean-arch/infrastructure/env"
	userRepository "github.com/Hajime3778/go-clean-arch/interface/database/user"
	usecase "github.com/Hajime3778/go-clean-arch/usecase/auth"
	"github.com/stretchr/testify/assert"
)

func TestSignUpCreate(t *testing.T) {
	env.NewEnv().LoadEnvFile("../../.env")
	sqlDriver := database.NewSqlConnenction()
	userRepo := userRepository.NewUserRepository(sqlDriver)

	mockUser := domain.User{
		Name:      "test user",
		Email:     generateRandomEmail(),
		Password:  "test password",
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}

	// mockUserRepo := &mock.MockUserRepo{
	// 	MockCreate: func(context.Context, domain.User) (int64, error) {
	// 		return int64(1), nil
	// 	},
	// }
	ctx := context.TODO()

	userUsecase := usecase.NewAuthUsecase(userRepo)
	token, err := userUsecase.SignUp(ctx, mockUser)
	if err != nil {
		t.Fatal(err)
	}
	assert.NotEmpty(t, token)

	token, err = userUsecase.SignIn(ctx, mockUser.Email, mockUser.Password)
	if err != nil {
		t.Fatal(err)
	}
	assert.NotEmpty(t, token)
}

func generateRandomEmail() string {
	return fmt.Sprintf("%d@example.com", time.Now().UnixNano())
}
