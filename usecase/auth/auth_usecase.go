package auth

import (
	"context"

	"github.com/Hajime3778/go-clean-arch/domain"
	repository "github.com/Hajime3778/go-clean-arch/interface/database/user"
)

type authUsecase struct {
	repo repository.UserRepository
}

// NewAuthUsecase タスク機能のUsecaseオブジェクトを作成します
func NewAuthUsecase(repo repository.UserRepository) AuthUsecase {
	return &authUsecase{repo}
}

// SignUp ユーザーのサインアップを行います
func (u *authUsecase) SignUp(ctx context.Context, user domain.User) (token string, err error) {
	panic("not implemented") // TODO: Implement
}

// SignIn ユーザーのサインインを行います
func (u *authUsecase) SignIn(ctx context.Context, email string, password string) (token string, err error) {
	panic("not implemented") // TODO: Implement
}

// VerifyAccessToken アクセストークンの検証を行います
func (u *authUsecase) VerifyAccessToken(ctx context.Context, token string) (bool, error) {
	panic("not implemented") // TODO: Implement
}
