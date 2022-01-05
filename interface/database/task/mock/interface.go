package mock

import (
	"context"

	"github.com/Hajime3778/go-clean-arch/domain"
	repo "github.com/Hajime3778/go-clean-arch/interface/database/task"
)

type MockTaskRepo struct {
	repo.TaskRepository
	MockFindByUserID func(ctx context.Context, userID int64, limit int64, offset int64) ([]domain.Task, error)
	MockGetByID      func(ctx context.Context, id int64) (domain.Task, error)
	MockCreate       func(ctx context.Context, task domain.Task) (int64, error)
	MockUpdate       func(ctx context.Context, task domain.Task) error
	MockDelete       func(ctx context.Context, id int64) error
}

func (m *MockTaskRepo) FindByUserID(ctx context.Context, userID int64, limit int64, offset int64) ([]domain.Task, error) {
	return m.MockFindByUserID(ctx, userID, limit, offset)
}

func (m *MockTaskRepo) GetByID(ctx context.Context, id int64) (domain.Task, error) {
	return m.MockGetByID(ctx, id)
}

func (m *MockTaskRepo) Create(ctx context.Context, task domain.Task) (int64, error) {
	return m.MockCreate(ctx, task)
}

func (m *MockTaskRepo) Update(ctx context.Context, task domain.Task) error {
	return m.MockUpdate(ctx, task)
}

func (m *MockTaskRepo) Delete(ctx context.Context, id int64) error {
	return m.MockDelete(ctx, id)
}
