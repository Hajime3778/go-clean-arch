package auth_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/Hajime3778/go-clean-arch/domain"
	"github.com/Hajime3778/go-clean-arch/infrastructure/database"
	"github.com/Hajime3778/go-clean-arch/infrastructure/env"
	interfaceDB "github.com/Hajime3778/go-clean-arch/interface/database"
	userRepository "github.com/Hajime3778/go-clean-arch/interface/database/user"
	authHandler "github.com/Hajime3778/go-clean-arch/interface/handlers/nethttp/auth"
	"github.com/form3tech-oss/jwt-go"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

const signUpURL = "http://localhost:8080/auth/sign_up"

// const signInURL = "http://localhost:8080/auth/sign_in"

var sqlDriver interfaceDB.SqlDriver

func TestMain(m *testing.M) {
	env.NewEnv().LoadEnvFile("../../.env")
	sqlDriver = database.NewSqlConnenction()
	exitVal := m.Run()
	os.Exit(exitVal)
}

func TestSignUp(t *testing.T) {
	t.Run("正常系 ユーザーの新規登録", func(t *testing.T) {
		ctx := context.TODO()
		signUpRequest := authHandler.SignUpRequest{
			Name:     "test user",
			Email:    generateRandomEmail(),
			Password: "password",
		}
		byteRequest, _ := json.Marshal(signUpRequest)
		req, _ := http.NewRequest("POST", signUpURL, bytes.NewBuffer(byteRequest))
		req.Header.Set("Content-Type", "application/json")
		client := new(http.Client)
		res, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()

		var response authHandler.SignUpResponse
		decoder := json.NewDecoder(res.Body)
		err = decoder.Decode(&response)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, http.StatusOK, res.StatusCode)

		userID := getUserIDFromToken(response.Token)
		userRepo := userRepository.NewUserRepository(sqlDriver)
		user, _ := userRepo.GetByID(ctx, userID)
		password := []byte(signUpRequest.Password + user.Salt)
		passwordErr := bcrypt.CompareHashAndPassword([]byte(user.Password), password)

		assert.NoError(t, passwordErr)
		assert.Equal(t, userID, user.ID)
		assert.Equal(t, signUpRequest.Name, user.Name)
		assert.Equal(t, signUpRequest.Email, user.Email)
	})
	t.Run("準正常系 リクエストパラメータが足りていない場合、400エラーとなること", func(t *testing.T) {})
	t.Run("準正常系 リクエスト形式が間違っている場合、400エラーとなること", func(t *testing.T) {})
}

func TestSignIn(t *testing.T) {
	t.Run("正常系 ユーザーの新規登録後、サインインできること", func(t *testing.T) {})
	t.Run("準正常系 存在しないEmailの場合、401エラーとなること", func(t *testing.T) {})
	t.Run("準正常系 パスワードが間違っている場合、401エラーとなること", func(t *testing.T) {})
	t.Run("準正常系 リクエストパラメータが足りていない場合、400エラーとなること", func(t *testing.T) {})
	t.Run("準正常系 リクエスト形式が間違っている場合、400エラーとなること", func(t *testing.T) {})
}

func getUserIDFromToken(strToken string) int64 {
	token, _ := jwt.ParseWithClaims(strToken, &domain.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET_KEY")), nil
	})
	claims := token.Claims.(*domain.Claims)
	return claims.UserID
}

func generateRandomEmail() string {
	return fmt.Sprintf("%d@example.com", time.Now().UnixNano())
}
