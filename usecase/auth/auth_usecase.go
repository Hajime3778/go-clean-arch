package auth

import (
	"context"

	"github.com/Hajime3778/go-clean-arch/domain"
	repository "github.com/Hajime3778/go-clean-arch/interface/database/user"
	"github.com/Hajime3778/go-clean-arch/util/string_util"
	"github.com/Hajime3778/go-clean-arch/util/token"
	"golang.org/x/crypto/bcrypt"
)

type authUsecase struct {
	repo repository.UserRepository
}

// NewAuthUsecase タスク機能のUsecaseオブジェクトを作成します
func NewAuthUsecase(repo repository.UserRepository) AuthUsecase {
	return &authUsecase{repo}
}

// SignUp ユーザーのサインアップを行います
// bcryptはsaltを内包しているので、saltを付与する必要はないのですが
// salt機能がないライブラリも多いので、自身の練習&参考用サンプルとしてsaltをつけてます。
// https://github.com/golang/crypto/blob/e495a2d5b3d3be43468d0ebb413f46eeaedf7eb3/bcrypt/bcrypt.go#L144
func (u *authUsecase) SignUp(ctx context.Context, user domain.User) (string, error) {
	_, err := u.repo.GetByEmail(ctx, user.Email)
	if err == nil {
		return "", domain.ErrExistEmail
	}
	if err != nil && err != domain.ErrRecordNotFound {
		return "", err
	}

	salt := generateSalt()
	password := []byte(user.Password + salt)
	hashed, _ := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)

	user.Password = string(hashed)
	user.Salt = salt

	userID, err := u.repo.Create(ctx, user)
	if err != nil {
		return "", err
	}
	user.ID = userID

	token := token.GenerateAccessToken(user)

	return token, nil
}

// SignIn ユーザーのサインインを行います
func (u *authUsecase) SignIn(ctx context.Context, email string, password string) (string, error) {
	user, err := u.repo.GetByEmail(ctx, email)
	if err == domain.ErrRecordNotFound {
		return "", domain.ErrFailedSignIn
	}
	if err != nil {
		return "", err
	}
	inputPassword := []byte(password + user.Salt)
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), inputPassword)
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return "", domain.ErrFailedSignIn
	}
	if err != nil {
		return "", err
	}
	token := token.GenerateAccessToken(user)
	return token, err
}

// generateSalt Saltを作成します(10桁のランダム文字列)
func generateSalt() string {
	return string_util.GenerateRundomString(10)
}
