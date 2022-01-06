package user

import (
	"context"

	"github.com/Hajime3778/go-clean-arch/domain"
)

// UserRepository
type UserRepository interface {
	GetByID(ctx context.Context, id int64) (domain.User, error)
	Create(ctx context.Context, user domain.User) (int64, error)
	Update(ctx context.Context, user domain.User) error
	Delete(ctx context.Context, id int64) error
}
