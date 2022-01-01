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
	t.Run("正常系", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "http://example.com/tasks", nil)
		w := httptest.NewRecorder()
		mockUsecase := &mock.MockTaskUsecase{}
		handler := task.NewTaskIndexHandler(mockUsecase)
		handler.Handler(w, r)
		res := w.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusOK, res.StatusCode)
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

	t.Run("準正常系 リクエストパラメータが足りていない場合、エラーとなり400が返却されること", func(t *testing.T) {
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

		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	})

	t.Run("準正常系 リクエスト形式が間違っている場合、エラーとなり400が返却されること", func(t *testing.T) {
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

		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
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

		assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
	})
}
