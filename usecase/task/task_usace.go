package task

import (
	"context"

	"github.com/Hajime3778/go-clean-arch/domain"
	_repository "github.com/Hajime3778/go-clean-arch/interface/database/task"
)

type taskUsecase struct {
	repo _repository.TaskRepository
}

// NewTaskUsecase タスク機能のUsecaseオブジェクトを作成します
func NewTaskUsecase(repo _repository.TaskRepository) TaskUsecase {
	return &taskUsecase{repo}
}

// NewTaskUsecase タスクを指定した範囲まで取得します
func (tu *taskUsecase) Fetch(ctx context.Context, cursor string, num int64) ([]domain.Task, string, error) {
	panic("not implemented") // TODO: Implement
}

// FetchByID IDでタスクを1件取得します
func (tu *taskUsecase) FetchByID(ctx context.Context, id int64) (domain.Task, error) {
	task, err := tu.repo.FetchByID(ctx, id)
	if err != nil {
		return domain.Task{}, err
	}
	return task, nil
}

// Create タスクを1件作成します
func (tu *taskUsecase) Create(ctx context.Context, task domain.Task) error {
	err := tu.repo.Create(ctx, task)
	if err != nil {
		return err
	}
	return nil
}

// Update IDでタスクを1件更新します
func (tu *taskUsecase) Update(ctx context.Context, task domain.Task) error {
	panic("not implemented") // TODO: Implement
}

// Delete IDでタスクを1件削除します
func (tu *taskUsecase) Delete(ctx context.Context, id int64) error {
	panic("not implemented") // TODO: Implement
}
