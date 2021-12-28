package task_test

import (
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
		r := httptest.NewRequest(http.MethodGet, "http://localhost:8080/tasks/hogehoge", nil)
		w := httptest.NewRecorder()
		mockUsecase := &mock.MockTaskUsecase{
			MockGetByID: func(ctx context.Context, id int64) (domain.Task, error) {
				return domain.Task{}, errors.New("test error")
			},
		}
		handler := task.NewTaskHandler(mockUsecase)
		handler.Handler(w, r)
		res := w.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	})
}

func TestGetByID(t *testing.T) {
	t.Run("正常系 存在するIDで1件取得", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "http://localhost:8080/tasks/5", nil)
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
			t.Error(err)
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

	t.Run("異常系 Usecase実行時にデータが存在しないエラーが発生した場合、404エラーとなること", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "http://localhost:8080/tasks/5", nil)
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

		assert.Equal(t, http.StatusNotFound, res.StatusCode)
	})

	t.Run("異常系 Usecase実行時に想定外のエラーが発生した場合、500エラーとなること", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "http://localhost:8080/tasks/5", nil)
		w := httptest.NewRecorder()
		mockUsecase := &mock.MockTaskUsecase{
			MockGetByID: func(ctx context.Context, id int64) (domain.Task, error) {
				return domain.Task{}, errors.New("test error")
			},
		}
		handler := task.NewTaskHandler(mockUsecase)
		handler.Handler(w, r)
		res := w.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
	})
}
