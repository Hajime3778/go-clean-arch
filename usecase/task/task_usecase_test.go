package task_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Hajime3778/go-clean-arch/domain"
	"github.com/Hajime3778/go-clean-arch/interface/database/task/mock"
	usecase "github.com/Hajime3778/go-clean-arch/usecase/task"
	"github.com/stretchr/testify/assert"
)

func TestFindByUserID(t *testing.T) {
	t.Run("正常系 指定したユーザーIDで取得", func(t *testing.T) {
		mockTasks := createMockTasks(5, int64(1))
		mockTaskRepo := &mock.MockTaskRepo{
			MockFindByUserID: func(ctx context.Context, userID int64, limit int64, offset int64) ([]domain.Task, error) {
				return mockTasks, nil
			},
		}
		taskUsecase := usecase.NewTaskUsecase(mockTaskRepo)
		result, err := taskUsecase.FindByUserID(context.TODO(), int64(1), int64(1))

		assert.NoError(t, err)
		assert.Equal(t, mockTasks, result)
	})
	t.Run("異常系 Repository実行時にエラーが発生した場合、エラーとなること", func(t *testing.T) {
		mockTaskRepo := &mock.MockTaskRepo{
			MockFindByUserID: func(ctx context.Context, userID int64, limit int64, offset int64) ([]domain.Task, error) {
				return nil, domain.ErrInternalServerError
			},
		}
		taskUsecase := usecase.NewTaskUsecase(mockTaskRepo)
		result, err := taskUsecase.FindByUserID(context.TODO(), int64(1), int64(1))

		assert.Equal(t, domain.ErrInternalServerError, err)
		assert.Nil(t, result)
	})
}

func TestGetByID(t *testing.T) {
	t.Run("正常系 存在するIDで1件取得", func(t *testing.T) {
		mockTask := domain.Task{
			ID:        1,
			UserID:    1,
			Title:     "test title",
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

	t.Run("異常系 Repository実行時にエラーが発生した場合、エラーとなること", func(t *testing.T) {
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

func TestCreate(t *testing.T) {
	t.Run("正常系 1件追加", func(t *testing.T) {
		mockTaskRepo := &mock.MockTaskRepo{
			MockCreate: func(ctx context.Context, task domain.Task) error {
				return nil
			},
		}
		taskUsecase := usecase.NewTaskUsecase(mockTaskRepo)
		err := taskUsecase.Create(context.TODO(), domain.Task{})

		assert.NoError(t, err)
	})

	t.Run("異常系 Repository実行時にエラーが発生した場合、エラーとなること", func(t *testing.T) {
		mockTaskRepo := &mock.MockTaskRepo{

			MockCreate: func(ctx context.Context, task domain.Task) error {
				return domain.ErrInternalServerError
			},
		}
		taskUsecase := usecase.NewTaskUsecase(mockTaskRepo)
		err := taskUsecase.Create(context.TODO(), domain.Task{})

		assert.Equal(t, domain.ErrInternalServerError, err)
	})
}

func TestUpdate(t *testing.T) {
	t.Run("正常系 1件更新", func(t *testing.T) {
		mockTaskRepo := &mock.MockTaskRepo{
			MockGetByID: func(ctx context.Context, id int64) (domain.Task, error) {
				return domain.Task{}, nil
			},
			MockUpdate: func(ctx context.Context, task domain.Task) error {
				return nil
			},
		}
		taskUsecase := usecase.NewTaskUsecase(mockTaskRepo)
		err := taskUsecase.Update(context.TODO(), domain.Task{})

		assert.NoError(t, err)
	})

	t.Run("異常系 存在しないIDが指定された場合、エラーとなること", func(t *testing.T) {
		mockTaskRepo := &mock.MockTaskRepo{
			MockGetByID: func(ctx context.Context, id int64) (domain.Task, error) {
				return domain.Task{}, domain.ErrRecordNotFound
			},
		}
		taskUsecase := usecase.NewTaskUsecase(mockTaskRepo)
		err := taskUsecase.Update(context.TODO(), domain.Task{})

		assert.Equal(t, domain.ErrRecordNotFound, err)
	})

	t.Run("異常系 Repository実行時にエラーが発生した場合、エラーとなること", func(t *testing.T) {
		mockTaskRepo := &mock.MockTaskRepo{
			MockGetByID: func(ctx context.Context, id int64) (domain.Task, error) {
				return domain.Task{}, nil
			},
			MockUpdate: func(ctx context.Context, task domain.Task) error {
				return domain.ErrInternalServerError
			},
		}
		taskUsecase := usecase.NewTaskUsecase(mockTaskRepo)
		err := taskUsecase.Update(context.TODO(), domain.Task{})

		assert.Equal(t, domain.ErrInternalServerError, err)
	})
}

func TestDelete(t *testing.T) {
	t.Run("正常系 1件削除", func(t *testing.T) {
		mockTaskRepo := &mock.MockTaskRepo{
			MockDelete: func(ctx context.Context, id int64) error {
				return nil
			},
		}
		taskUsecase := usecase.NewTaskUsecase(mockTaskRepo)
		err := taskUsecase.Delete(context.TODO(), int64(1))

		assert.NoError(t, err)
	})

	t.Run("異常系 Repository実行時にエラーが発生した場合、エラーとなること", func(t *testing.T) {
		mockTaskRepo := &mock.MockTaskRepo{
			MockDelete: func(ctx context.Context, id int64) error {
				return domain.ErrInternalServerError
			},
		}
		taskUsecase := usecase.NewTaskUsecase(mockTaskRepo)
		err := taskUsecase.Delete(context.TODO(), int64(1))

		assert.Equal(t, domain.ErrInternalServerError, err)
	})
}

// createMockTasks モックのタスクを指定したユーザーIDで作成します
func createMockTasks(num int, userID int64) []domain.Task {
	mockTasks := make([]domain.Task, 0)
	for i := 0; i < num; i++ {
		id := int64(i + 1)
		task := domain.Task{
			ID:        id,
			UserID:    userID,
			Title:     fmt.Sprintf("test title%d", id),
			Content:   fmt.Sprintf("test content%d", id),
			DueDate:   time.Now(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		mockTasks = append(mockTasks, task)
	}
	return mockTasks
}
