package task

import (
	"context"

	"github.com/Hajime3778/go-clean-arch/domain"
)

type TaskUsecase interface {
	Fetch(ctx context.Context, cursor string, num int64) ([]domain.Task, string, error)
	FetchByID(ctx context.Context, id int64) (domain.Task, error)
	Create(ctx context.Context, task domain.Task) error
	Update(ctx context.Context, task domain.Task) error
	Delete(ctx context.Context, id int64) error
}
