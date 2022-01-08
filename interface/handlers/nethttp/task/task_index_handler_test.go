package task_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Hajime3778/go-clean-arch/domain"
	"github.com/Hajime3778/go-clean-arch/interface/handlers/nethttp/task"
	"github.com/Hajime3778/go-clean-arch/usecase/task/mock"
	"github.com/stretchr/testify/assert"
)

func TestTaskIndexHandlerTest(t *testing.T) {
	t.Run("異常系 実装していないメソッドでリクエストした場合、404エラーとなること", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodOptions, "http://example.com/tasks", nil)
		w := httptest.NewRecorder()
		mockUsecase := &mock.MockTaskUsecase{}
		handler := task.NewTaskIndexHandler(mockUsecase)
		handler.Handler(w, r)
		res := w.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusNotFound, res.StatusCode)
	})
}

func TestFindByIDTest(t *testing.T) {
	t.Run("正常系 複数取得", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "http://example.com/tasks?limit=10&offset=0", nil)
		w := httptest.NewRecorder()
		mockUsecase := &mock.MockTaskUsecase{
			MockFindByUserID: func(ctx context.Context, limit int64, offset int64) ([]domain.Task, error) {
				return make([]domain.Task, 0), nil
			},
		}
		handler := task.NewTaskIndexHandler(mockUsecase)
		handler.Handler(w, r)
		res := w.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusOK, res.StatusCode)
	})

	t.Run("準正常系 パラメータが指定されていない場合、400エラーとなること", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "http://example.com/tasks", nil)
		w := httptest.NewRecorder()
		mockUsecase := &mock.MockTaskUsecase{
			MockFindByUserID: func(ctx context.Context, limit int64, offset int64) ([]domain.Task, error) {
				return make([]domain.Task, 0), nil
			},
		}
		handler := task.NewTaskIndexHandler(mockUsecase)
		handler.Handler(w, r)
		res := w.Result()
		defer res.Body.Close()

		var resError domain.ErrorResponse
		decoder := json.NewDecoder(res.Body)
		err := decoder.Decode(&resError)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
		assert.NotEmpty(t, resError.Message)
	})

	t.Run("準正常系 limitが数字でない場合、400エラーとなること", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "http://example.com/tasks?limit=foo&offset=0", nil)
		w := httptest.NewRecorder()
		mockUsecase := &mock.MockTaskUsecase{
			MockFindByUserID: func(ctx context.Context, limit int64, offset int64) ([]domain.Task, error) {
				return make([]domain.Task, 0), nil
			},
		}
		handler := task.NewTaskIndexHandler(mockUsecase)
		handler.Handler(w, r)
		res := w.Result()
		defer res.Body.Close()

		var resError domain.ErrorResponse
		decoder := json.NewDecoder(res.Body)
		err := decoder.Decode(&resError)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
		assert.NotEmpty(t, resError.Message)
	})

	t.Run("準正常系 offsetが数字でない場合、400エラーとなること", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "http://example.com/tasks?limit=0&offset=foo", nil)
		w := httptest.NewRecorder()
		mockUsecase := &mock.MockTaskUsecase{
			MockFindByUserID: func(ctx context.Context, limit int64, offset int64) ([]domain.Task, error) {
				return make([]domain.Task, 0), nil
			},
		}
		handler := task.NewTaskIndexHandler(mockUsecase)
		handler.Handler(w, r)
		res := w.Result()
		defer res.Body.Close()

		var resError domain.ErrorResponse
		decoder := json.NewDecoder(res.Body)
		err := decoder.Decode(&resError)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
		assert.NotEmpty(t, resError.Message)
	})

	t.Run("異常系 Usecase実行時にエラーが発生した場合、エラーとなること", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "http://example.com/tasks?limit=0&offset=0", nil)
		w := httptest.NewRecorder()
		mockErr := errors.New("test error")
		mockUsecase := &mock.MockTaskUsecase{
			MockFindByUserID: func(ctx context.Context, limit int64, offset int64) ([]domain.Task, error) {
				return nil, mockErr
			},
		}
		handler := task.NewTaskIndexHandler(mockUsecase)
		handler.Handler(w, r)
		res := w.Result()
		defer res.Body.Close()

		var resError domain.ErrorResponse
		decoder := json.NewDecoder(res.Body)
		err := decoder.Decode(&resError)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
		assert.Equal(t, domain.ErrorResponse{Message: mockErr.Error()}, resError)
	})
}

func TestCreate(t *testing.T) {
	t.Run("正常系 1件追加", func(t *testing.T) {
		reqTask := task.CreateTaskRequest{
			Title:   "test title",
			Content: "test content",
			DueDate: time.Now(),
		}
		byteTask, _ := json.Marshal(reqTask)
		r := httptest.NewRequest(http.MethodPost, "http://example.com/tasks",
			bytes.NewBuffer(byteTask),
		)
		w := httptest.NewRecorder()
		mockUsecase := &mock.MockTaskUsecase{
			MockCreate: func(ctx context.Context, task domain.Task) error {
				return nil
			},
		}
		handler := task.NewTaskIndexHandler(mockUsecase)
		handler.Handler(w, r)
		res := w.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusCreated, res.StatusCode)
	})

	t.Run("準正常系 リクエストパラメータが足りていない場合、400エラーとなること", func(t *testing.T) {
		reqTask := task.CreateTaskRequest{
			Title:   "test title",
			DueDate: time.Now(),
		}
		byteTask, _ := json.Marshal(reqTask)
		r := httptest.NewRequest(http.MethodPost, "http://example.com/tasks",
			bytes.NewBuffer(byteTask),
		)
		w := httptest.NewRecorder()
		mockUsecase := &mock.MockTaskUsecase{
			MockCreate: func(ctx context.Context, task domain.Task) error {
				return nil
			},
		}
		handler := task.NewTaskIndexHandler(mockUsecase)
		handler.Handler(w, r)
		res := w.Result()
		defer res.Body.Close()

		var resError domain.ErrorResponse
		decoder := json.NewDecoder(res.Body)
		err := decoder.Decode(&resError)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
		assert.NotEmpty(t, resError.Message)
	})

	t.Run("準正常系 リクエスト形式が間違っている場合、400エラーとなること", func(t *testing.T) {
		req := domain.ErrorResponse{
			Message: "test",
		}
		byteTask, _ := json.Marshal(req)
		r := httptest.NewRequest(http.MethodPost, "http://example.com/tasks",
			bytes.NewBuffer(byteTask),
		)
		w := httptest.NewRecorder()
		mockUsecase := &mock.MockTaskUsecase{
			MockCreate: func(ctx context.Context, task domain.Task) error {
				return nil
			},
		}
		handler := task.NewTaskIndexHandler(mockUsecase)
		handler.Handler(w, r)
		res := w.Result()
		defer res.Body.Close()

		var resError domain.ErrorResponse
		decoder := json.NewDecoder(res.Body)
		err := decoder.Decode(&resError)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
		assert.NotEmpty(t, resError.Message)
	})

	t.Run("異常系 Usecase実行時にエラーが発生した場合、エラーとなること", func(t *testing.T) {
		reqTask := task.CreateTaskRequest{
			Title:   "test title",
			Content: "test content",
			DueDate: time.Now(),
		}
		byteTask, _ := json.Marshal(reqTask)
		r := httptest.NewRequest(http.MethodPost, "http://example.com/tasks",
			bytes.NewBuffer(byteTask),
		)
		w := httptest.NewRecorder()
		mockUsecase := &mock.MockTaskUsecase{
			MockCreate: func(ctx context.Context, task domain.Task) error {
				return errors.New("test error")
			},
		}
		handler := task.NewTaskIndexHandler(mockUsecase)
		handler.Handler(w, r)
		res := w.Result()
		defer res.Body.Close()

		var resError domain.ErrorResponse
		decoder := json.NewDecoder(res.Body)
		err := decoder.Decode(&resError)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
		assert.NotEmpty(t, resError.Message)
	})
}
