package auth_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Hajime3778/go-clean-arch/domain"
	"github.com/Hajime3778/go-clean-arch/interface/handlers/nethttp/auth"
	"github.com/Hajime3778/go-clean-arch/usecase/auth/mock"
	"github.com/stretchr/testify/assert"
)

func TestSignUp(t *testing.T) {
	t.Run("正常系 新規登録成功", func(t *testing.T) {
		req := auth.SignUpRequest{
			Name:     "test name",
			Email:    "test@example.com",
			Password: "test password",
		}
		byteReq, _ := json.Marshal(req)
		r := httptest.NewRequest(http.MethodPost, "http://example.com/auth/sign_up",
			bytes.NewBuffer(byteReq),
		)
		w := httptest.NewRecorder()
		mockUsecase := &mock.MockAuthUsecase{
			MockSignUp: func(ctx context.Context, task domain.User) (string, error) {
				return "mock token", nil
			},
		}
		handler := auth.NewAuthHandler(mockUsecase)
		handler.SignUpHandler(w, r)
		res := w.Result()
		defer res.Body.Close()

		var response auth.SignUpResponse
		decoder := json.NewDecoder(res.Body)
		decoder.DisallowUnknownFields()
		err := decoder.Decode(&response)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Equal(t, "mock token", response.Token)
	})

	t.Run("異常系 実装していないメソッドでリクエストした場合、404エラーとなること", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodOptions, "http://example.com/auth/sign_up", nil)
		w := httptest.NewRecorder()
		mockUsecase := &mock.MockAuthUsecase{}
		handler := auth.NewAuthHandler(mockUsecase)
		handler.SignUpHandler(w, r)
		res := w.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusNotFound, res.StatusCode)
	})

	t.Run("準正常系 パラメータが指定されていない場合、400エラーとなること", func(t *testing.T) {
		req := auth.SignUpRequest{
			Email:    "test@example.com",
			Password: "test password",
		}
		byteReq, _ := json.Marshal(req)
		r := httptest.NewRequest(http.MethodPost, "http://example.com/auth/sign_up",
			bytes.NewBuffer(byteReq),
		)
		w := httptest.NewRecorder()
		mockUsecase := &mock.MockAuthUsecase{
			MockSignUp: func(ctx context.Context, task domain.User) (string, error) {
				return "mock token", nil
			},
		}
		handler := auth.NewAuthHandler(mockUsecase)
		handler.SignUpHandler(w, r)
		res := w.Result()
		defer res.Body.Close()

		var response domain.ErrorResponse
		decoder := json.NewDecoder(res.Body)
		err := decoder.Decode(&response)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
		assert.NotEmpty(t, response.Message)
	})

	t.Run("準正常系 リクエスト形式が間違っている場合、400エラーとなること", func(t *testing.T) {
		req := domain.ErrorResponse{
			Message: "test",
		}
		byteReq, _ := json.Marshal(req)
		r := httptest.NewRequest(http.MethodPost, "http://example.com/auth/sign_up",
			bytes.NewBuffer(byteReq),
		)
		w := httptest.NewRecorder()
		mockUsecase := &mock.MockAuthUsecase{
			MockSignUp: func(ctx context.Context, task domain.User) (string, error) {
				return "mock token", nil
			},
		}
		handler := auth.NewAuthHandler(mockUsecase)
		handler.SignUpHandler(w, r)
		res := w.Result()
		defer res.Body.Close()

		var response domain.ErrorResponse
		decoder := json.NewDecoder(res.Body)
		err := decoder.Decode(&response)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
		assert.NotEmpty(t, response.Message)
	})

	t.Run("異常系 Usecase実行時にエラーが発生した場合、エラーとなること", func(t *testing.T) {
		req := auth.SignUpRequest{
			Name:     "test name",
			Email:    "test@example.com",
			Password: "test password",
		}
		byteReq, _ := json.Marshal(req)
		r := httptest.NewRequest(http.MethodPost, "http://example.com/auth/sign_up",
			bytes.NewBuffer(byteReq),
		)
		mockErr := errors.New("test error")
		w := httptest.NewRecorder()
		mockUsecase := &mock.MockAuthUsecase{
			MockSignUp: func(ctx context.Context, task domain.User) (string, error) {
				return "", mockErr
			},
		}
		handler := auth.NewAuthHandler(mockUsecase)
		handler.SignUpHandler(w, r)
		res := w.Result()
		defer res.Body.Close()

		var response domain.ErrorResponse
		decoder := json.NewDecoder(res.Body)
		err := decoder.Decode(&response)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
		assert.Equal(t, domain.ErrorResponse{Message: mockErr.Error()}, response)
	})
}

func TestSignIn(t *testing.T) {
	t.Run("正常系 サインイン成功", func(t *testing.T) {
		req := auth.SignInRequest{
			Email:    "test@example.com",
			Password: "test password",
		}
		byteReq, _ := json.Marshal(req)
		r := httptest.NewRequest(http.MethodPost, "http://example.com/auth/sign_in",
			bytes.NewBuffer(byteReq),
		)
		w := httptest.NewRecorder()
		mockUsecase := &mock.MockAuthUsecase{
			MockSignIn: func(ctx context.Context, email string, password string) (string, error) {
				return "mock token", nil
			},
		}
		handler := auth.NewAuthHandler(mockUsecase)
		handler.SignInHandler(w, r)
		res := w.Result()
		defer res.Body.Close()

		var response auth.SignInResponse
		decoder := json.NewDecoder(res.Body)
		decoder.DisallowUnknownFields()
		err := decoder.Decode(&response)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Equal(t, "mock token", response.Token)
	})

	t.Run("異常系 実装していないメソッドでリクエストした場合、404エラーとなること", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodOptions, "http://example.com/auth/sign_in", nil)
		w := httptest.NewRecorder()
		mockUsecase := &mock.MockAuthUsecase{}
		handler := auth.NewAuthHandler(mockUsecase)
		handler.SignInHandler(w, r)
		res := w.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusNotFound, res.StatusCode)
	})

	t.Run("準正常系 パラメータが指定されていない場合、400エラーとなること", func(t *testing.T) {
		req := auth.SignInRequest{
			Password: "test password",
		}
		byteReq, _ := json.Marshal(req)
		r := httptest.NewRequest(http.MethodPost, "http://example.com/auth/sign_in",
			bytes.NewBuffer(byteReq),
		)
		w := httptest.NewRecorder()
		mockUsecase := &mock.MockAuthUsecase{
			MockSignIn: func(ctx context.Context, email string, password string) (string, error) {
				return "mock token", nil
			},
		}
		handler := auth.NewAuthHandler(mockUsecase)
		handler.SignInHandler(w, r)
		res := w.Result()
		defer res.Body.Close()

		var response domain.ErrorResponse
		decoder := json.NewDecoder(res.Body)
		err := decoder.Decode(&response)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
		assert.NotEmpty(t, response.Message)
	})

	t.Run("準正常系 リクエスト形式が間違っている場合、400エラーとなること", func(t *testing.T) {
		req := domain.ErrorResponse{
			Message: "test",
		}
		byteReq, _ := json.Marshal(req)
		r := httptest.NewRequest(http.MethodPost, "http://example.com/auth/sign_in",
			bytes.NewBuffer(byteReq),
		)
		w := httptest.NewRecorder()
		mockUsecase := &mock.MockAuthUsecase{
			MockSignIn: func(ctx context.Context, email string, password string) (string, error) {
				return "mock token", nil
			},
		}
		handler := auth.NewAuthHandler(mockUsecase)
		handler.SignInHandler(w, r)
		res := w.Result()
		defer res.Body.Close()

		var response domain.ErrorResponse
		decoder := json.NewDecoder(res.Body)
		err := decoder.Decode(&response)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
		assert.NotEmpty(t, response.Message)
	})

	t.Run("異常系 Usecase実行時にエラーが発生した場合、エラーとなること", func(t *testing.T) {
		req := auth.SignInRequest{
			Email:    "test@example.com",
			Password: "test password",
		}
		byteReq, _ := json.Marshal(req)
		r := httptest.NewRequest(http.MethodPost, "http://example.com/auth/sign_in",
			bytes.NewBuffer(byteReq),
		)
		mockErr := errors.New("test error")
		w := httptest.NewRecorder()
		mockUsecase := &mock.MockAuthUsecase{
			MockSignIn: func(ctx context.Context, email string, password string) (string, error) {
				return "", mockErr
			},
		}
		handler := auth.NewAuthHandler(mockUsecase)
		handler.SignInHandler(w, r)
		res := w.Result()
		defer res.Body.Close()

		var response domain.ErrorResponse
		decoder := json.NewDecoder(res.Body)
		err := decoder.Decode(&response)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
		assert.Equal(t, domain.ErrorResponse{Message: mockErr.Error()}, response)
	})
}
