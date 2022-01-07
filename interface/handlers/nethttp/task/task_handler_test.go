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

func TestTaskHandlerTest(t *testing.T) {
	t.Run("異常系 パラメータが読み取れない場合 400エラーとなること", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "http://example.com/tasks/hogehoge", nil)
		w := httptest.NewRecorder()
		mockErr := errors.New("test error")
		mockUsecase := &mock.MockTaskUsecase{
			MockGetByID: func(ctx context.Context, id int64) (domain.Task, error) {
				return domain.Task{}, mockErr
			},
		}
		handler := task.NewTaskHandler(mockUsecase)
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

	t.Run("異常系 実装していないメソッドでリクエストした場合、404エラーとなること", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodOptions, "http://example.com/tasks/5", nil)
		w := httptest.NewRecorder()
		mockUsecase := &mock.MockTaskUsecase{}
		handler := task.NewTaskHandler(mockUsecase)
		handler.Handler(w, r)
		res := w.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusNotFound, res.StatusCode)
	})
}

func TestGetByID(t *testing.T) {
	t.Run("正常系 存在するIDで1件取得", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "http://example.com/tasks/5", nil)
		w := httptest.NewRecorder()
		mockTask := domain.Task{
			ID:        1,
			UserID:    1,
			Title:     "test title",
			Content:   "test content",
			DueDate:   time.Now(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		mockUsecase := &mock.MockTaskUsecase{
			MockGetByID: func(ctx context.Context, id int64) (domain.Task, error) {
				return mockTask, nil
			},
		}
		handler := task.NewTaskHandler(mockUsecase)
		handler.Handler(w, r)
		res := w.Result()
		defer res.Body.Close()

		var resTask domain.Task
		decoder := json.NewDecoder(res.Body)
		decoder.DisallowUnknownFields()
		err := decoder.Decode(&resTask)
		if err != nil {
			t.Fatal(err)
		}

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, res.StatusCode)

		assert.Equal(t, mockTask.ID, resTask.ID)
		assert.Equal(t, mockTask.Title, resTask.Title)
		assert.Equal(t, mockTask.Content, resTask.Content)
		assert.True(t, mockTask.DueDate.Equal(resTask.DueDate))
		assert.True(t, mockTask.CreatedAt.Equal(resTask.CreatedAt))
		assert.True(t, mockTask.UpdatedAt.Equal(resTask.UpdatedAt))
	})

	t.Run("準正常系 Usecase実行時にデータが存在しないエラーが発生した場合、404エラーとなること", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "http://example.com/tasks/5", nil)
		w := httptest.NewRecorder()
		mockUsecase := &mock.MockTaskUsecase{
			MockGetByID: func(ctx context.Context, id int64) (domain.Task, error) {
				return domain.Task{}, domain.ErrRecordNotFound
			},
		}
		handler := task.NewTaskHandler(mockUsecase)
		handler.Handler(w, r)
		res := w.Result()
		defer res.Body.Close()

		var resError domain.ErrorResponse
		decoder := json.NewDecoder(res.Body)
		err := decoder.Decode(&resError)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, http.StatusNotFound, res.StatusCode)
		assert.Equal(t, domain.ErrorResponse{Message: domain.ErrRecordNotFound.Error()}, resError)
	})

	t.Run("異常系 Usecase実行時に想定外のエラーが発生した場合、500エラーとなること", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "http://example.com/tasks/5", nil)
		w := httptest.NewRecorder()
		mockErr := errors.New("test error")
		mockUsecase := &mock.MockTaskUsecase{
			MockGetByID: func(ctx context.Context, id int64) (domain.Task, error) {
				return domain.Task{}, mockErr
			},
		}
		handler := task.NewTaskHandler(mockUsecase)
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

func TestUpdate(t *testing.T) {
	t.Run("正常系 1件更新", func(t *testing.T) {
		reqTask := task.UpdateTaskRequest{
			Title:   "test title",
			Content: "test content",
			DueDate: time.Now(),
		}
		byteTask, _ := json.Marshal(reqTask)
		r := httptest.NewRequest(http.MethodPut, "http://example.com/tasks/5",
			bytes.NewBuffer(byteTask),
		)
		w := httptest.NewRecorder()
		mockUsecase := &mock.MockTaskUsecase{
			MockUpdate: func(ctx context.Context, task domain.Task) error {
				return nil
			},
		}
		handler := task.NewTaskHandler(mockUsecase)
		handler.Handler(w, r)
		res := w.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusOK, res.StatusCode)
	})

	t.Run("準正常系 リクエストパラメータが足りていない場合、400エラーとなること", func(t *testing.T) {
		reqTask := task.UpdateTaskRequest{
			Title:   "test title",
			DueDate: time.Now(),
		}
		byteTask, _ := json.Marshal(reqTask)
		r := httptest.NewRequest(http.MethodPut, "http://example.com/tasks/5",
			bytes.NewBuffer(byteTask),
		)
		w := httptest.NewRecorder()
		mockUsecase := &mock.MockTaskUsecase{
			MockUpdate: func(ctx context.Context, task domain.Task) error {
				return nil
			},
		}
		handler := task.NewTaskHandler(mockUsecase)
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
		r := httptest.NewRequest(http.MethodPut, "http://example.com/tasks/5",
			bytes.NewBuffer(byteTask),
		)
		w := httptest.NewRecorder()
		mockUsecase := &mock.MockTaskUsecase{
			MockUpdate: func(ctx context.Context, task domain.Task) error {
				return nil
			},
		}
		handler := task.NewTaskHandler(mockUsecase)
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
		reqTask := task.UpdateTaskRequest{
			Title:   "test title",
			Content: "test content",
			DueDate: time.Now(),
		}
		byteTask, _ := json.Marshal(reqTask)
		r := httptest.NewRequest(http.MethodPut, "http://example.com/tasks/5",
			bytes.NewBuffer(byteTask),
		)
		mockErr := errors.New("test error")
		w := httptest.NewRecorder()
		mockUsecase := &mock.MockTaskUsecase{
			MockUpdate: func(ctx context.Context, task domain.Task) error {
				return mockErr
			},
		}
		handler := task.NewTaskHandler(mockUsecase)
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

func TestDelete(t *testing.T) {
	t.Run("正常系 1件削除", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodDelete, "http://example.com/tasks/5", nil)
		w := httptest.NewRecorder()
		mockUsecase := &mock.MockTaskUsecase{
			MockDelete: func(ctx context.Context, id int64) error {
				return nil
			},
		}
		handler := task.NewTaskHandler(mockUsecase)
		handler.Handler(w, r)
		res := w.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusNoContent, res.StatusCode)
	})

	t.Run("異常系 Usecase実行時にエラーが発生した場合、エラーとなること", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodDelete, "http://example.com/tasks/5", nil)
		w := httptest.NewRecorder()
		mockErr := errors.New("test error")
		mockUsecase := &mock.MockTaskUsecase{
			MockDelete: func(ctx context.Context, id int64) error {
				return mockErr
			},
		}
		handler := task.NewTaskHandler(mockUsecase)
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
