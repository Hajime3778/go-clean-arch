package task

import (
	"context"
	"log"

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

// GetByID IDでタスクを1件取得します
func (tr *taskRepository) GetByID(ctx context.Context, id int64) (task domain.Task, err error) {
	query := `
		SELECT 
			* 
		FROM 
			tasks
		WHERE 
			id = ?
	`
	rows, err := tr.SqlDriver.QueryContext(ctx, query, id)
	if err != nil {
		return domain.Task{}, err
	}

	defer func() {
		err := rows.Close()
		if err != nil {
			log.Println(err.Error())
		}
	}()

	if !rows.Next() {
		return domain.Task{}, domain.ErrRecordNotFound
	}

	err = rows.Scan(
		&task.ID,
		&task.UserID,
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
		INSERT INTO tasks(user_id,title,content,due_date) VALUES(?,?,?,?)
	`
	_, err := tr.SqlDriver.ExecuteContext(ctx, query, task.UserID, task.Title, task.Content, task.DueDate)
	if err != nil {
		return err
	}

	return nil
}

// Update IDでタスクを1件更新します
func (tr *taskRepository) Update(ctx context.Context, task domain.Task) error {
	query := `
		UPDATE tasks SET title = ?, content = ?, due_date = ? where id = ? 
	`
	_, err := tr.SqlDriver.ExecuteContext(ctx, query, task.Title, task.Content, task.DueDate, task.ID)
	if err != nil {
		return err
	}

	return nil
}

// Delete IDでタスクを1件削除します
func (tr *taskRepository) Delete(ctx context.Context, id int64) error {
	query := `
		DELETE FROM tasks where id = ? 
	`
	_, err := tr.SqlDriver.ExecuteContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}
