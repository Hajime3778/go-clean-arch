package task_test

import (
	"context"
	"testing"
	"time"

	"github.com/Hajime3778/go-clean-arch/domain"
	usecase "github.com/Hajime3778/go-clean-arch/usecase/task"
	"github.com/Hajime3778/go-clean-arch/usecase/task/mock"
	"github.com/stretchr/testify/assert"
)

func TestGetByID(t *testing.T) {
	t.Run("正常系 存在するIDで1件取得", func(t *testing.T) {
		mockTask := domain.Task{
			ID:        1,
			UserID:    1,
			Title:     "test titke",
			Content:   "test content",
			DueDate:   time.Now(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		mockTaskRepo := &mock.MockTaskRepo{
			MockGetByID: func(ctx context.Context, id int64) (domain.Task, error) {
				return mockTask, nil
			},
		}
		taskUsecase := usecase.NewTaskUsecase(mockTaskRepo)
		result, err := taskUsecase.GetByID(context.TODO(), mockTask.ID)

		assert.NoError(t, err)
		assert.Equal(t, mockTask, result)
	})

	t.Run("正常系 Repository実行時にエラーが発生した場合、エラーとなること", func(t *testing.T) {
		mockTaskRepo := &mock.MockTaskRepo{
			MockGetByID: func(ctx context.Context, id int64) (domain.Task, error) {
				return domain.Task{}, domain.ErrRecordNotFound
			},
		}
		taskUsecase := usecase.NewTaskUsecase(mockTaskRepo)
		result, err := taskUsecase.GetByID(context.TODO(), int64(1))

		assert.Equal(t, domain.ErrRecordNotFound, err)
		assert.Equal(t, domain.Task{}, result)
	})
}
