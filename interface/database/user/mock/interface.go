package mock

import (
	"context"

	"github.com/Hajime3778/go-clean-arch/domain"
	repo "github.com/Hajime3778/go-clean-arch/interface/database/user"
)

type MockUserRepo struct {
	repo.UserRepository
	MockGetByID    func(ctx context.Context, id int64) (domain.User, error)
	MockGetByEmail func(ctx context.Context, email string) (domain.User, error)
	MockCreate     func(ctx context.Context, task domain.User) (int64, error)
}

func (m *MockUserRepo) GetByID(ctx context.Context, id int64) (domain.User, error) {
	return m.MockGetByID(ctx, id)
}

func (m *MockUserRepo) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	return m.MockGetByEmail(ctx, email)
}

func (m *MockUserRepo) Create(ctx context.Context, user domain.User) (int64, error) {
	return m.MockCreate(ctx, user)
}
