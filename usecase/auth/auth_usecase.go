package auth

import (
	"context"
	"crypto/rand"
	"os"
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

type CustomClaims struct {
	UserID   int64  `json:"user_id"`
	UserName string `json:"name"`
	jwt.StandardClaims
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

// GenerateAccessToken アクセストークンを発行します
func GenerateAccessToken(user domain.User) string {
	claims := CustomClaims{
		user.ID,
		user.Name,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	return tokenString
}

// generateSalt Saltを作成します(10桁のランダム文字列)
func generateSalt() string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	b := make([]byte, 10)
	rand.Read(b)

	var result string
	for _, v := range b {
		result += string(letters[int(v)%len(letters)])
	}
	return result
}
