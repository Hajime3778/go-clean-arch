package auth_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Hajime3778/go-clean-arch/domain"
	"github.com/Hajime3778/go-clean-arch/interface/database/user/mock"
	usecase "github.com/Hajime3778/go-clean-arch/usecase/auth"
	"github.com/stretchr/testify/assert"
)

func TestSignUpCreate(t *testing.T) {
	mockUser := domain.User{
		Name:      "test user",
		Email:     "test email",
		Password:  generateRandomEmail(),
		Salt:      "salt",
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}

	mockUserRepo := &mock.MockUserRepo{
		MockCreate: func(context.Context, domain.User) (int64, error) {
			return int64(1), nil
		},
	}
	userUsecase := usecase.NewAuthUsecase(mockUserRepo)
	token, err := userUsecase.SignUp(context.TODO(), mockUser)
	if err != nil {
		t.Fatal(err)
	}
	assert.NotEmpty(t, token)
}

func generateRandomEmail() string {
	return fmt.Sprintf("%d@example.com", time.Now().UnixNano())
}
