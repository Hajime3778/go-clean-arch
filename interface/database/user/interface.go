package user

import (
	"context"

	"github.com/Hajime3778/go-clean-arch/domain"
)

// UserRepository
type UserRepository interface {
	GetByID(ctx context.Context, id int64) (domain.User, error)
	GetByEmailAndPassword(ctx context.Context, email string, password string) (domain.User, error)
	Create(ctx context.Context, user domain.User) (int64, error)
}
