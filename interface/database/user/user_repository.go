package user

import (
	"context"
	"log"

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

func (ur *userRepository) GetByID(ctx context.Context, id int64) (domain.User, error) {
	query := `
		SELECT 
			* 
		FROM 
			users
		WHERE 
			id = ?
	`
	rows, err := ur.SqlDriver.QueryContext(ctx, query, id)
	if err != nil {
		return domain.User{}, err
	}

	defer func() {
		err := rows.Close()
		if err != nil {
			log.Println(err.Error())
		}
	}()

	if !rows.Next() {
		return domain.User{}, domain.ErrRecordNotFound
	}

	user := domain.User{}
	err = rows.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.Salt,
		&user.UpdatedAt,
		&user.CreatedAt,
	)

	if err != nil {
		return user, err
	}

	return user, nil
}

func (ur *userRepository) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	query := `
		SELECT 
			* 
		FROM 
			users
		WHERE 
			email = ?
	`
	rows, err := ur.SqlDriver.QueryContext(ctx, query, email)
	if err != nil {
		return domain.User{}, err
	}

	defer func() {
		err := rows.Close()
		if err != nil {
			log.Println(err.Error())
		}
	}()

	if !rows.Next() {
		return domain.User{}, domain.ErrRecordNotFound
	}

	user := domain.User{}
	err = rows.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.Salt,
		&user.UpdatedAt,
		&user.CreatedAt,
	)

	if err != nil {
		return user, err
	}

	return user, nil
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
