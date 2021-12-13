package task

import (
	"context"

	"github.com/Hajime3778/go-clean-arch/domain"
)

// TaskRepository
type TaskRepository interface {
	Fetch(ctx context.Context, cursor string, num int64) (tasks []domain.Task, nextCursor string, err error)
	FetchByID(ctx context.Context, id int64) (task domain.Task, err error)
	Create(ctx context.Context, task domain.Task) error
	Update(ctx context.Context, task domain.Task) error
	Delete(ctx context.Context, id int64) error
}
