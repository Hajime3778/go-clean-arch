package task_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/Hajime3778/go-clean-arch/domain"
	infrastructure "github.com/Hajime3778/go-clean-arch/infrastructure/database"
	taskRepository "github.com/Hajime3778/go-clean-arch/interface/database/task"
	"github.com/stretchr/testify/assert"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestFindByUserID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock error: '%s'", err)
	}

	sqlDriver := new(infrastructure.SqlDriver)
	sqlDriver.Conn = db

	repo := taskRepository.NewTaskRepository(sqlDriver)

	t.Run("正常系 指定したユーザーIDのみすべて取得されていること", func(t *testing.T) {
		userID := int64(1)

		rows := sqlmock.NewRows([]string{"id", "user_id", "title", "content", "due_date", "updated_at", "created_at"})

		// ユーザーID=1でデータを5件作成
		mockTasks := createMockTasks(5, userID)
		for _, mockTask := range mockTasks {
			rows.AddRow(mockTask.ID, mockTask.UserID, mockTask.Title, mockTask.Content, mockTask.DueDate, mockTask.UpdatedAt, mockTask.CreatedAt)
		}
		//rows.AddRow(int64(6), int64(2), "dummy title", "dummy content", time.Now(), time.Now(), time.Now())

		query := `
			SELECT
				*
			FROM
				tasks
			WHERE
				user_id = ?
			ORDER BY
				due_date
			LIMIT ? OFFSET ?
		`
		mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(userID, 3, 1).WillReturnRows(rows)

		got, err := repo.FindByUserID(context.TODO(), userID, 3, 1)
		assert.NoError(t, err)
		assert.Equal(t, mockTasks, got)
	})

	t.Run("正常系 指定した範囲のデータが取得できていること", func(t *testing.T) {
	})

	t.Run("準正常系 データが存在しない場合、エラーとならないこと", func(t *testing.T) {
		userID := int64(1)

		rows := sqlmock.NewRows([]string{"id", "user_id", "title", "content", "due_date", "updated_at", "created_at"})
		query := `
			SELECT
				*
			FROM
				tasks
			WHERE
				user_id = ?
			ORDER BY
				due_date
			LIMIT ? OFFSET ?
		`
		mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(userID, 5, 0).WillReturnRows(rows)

		got, err := repo.FindByUserID(context.TODO(), userID, 5, 0)
		assert.NoError(t, err)
		assert.Equal(t, []domain.Task{}, got)
	})
}

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
		assert.Equal(t, domain.ErrRecordNotFound, err)
		assert.Equal(t, domain.Task{}, got)
	})
}

func TestCreate(t *testing.T) {
	// db, mock, err := sqlmock.New()
	// if err != nil {
	// 	t.Fatalf("sqlmock error: '%s'", err)
	// }

	// sqlDriver := new(infrastructure.SqlDriver)
	// sqlDriver.Conn = db

	// repo := taskRepository.NewTaskRepository(sqlDriver)

	t.Run("正常系 存在するIDで1件取得", func(t *testing.T) {
	})

	t.Run("準正常系 存在しないIDで検索してエラーとなること", func(t *testing.T) {
	})
}

func TestUpdate(t *testing.T) {
	// db, mock, err := sqlmock.New()
	// if err != nil {
	// 	t.Fatalf("sqlmock error: '%s'", err)
	// }

	// sqlDriver := new(infrastructure.SqlDriver)
	// sqlDriver.Conn = db

	// repo := taskRepository.NewTaskRepository(sqlDriver)

	t.Run("正常系 存在するIDで1件取得", func(t *testing.T) {
	})

	t.Run("準正常系 存在しないIDで検索してエラーとなること", func(t *testing.T) {
	})
}

func TestDelete(t *testing.T) {
	// db, mock, err := sqlmock.New()
	// if err != nil {
	// 	t.Fatalf("sqlmock error: '%s'", err)
	// }

	// sqlDriver := new(infrastructure.SqlDriver)
	// sqlDriver.Conn = db

	// repo := taskRepository.NewTaskRepository(sqlDriver)

	t.Run("正常系 存在するIDで1件取得", func(t *testing.T) {
	})

	t.Run("準正常系 存在しないIDで検索してエラーとなること", func(t *testing.T) {
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
