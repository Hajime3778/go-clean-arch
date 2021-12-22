package user

import (
	"context"

	"github.com/Hajime3778/go-clean-arch/domain"
)

type UserRepository interface {
	FetchByID(ctx context.Context, id int64)
	FetchByEmailAndPassword(ctx context.Context, email string, password string)
	Create(ctx context.Context, task domain.Task) error
}
