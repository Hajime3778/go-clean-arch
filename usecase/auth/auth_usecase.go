package auth

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/Hajime3778/go-clean-arch/domain"
	repository "github.com/Hajime3778/go-clean-arch/interface/database/user"
	jwt "github.com/form3tech-oss/jwt-go"
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
func (u *authUsecase) SignUp(ctx context.Context, user domain.User) (token string, err error) {
	// bcryptはsaltを内包しているので、saltを付与する必要はないのですが
	// salt機能がないライブラリも多いので、自身の練習&参考用サンプルとしてsaltをつけてます。
	// https://github.com/golang/crypto/blob/e495a2d5b3d3be43468d0ebb413f46eeaedf7eb3/bcrypt/bcrypt.go#L144
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

	token = GenerateAccessToken(user)

	return token, nil
}

// SignIn ユーザーのサインインを行います
func (u *authUsecase) SignIn(ctx context.Context, email string, password string) (string, error) {
	user, err := u.repo.GetByEmail(ctx, email)
	if err != nil {
		return "", err
	}
	inputPassword := []byte(password + user.Salt)
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), inputPassword)
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return "", domain.ErrMismatchedPassword
	}
	if err != nil {
		return "", err
	}
	token := GenerateAccessToken(user)
	return token, err
}

// VerifyAccessToken アクセストークンの検証を行います
func (u *authUsecase) VerifyAccessToken(ctx context.Context, token string) (bool, error) {
	panic("not implemented") // TODO: Implement
}

func GenerateAccessToken(user domain.User) string {
	// headerのセット
	token := jwt.New(jwt.SigningMethodHS256)

	// claimsのセット
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = fmt.Sprint(user.ID)
	claims["name"] = user.Name
	claims["expires_at"] = time.Now().Add(time.Hour * 24).Unix()

	// 電子署名
	tokenString, _ := token.SignedString([]byte("TODO: secret-key"))
	return tokenString
}

// generateSalt Saltを作成します(10桁のランダム文字列)
func generateSalt() string {
	var letter = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	b := make([]rune, 10)
	for i := range b {
		b[i] = letter[rand.Intn(len(letter))]
	}
	return string(b)
}

// type jwtCustomClaims struct {
// 	UserID int64  `json:"user_id"`
// 	Name   string `json:"name"`
// 	jwt.StandardClaims
// }
