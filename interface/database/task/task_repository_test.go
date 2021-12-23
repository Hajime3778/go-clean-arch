package task_test

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/Hajime3778/go-clean-arch/domain"
	infrastructure "github.com/Hajime3778/go-clean-arch/infrastructure/database"
	taskRepository "github.com/Hajime3778/go-clean-arch/interface/database/task"
	"github.com/stretchr/testify/assert"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestGetByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock error: '%s'", err)
	}

	sqlDriver := new(infrastructure.SqlDriver)
	sqlDriver.Conn = db

	mockTask := domain.Task{
		ID:        1,
		UserID:    1,
		Title:     "test title",
		Content:   "test content",
		DueDate:   time.Now(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	repo := taskRepository.NewTaskRepository(sqlDriver)

	t.Run("正常系 存在するIDで1件取得", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "user_id", "title", "content", "due_date", "updated_at", "created_at"}).
			AddRow(mockTask.ID, mockTask.UserID, mockTask.Title, mockTask.Content, mockTask.DueDate, mockTask.UpdatedAt, mockTask.CreatedAt)

		query := "SELECT * FROM tasks WHERE id = ?"
		mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(mockTask.ID).WillReturnRows(rows)

		got, err := repo.GetByID(context.TODO(), mockTask.ID)
		assert.NoError(t, err)
		assert.Equal(t, mockTask, got)
	})

	t.Run("準正常系 存在しないIDで検索してエラーとなること", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "user_id", "title", "content", "due_date", "updated_at", "created_at"})

		query := "SELECT * FROM tasks WHERE id = ?"
		mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(int64(2)).WillReturnRows(rows)

		got, err := repo.GetByID(context.TODO(), int64(2))
		assert.NoError(t, err)
		assert.Equal(t, domain.ErrRecordNotFound, err)
		assert.Equal(t, domain.Task{}, got)
	})
}
