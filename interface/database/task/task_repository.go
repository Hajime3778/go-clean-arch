package task

import (
	"context"

	"github.com/Hajime3778/go-clean-arch/domain"
	"github.com/Hajime3778/go-clean-arch/interface/database"
)

type taskRepository struct {
	SqlDriver database.SqlDriver
}

// NewTaskRepository タスク機能のRepositoryオブジェクトを作成します
func NewTaskRepository(sqlDriver database.SqlDriver) TaskRepository {
	return &taskRepository{sqlDriver}
}

// NewTaskUsecase タスクを指定した範囲まで取得します
func (tr *taskRepository) Fetch(ctx context.Context, cursor string, num int64) (res []domain.Task, nextCursor string, err error) {
	panic("not implemented") // TODO: Implement
}

// FetchByID IDでタスクを1件取得します
func (tr *taskRepository) FetchByID(ctx context.Context, id int64) (domain.Task, error) {
	query := `
		SELECT 
			* 
		FROM 
			tasks
		WHERE 
			id = ?
		ORDER BY id 
		LIMIT 1
	`
	row := tr.SqlDriver.QueryRow(query, id)

	task := domain.Task{}
	err := row.Scan(
		&task.ID,
		&task.Title,
		&task.Content,
		&task.DueDate,
		&task.UpdatedAt,
		&task.CreatedAt,
	)

	if err != nil {
		return task, err
	}

	return task, nil
}

// Create タスクを1件作成します
func (tr *taskRepository) Create(ctx context.Context, task domain.Task) error {
	query := `
		INSERT INTO tasks(title,content,due_date) VALUES(?,?,?)
	`

	_, err := tr.SqlDriver.Execute(query, task.Title, task.Content, task.DueDate)
	if err != nil {
		return err
	}

	return nil
}

// Update IDでタスクを1件更新します
func (tr *taskRepository) Update(ctx context.Context, task domain.Task) error {
	panic("not implemented") // TODO: Implement
}

// Delete IDでタスクを1件削除します
func (tr *taskRepository) Delete(ctx context.Context, id int64) error {
	panic("not implemented") // TODO: Implement
}
