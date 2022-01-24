package token

import (
	"os"
	"time"

	"github.com/Hajime3778/go-clean-arch/domain"
	"github.com/form3tech-oss/jwt-go"
)

// GenerateAccessToken アクセストークンを発行します
func GenerateAccessToken(user domain.User) string {
	claims := domain.Claims{
		UserID:   user.ID,
		UserName: user.Name,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	return tokenString
}
