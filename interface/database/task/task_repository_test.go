package task_test

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/Hajime3778/go-clean-arch/domain"
	infrastructure "github.com/Hajime3778/go-clean-arch/infrastructure/database"
	"github.com/Hajime3778/go-clean-arch/interface/database"
	mockSqlDriver "github.com/Hajime3778/go-clean-arch/interface/database/mock"
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

	t.Run("正常系 指定したユーザーIDで取得", func(t *testing.T) {
		userID := int64(1)
		rows := sqlmock.NewRows([]string{"id", "user_id", "title", "content", "due_date", "updated_at", "created_at"})
		mockTasks := createMockTasks(5, userID)
		for _, mockTask := range mockTasks {
			rows.AddRow(mockTask.ID, mockTask.UserID, mockTask.Title, mockTask.Content, mockTask.DueDate, mockTask.UpdatedAt, mockTask.CreatedAt)
		}
		mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(userID, 3, 1).WillReturnRows(rows)

		got, err := repo.FindByUserID(context.TODO(), userID, 3, 1)
		assert.NoError(t, err)
		assert.Equal(t, mockTasks, got)
	})

	t.Run("準正常系 データが存在しない場合、エラーとならないこと", func(t *testing.T) {
		userID := int64(1)
		rows := sqlmock.NewRows([]string{"id", "user_id", "title", "content", "due_date", "updated_at", "created_at"})
		mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(userID, 5, 0).WillReturnRows(rows)

		got, err := repo.FindByUserID(context.TODO(), userID, 5, 0)
		assert.NoError(t, err)
		assert.Equal(t, []domain.Task{}, got)
	})

	t.Run("異常系 クエリ実行で失敗した場合エラーが返却されること", func(t *testing.T) {
		userID := int64(1)
		mockErr := errors.New("query failed error")
		mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(userID, 5, 0).
			WillReturnError(mockErr)

		got, err := repo.FindByUserID(context.TODO(), userID, 5, 0)
		assert.Equal(t, mockErr, err)
		assert.Nil(t, got)
	})

	t.Run("異常系 Scan実行で失敗した場合エラーが返却されること", func(t *testing.T) {
		userID := int64(1)
		rows := sqlmock.NewRows([]string{"foo"}).AddRow("bar")
		mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(userID, 5, 0).WillReturnRows(rows)

		got, err := repo.FindByUserID(context.TODO(), userID, 5, 0)
		assert.NotNil(t, err)
		assert.Nil(t, got)
	})

	t.Run("異常系 Rows.Close実行で失敗した場合にログが出力されること", func(t *testing.T) {
		userID := int64(1)
		mockErr := errors.New("rows close error")
		rows := sqlmock.NewRows([]string{"foo"}).AddRow("bar").CloseError(mockErr)
		mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(userID, 5, 0).WillReturnRows(rows)

		got, err := repo.FindByUserID(context.TODO(), userID, 5, 0)
		assert.NotNil(t, err)
		assert.Nil(t, got)
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
	query := "SELECT * FROM tasks WHERE id = ?"

	t.Run("正常系 存在するIDで1件取得", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "user_id", "title", "content", "due_date", "updated_at", "created_at"}).
			AddRow(mockTask.ID, mockTask.UserID, mockTask.Title, mockTask.Content, mockTask.DueDate, mockTask.UpdatedAt, mockTask.CreatedAt)
		mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(mockTask.ID).WillReturnRows(rows)

		got, err := repo.GetByID(context.TODO(), mockTask.ID)
		assert.NoError(t, err)
		assert.Equal(t, mockTask, got)
	})

	t.Run("準正常系 存在しないIDで検索してエラーとなること", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "user_id", "title", "content", "due_date", "updated_at", "created_at"})
		mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(int64(2)).WillReturnRows(rows)

		got, err := repo.GetByID(context.TODO(), int64(2))
		assert.Equal(t, domain.ErrRecordNotFound, err)
		assert.Equal(t, domain.Task{}, got)
	})

	t.Run("異常系 クエリ実行で失敗した場合エラーが返却されること", func(t *testing.T) {
		mockErr := errors.New("query failed error")
		mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(int64(2)).WillReturnError(mockErr)

		got, err := repo.GetByID(context.TODO(), int64(2))
		assert.Equal(t, mockErr, err)
		assert.Equal(t, domain.Task{}, got)
	})

	t.Run("異常系 Scan実行で失敗した場合エラーが返却されること", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"foo"}).AddRow("bar")
		mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(int64(2)).WillReturnRows(rows)

		got, err := repo.GetByID(context.TODO(), int64(2))
		assert.NotNil(t, err)
		assert.Equal(t, domain.Task{}, got)
	})

	t.Run("異常系 Rows.Close実行で失敗した場合にログが出力されること", func(t *testing.T) {
		mockErr := errors.New("rows close error")
		rows := sqlmock.NewRows([]string{"foo"}).AddRow("bar").CloseError(mockErr)
		mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(int64(2)).WillReturnRows(rows)

		got, err := repo.GetByID(context.TODO(), int64(2))
		assert.NotNil(t, err)
		assert.Equal(t, domain.Task{}, got)
	})
}

func TestCreate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock error: '%s'", err)
	}

	t.Run("正常系 1件追加", func(t *testing.T) {
		sqlDriver := new(infrastructure.SqlDriver)
		sqlDriver.Conn = db

		repo := taskRepository.NewTaskRepository(sqlDriver)
		query := "INSERT INTO tasks(user_id,title,content,due_date) VALUES(?,?,?,?)"

		mockTask := domain.Task{
			UserID:  1,
			Title:   "test title",
			Content: "test content",
			DueDate: time.Now(),
		}
		prep := mock.ExpectPrepare(regexp.QuoteMeta(query))
		prep.ExpectExec().
			WithArgs(mockTask.UserID, mockTask.Title, mockTask.Content, mockTask.DueDate).
			WillReturnResult(sqlmock.NewResult(12, 1))

		id, err := repo.Create(context.TODO(), mockTask)
		assert.NoError(t, err)
		assert.NotEqual(t, int64(0), id)
	})

	t.Run("異常系 クエリ実行で失敗した場合エラーが返却されること", func(t *testing.T) {
		sqlDriver := new(infrastructure.SqlDriver)
		sqlDriver.Conn = db

		repo := taskRepository.NewTaskRepository(sqlDriver)
		query := "INSERT INTO tasks(user_id,title,content,due_date) VALUES(?,?,?,?)"

		mockTask := domain.Task{
			UserID:    1,
			Title:     "test title",
			Content:   "test content",
			DueDate:   time.Now(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		mockErr := errors.New("query failed error")
		prep := mock.ExpectPrepare(regexp.QuoteMeta(query))
		prep.ExpectExec().
			WithArgs(mockTask.UserID, mockTask.Title, mockTask.Content, mockTask.DueDate).
			WillReturnError(mockErr)

		id, err := repo.Create(context.TODO(), mockTask)
		assert.Equal(t, mockErr, err)
		assert.Equal(t, int64(0), id)
	})

	t.Run("異常系 追加後IDで失敗した場合エラーが返却されtること", func(t *testing.T) {

		mockErr := errors.New("test error")
		mockResult := &mockSqlDriver.MockResult{
			MockLastInsertId: func() (int64, error) {
				return 0, mockErr
			},
		}
		mockDriver := &mockSqlDriver.MockSqlDriver{
			MockExecuteContext: func(context.Context, string, ...interface{}) (database.Result, error) {
				return mockResult, nil
			},
		}
		mockDriver.Conn = db

		repo := taskRepository.NewTaskRepository(mockDriver)
		query := "INSERT INTO tasks(user_id,title,content,due_date) VALUES(?,?,?,?)"

		mockTask := domain.Task{
			UserID:    1,
			Title:     "test title",
			Content:   "test content",
			DueDate:   time.Now(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		prep := mock.ExpectPrepare(regexp.QuoteMeta(query))
		prep.ExpectExec().
			WithArgs(mockTask.UserID, mockTask.Title, mockTask.Content, mockTask.DueDate).
			WillReturnResult(sqlmock.NewResult(12, 1))

		id, err := repo.Create(context.TODO(), mockTask)
		assert.Equal(t, mockErr, err)
		assert.Equal(t, int64(0), id)
	})
}

func TestUpdate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock error: '%s'", err)
	}

	sqlDriver := new(infrastructure.SqlDriver)
	sqlDriver.Conn = db

	repo := taskRepository.NewTaskRepository(sqlDriver)
	query := "UPDATE tasks SET title = ?, content = ?, due_date = ? where id = ?"

	t.Run("正常系 1件更新", func(t *testing.T) {
		mockTask := domain.Task{
			ID:      1,
			Title:   "test title",
			Content: "test content",
			DueDate: time.Now(),
		}
		prep := mock.ExpectPrepare(regexp.QuoteMeta(query))
		prep.ExpectExec().
			WithArgs(mockTask.Title, mockTask.Content, mockTask.DueDate, mockTask.ID).
			WillReturnResult(sqlmock.NewResult(12, 1))

		err = repo.Update(context.TODO(), mockTask)
		assert.NoError(t, err)
	})

	t.Run("異常系 クエリ実行で失敗した場合エラーが返却されること", func(t *testing.T) {
		mockTask := domain.Task{
			ID:      1,
			Title:   "test title",
			Content: "test content",
			DueDate: time.Now(),
		}
		mockErr := errors.New("query failed error")
		prep := mock.ExpectPrepare(regexp.QuoteMeta(query))
		prep.ExpectExec().
			WithArgs(mockTask.Title, mockTask.Content, mockTask.DueDate, mockTask.ID).
			WillReturnError(mockErr)

		err = repo.Update(context.TODO(), mockTask)
		assert.Equal(t, mockErr, err)
	})
}

func TestDelete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock error: '%s'", err)
	}

	sqlDriver := new(infrastructure.SqlDriver)
	sqlDriver.Conn = db

	repo := taskRepository.NewTaskRepository(sqlDriver)
	query := "DELETE FROM tasks where id = ? "

	t.Run("正常系 1件削除", func(t *testing.T) {
		prep := mock.ExpectPrepare(regexp.QuoteMeta(query))
		prep.ExpectExec().
			WithArgs(int64(1)).
			WillReturnResult(sqlmock.NewResult(12, 1))

		err = repo.Delete(context.TODO(), int64(1))
		assert.NoError(t, err)
	})

	t.Run("異常系 クエリ実行で失敗した場合エラーが返却されること", func(t *testing.T) {
		mockErr := errors.New("query failed error")
		prep := mock.ExpectPrepare(regexp.QuoteMeta(query))
		prep.ExpectExec().
			WithArgs(int64(1)).
			WillReturnError(mockErr)

		err = repo.Delete(context.TODO(), int64(1))
		assert.Equal(t, mockErr, err)
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
