package mock

import (
	"context"

	"github.com/Hajime3778/go-clean-arch/domain"
	usecase "github.com/Hajime3778/go-clean-arch/usecase/auth"
)

type MockAuthUsecase struct {
	usecase.AuthUsecase
	MockSignUp            func(ctx context.Context, user domain.User) (token string, err error)
	MockSignIn            func(ctx context.Context, email string, password string) (token string, err error)
	MockVerifyAccessToken func(ctx context.Context, token string) (bool, error)
}

func (m *MockAuthUsecase) SignUp(ctx context.Context, user domain.User) (token string, err error) {
	return m.MockSignUp(ctx, user)
}

func (m *MockAuthUsecase) SignIn(ctx context.Context, email string, password string) (token string, err error) {
	return m.MockSignIn(ctx, email, password)
}

func (m *MockAuthUsecase) VerifyAccessToken(ctx context.Context, token string) (bool, error) {
	return m.MockVerifyAccessToken(ctx, token)
}
