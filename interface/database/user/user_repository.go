package user

import (
	"context"

	"github.com/Hajime3778/go-clean-arch/domain"
	"github.com/Hajime3778/go-clean-arch/interface/database"
)

type userRepository struct {
	SqlDriver database.SqlDriver
}

// NewUserRepository ユーザー機能のRepositoryオブジェクトを作成します
func NewUserRepository(sqlDriver database.SqlDriver) UserRepository {
	return &userRepository{sqlDriver}
}

func (u *userRepository) FetchByID(ctx context.Context, id int64) {
	panic("not implemented") // TODO: Implement
}

func (u *userRepository) FetchByEmailAndPassword(ctx context.Context, email string, password string) {
	panic("not implemented") // TODO: Implement
}

func (u *userRepository) Create(ctx context.Context, task domain.Task) error {
	panic("not implemented") // TODO: Implement
}
