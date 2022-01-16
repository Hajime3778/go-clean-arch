package auth

import (
	"context"

	"github.com/Hajime3778/go-clean-arch/domain"
)

type AuthUsecase interface {
	SignUp(ctx context.Context, user domain.User) (token string, err error)
	SignIn(ctx context.Context, email string, password string) (token string, err error)
}
