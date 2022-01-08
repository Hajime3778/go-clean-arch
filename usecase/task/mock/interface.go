package mock

import (
	"context"

	"github.com/Hajime3778/go-clean-arch/domain"
	usecase "github.com/Hajime3778/go-clean-arch/usecase/task"
)

type MockTaskUsecase struct {
	usecase.TaskUsecase
	MockFindByUserID func(ctx context.Context, limit int64, offset int64) ([]domain.Task, error)
	MockGetByID      func(ctx context.Context, id int64) (domain.Task, error)
	MockCreate       func(ctx context.Context, task domain.Task) error
	MockUpdate       func(ctx context.Context, task domain.Task) error
	MockDelete       func(ctx context.Context, id int64) error
}

func (m *MockTaskUsecase) FindByUserID(ctx context.Context, limit int64, offset int64) ([]domain.Task, error) {
	return m.MockFindByUserID(ctx, limit, offset)
}

func (m *MockTaskUsecase) GetByID(ctx context.Context, id int64) (domain.Task, error) {
	return m.MockGetByID(ctx, id)
}

func (m *MockTaskUsecase) Create(ctx context.Context, task domain.Task) error {
	return m.MockCreate(ctx, task)
}

func (m *MockTaskUsecase) Update(ctx context.Context, task domain.Task) error {
	return m.MockUpdate(ctx, task)
}

func (m *MockTaskUsecase) Delete(ctx context.Context, id int64) error {
	return m.MockDelete(ctx, id)
}
