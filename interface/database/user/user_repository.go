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

func (u *userRepository) GetByID(ctx context.Context, id int64) (domain.User, error) {
	panic("not implemented") // TODO: Implement
}

func (u *userRepository) GetByEmailAndPassword(ctx context.Context, email string, password string) (domain.User, error) {
	panic("not implemented") // TODO: Implement
}

func (ur *userRepository) Create(ctx context.Context, user domain.User) (int64, error) {
	query := `
		INSERT INTO users(name,email,password,salt) VALUES(?,?,?,?)
	`
	result, err := ur.SqlDriver.ExecuteContext(ctx, query, user.Name, user.Email, user.Password, user.Salt)
	if err != nil {
		return 0, err
	}

	createdId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return createdId, nil
}
