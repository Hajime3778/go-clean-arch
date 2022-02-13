package token_test

import (
	"os"
	"testing"

	"github.com/Hajime3778/go-clean-arch/domain"
	"github.com/Hajime3778/go-clean-arch/util/token"
	"github.com/form3tech-oss/jwt-go"
	"github.com/stretchr/testify/assert"
)

func TestGenerateRundomString(t *testing.T) {
	t.Run("正常系 トークンの値が正しいこと", func(t *testing.T) {
		user := domain.User{
			ID:   1,
			Name: "test name",
		}
		tokenString := token.GenerateAccessToken(user)
		token, _ := jwt.ParseWithClaims(tokenString, &domain.Claims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("SECRET_KEY")), nil
		})
		claims := token.Claims.(*domain.Claims)
		assert.Equal(t, user.ID, claims.UserID)
		assert.Equal(t, user.Name, claims.UserName)
	})
}
